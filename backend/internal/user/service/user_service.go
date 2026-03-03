package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"swiftly/backend/internal/database"
	"swiftly/backend/internal/pkg/auth"
	"swiftly/backend/internal/pkg/captcha"
	"swiftly/backend/internal/pkg/sanitizer"
	"swiftly/backend/internal/pkg/socialauth"
	"swiftly/backend/internal/pkg/storage"
	"swiftly/backend/internal/user"
	"swiftly/backend/internal/user/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo         *repository.Repository
	activityRepo *repository.ActivityRepository
	uploader     storage.Uploader
}

func NewService(repo *repository.Repository, activityRepo *repository.ActivityRepository, uploader storage.Uploader) *Service {
	return &Service{
		repo:         repo,
		activityRepo: activityRepo,
		uploader:     uploader,
	}
}

func (s *Service) auditLog(userID, activityType, ip, userAgent string, metadata map[string]interface{}) {
	_ = s.activityRepo.Log(userID, activityType, ip, userAgent, metadata)
}

func (s *Service) CreateUser(req user.CreateUserRequest, ip, userAgent string) (*user.TokenResponse, error) {
	valid, err := captcha.VerifyToken(req.CaptchaToken)
	if err != nil || !valid {
		return nil, errors.New("bot detection failed. please try again")
	}

	cleanEmail := sanitizer.Email(req.Email)
	cleanFullName := sanitizer.Text(req.FullName)
	cleanUsername := sanitizer.Username(req.Username)
	cleanPhone := sanitizer.Phone(req.PhoneNumber)

	if cleanEmail == "" || cleanFullName == "" || req.Password == "" {
		return nil, errors.New("email, full name, and password are required")
	}

	exists, _ := s.IsEmailAvailable(cleanEmail)
	if !exists {
		return nil, errors.New("email is already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &user.User{
		Email:       cleanEmail,
		Username:    cleanUsername,
		PhoneNumber: cleanPhone,
		FullName:    cleanFullName,
		Password:    string(hashedPassword),
		Role:        "customer",
		Status:      "active",
	}

	err = s.repo.Create(u)
	if err != nil {
		return nil, err
	}

	s.GenerateAndStoreOTP(cleanEmail)
	s.auditLog(u.ID, "USER_REGISTERED", ip, userAgent, nil)

	sessionID := uuid.New().String()
	accessToken, refreshToken, err := auth.GenerateTokens(u.ID, sessionID)
	if err != nil {
		return nil, err
	}

	// Store session in Redis
	s.storeSession(u.ID, sessionID, refreshToken, ip, userAgent)

	return &user.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) GenerateAndStoreOTP(email string) (string, error) {
	b := make([]byte, 3)
	rand.Read(b)
	otp := fmt.Sprintf("%06d", (uint32(b[0])<<16|uint32(b[1])<<8|uint32(b[2]))%1000000)

	ctx := context.Background()
	rdb := database.GetRedis()
	err := rdb.Set(ctx, "otp:"+email, otp, 5*time.Minute).Err()
	if err != nil {
		return "", err
	}

	fmt.Printf("DEBUG: OTP for %s is %s\n", email, otp)
	return otp, nil
}

func (s *Service) VerifyOTP(email, otp string) error {
	ctx := context.Background()
	rdb := database.GetRedis()

	storedOTP, err := rdb.Get(ctx, "otp:"+email).Result()
	if err != nil {
		return errors.New("OTP expired or not found. please resend")
	}

	if storedOTP != otp {
		return errors.New("invalid OTP code")
	}

	err = s.repo.VerifyPhone(email)
	if err != nil {
		return err
	}

	u, _ := s.repo.GetByEmail(email)
	if u != nil {
		s.auditLog(u.ID, "PHONE_VERIFIED", "", "", nil)
	}

	rdb.Del(ctx, "otp:"+email)
	return nil
}

func (s *Service) Login(req user.LoginRequest, ip, userAgent string) (*user.TokenResponse, error) {
	valid, err := captcha.VerifyToken(req.CaptchaToken)
	if err != nil || !valid {
		return nil, errors.New("bot detection failed. please try again")
	}

	cleanEmail := sanitizer.Email(req.Email)

	u, err := s.repo.GetByEmail(cleanEmail)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if u.Status != "active" {
		return nil, errors.New("account is " + u.Status)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
	if err != nil {
		s.auditLog(u.ID, "LOGIN_FAILED", ip, userAgent, map[string]interface{}{"reason": "invalid_password"})
		return nil, errors.New("invalid email or password")
	}

	sessionID := uuid.New().String()
	accessToken, refreshToken, err := auth.GenerateTokens(u.ID, sessionID)
	if err != nil {
		return nil, err
	}

	s.storeSession(u.ID, sessionID, refreshToken, ip, userAgent)
	s.auditLog(u.ID, "LOGIN_SUCCESS", ip, userAgent, nil)
	s.ensureFullURL(u)

	return &user.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken handles the rotation of refresh tokens and provides replay protection.
// It verifies the old token, checks if it has been blacklisted (indicating reuse),
// and issues a new pair of tokens if the session is still active.
func (s *Service) RefreshToken(oldRefreshToken, ip, userAgent string) (*user.TokenResponse, error) {
	claims, err := auth.ValidateToken(oldRefreshToken)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	rdb := database.GetRedis()

	// 1. Check for Refresh Token Reuse (Replay Protection)
	// If the token is found in the blacklist, it means it has been used before.
	blacklisted, _ := rdb.Exists(ctx, "blacklist:"+oldRefreshToken).Result()
	if blacklisted > 0 {
		// DANGER: Refresh token reuse detected! This usually means the token was stolen.
		// As a security precaution, we revoke ALL active sessions for this user.
		s.RevokeAllSessions(claims.UserID)
		return nil, errors.New("security alert: session compromised. please login again")
	}

	// 2. Validate that the specific session still exists in Redis
	sessionKey := fmt.Sprintf("session:%s:%s", claims.UserID, claims.SessionID)
	exists, _ := rdb.Exists(ctx, sessionKey).Result()
	if exists == 0 {
		return nil, errors.New("session expired or revoked")
	}

	// 3. Generate new token pair (Rotation)
	accessToken, newRefreshToken, err := auth.GenerateTokens(claims.UserID, claims.SessionID)
	if err != nil {
		return nil, err
	}

	// 4. Update the session in Redis and blacklist the consumed refresh token
	s.storeSession(claims.UserID, claims.SessionID, newRefreshToken, ip, userAgent)
	
	// Blacklist the old token for 7 days (matching its max lifetime) to prevent replay
	rdb.Set(ctx, "blacklist:"+oldRefreshToken, "true", 7*24*time.Hour)

	return &user.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *Service) SocialLogin(socialUser *socialauth.SocialUser, ip, userAgent string) (*user.TokenResponse, error) {
	cleanEmail := sanitizer.Email(socialUser.Email)
	cleanFullName := sanitizer.Text(socialUser.FullName)

	u, err := s.repo.GetByEmail(cleanEmail)
	if err != nil {
		if err.Error() == "user not found" {
			now := time.Now()
			u = &user.User{
				Email:    cleanEmail,
				FullName: cleanFullName,
				Password: "SOCIAL_AUTH_NO_PASSWORD",
				Role:     "customer",
				Status:   "active",
				EmailVerifiedAt: &now,
			}
			err = s.repo.Create(u)
			if err != nil {
				return nil, err
			}
			s.auditLog(u.ID, "USER_REGISTERED_SOCIAL", ip, userAgent, nil)
		} else {
			return nil, err
		}
	} else if u.Status != "active" {
		return nil, errors.New("account is " + u.Status)
	}

	sessionID := uuid.New().String()
	accessToken, refreshToken, err := auth.GenerateTokens(u.ID, sessionID)
	if err != nil {
		return nil, err
	}

	s.storeSession(u.ID, sessionID, refreshToken, ip, userAgent)
	s.auditLog(u.ID, "LOGIN_SUCCESS_SOCIAL", ip, userAgent, nil)
	s.ensureFullURL(u)

	return &user.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) storeSession(userID, sessionID, refreshToken, ip, userAgent string) {
	ctx := context.Background()
	rdb := database.GetRedis()
	
	session := user.Session{
		ID:           sessionID,
		UserID:       userID,
		IPAddress:    ip,
		UserAgent:    userAgent,
		LastActiveAt: time.Now(),
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	// Dynamic Location Detection
	session.Location = s.GetLocationFromIP(ip)

	if strings.Contains(strings.ToLower(userAgent), "mobile") {
		session.DeviceType = "Mobile"
	} else {
		session.DeviceType = "Desktop"
	}

	data, _ := json.Marshal(session)
	sessionKey := fmt.Sprintf("session:%s:%s", userID, sessionID)
	rdb.Set(ctx, sessionKey, data, 7*24*time.Hour)
}

func (s *Service) GetLocationFromIP(ip string) string {
	// Clean IP from port if present (e.g., "172.18.0.1:60714")
	cleanIP := strings.Split(ip, ":")[0]

	// 1. Handle Localhost
	if cleanIP == "127.0.0.1" || cleanIP == "::1" || cleanIP == "" {
		return "Localhost (Dev)"
	}

	// 2. Handle Private/Internal IP Ranges (RFC 1918 & Docker)
	// 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
	if strings.HasPrefix(cleanIP, "10.") || 
	   strings.HasPrefix(cleanIP, "192.168.") ||
	   (strings.HasPrefix(cleanIP, "172.") && s.isPrivate172(cleanIP)) {
		return "Internal Network (Docker/VPN)"
	}

	// 3. Use ip-api.com for Public IPs
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=city,country", cleanIP)
	
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "Unknown Location"
	}
	defer resp.Body.Close()

	var result struct {
		City    string `json:"city"`
		Country string `json:"country"`
		Status  string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || result.Status == "fail" {
		return "Unknown Location"
	}

	if result.City != "" && result.Country != "" {
		return fmt.Sprintf("%s, %s", result.City, result.Country)
	}
	
	return "Unknown Location"
}

// Helper to check if 172.x IP is in the private range 172.16.0.0 - 172.31.255.255
func (s *Service) isPrivate172(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) < 2 {
		return false
	}
	secondOctet, _ := strconv.Atoi(parts[1])
	return secondOctet >= 16 && secondOctet <= 31
}

func (s *Service) GetActiveSessions(userID, currentSessionID string) ([]user.Session, error) {
	ctx := context.Background()
	rdb := database.GetRedis()
	
	pattern := fmt.Sprintf("session:%s:*", userID)
	keys, err := rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	sessions := make([]user.Session, 0)
	for _, key := range keys {
		data, err := rdb.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var sess user.Session
		if err := json.Unmarshal([]byte(data), &sess); err != nil {
			continue
		}

		if sess.ID == currentSessionID {
			sess.IsCurrent = true
		}
		sessions = append(sessions, sess)
	}
	return sessions, nil
}

func (s *Service) RevokeSession(userID, sessionID string) error {
	ctx := context.Background()
	rdb := database.GetRedis()
	sessionKey := fmt.Sprintf("session:%s:%s", userID, sessionID)
	return rdb.Del(ctx, sessionKey).Err()
}

func (s *Service) RevokeAllSessions(userID string) error {
	ctx := context.Background()
	rdb := database.GetRedis()
	pattern := fmt.Sprintf("session:%s:*", userID)
	keys, _ := rdb.Keys(ctx, pattern).Result()
	if len(keys) > 0 {
		return rdb.Del(ctx, keys...).Err()
	}
	return nil
}

func (s *Service) GetUserByID(id string) (*user.User, error) {
	u, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	s.ensureFullURL(u)
	return u, nil
}

// ensureFullURL memastikan avatar_url memiliki URL lengkap jika hanya tersimpan path relatifnya.
func (s *Service) ensureFullURL(u *user.User) {
	if u.AvatarURL != "" && !strings.HasPrefix(u.AvatarURL, "http") && s.uploader != nil {
		u.AvatarURL = fmt.Sprintf("%s/%s", s.uploader.GetBaseURL(), u.AvatarURL)
	}
}


func (s *Service) UpdateProfile(id string, req user.UpdateProfileRequest) error {
	req.FullName = sanitizer.Text(req.FullName)
	req.Username = sanitizer.Username(req.Username)
	req.PhoneNumber = sanitizer.Phone(req.PhoneNumber)
	req.Bio = sanitizer.Text(req.Bio)

	if req.FullName == "" {
		return errors.New("full name cannot be empty")
	}

	err := s.repo.UpdateProfile(id, &req)
	if err != nil {
		if strings.Contains(err.Error(), "users_username_key") {
			return errors.New("username is already taken")
		}
		if strings.Contains(err.Error(), "users_phone_number_key") {
			return errors.New("phone number is already taken")
		}
		return err
	}

	s.auditLog(id, "PROFILE_UPDATED", "", "", nil)
	return nil
}

func (s *Service) UploadAvatar(ctx context.Context, userID string, file io.Reader, contentType string, size int64) (string, error) {
	if s.uploader == nil {
		return "", errors.New("storage uploader is not configured")
	}

	if size > 2*1024*1024 { // 2MB Limit
		return "", errors.New("file size exceeds 2MB limit")
	}

	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
		return "", errors.New("invalid file type, only JPEG, PNG and WebP are allowed")
	}

	ext := "jpg"
	if contentType == "image/png" {
		ext = "png"
	} else if contentType == "image/webp" {
		ext = "webp"
	}

	objectName := fmt.Sprintf("avatars/%s-%s.%s", userID, uuid.New().String()[:8], ext)

	fileURL, err := s.uploader.UploadFile(ctx, file, objectName, contentType, size)
	if err != nil {
		return "", err
	}

	err = s.repo.UpdateAvatar(userID, fileURL)
	if err != nil {
		// Attempt to delete if db update fails
		_ = s.uploader.DeleteFile(ctx, objectName)
		return "", err
	}

	s.auditLog(userID, "AVATAR_UPDATED", "", "", nil)
	
	// Kembalikan URL lengkap untuk kebutuhan UI Frontend
	fullURL := fmt.Sprintf("%s/%s", s.uploader.GetBaseURL(), fileURL)
	return fullURL, nil
}

func (s *Service) Logout(token string) error {
	ctx := context.Background()
	rdb := database.GetRedis()
	return rdb.Set(ctx, "blacklist:"+token, "true", 24*time.Hour).Err()
}

func (s *Service) IsEmailAvailable(email string) (bool, error) {
	cleanEmail := sanitizer.Email(email)
	_, err := s.repo.GetByEmail(cleanEmail)
	if err != nil {
		if err.Error() == "user not found" {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (s *Service) ForgotPassword(identifier string) (string, string, error) {
	u, err := s.repo.GetByIdentifier(identifier)
	if err != nil {
		return "", "", errors.New("if account exists, reset instructions will be sent")
	}

	token := uuid.New().String()
	ctx := context.Background()
	rdb := database.GetRedis()
	
	err = rdb.Set(ctx, "reset_token:"+token, u.Email, 15*time.Minute).Err()
	if err != nil {
		return "", "", err
	}

	s.auditLog(u.ID, "PASSWORD_RESET_REQUESTED", "", "", map[string]interface{}{"identifier": identifier})

	parts := strings.Split(u.Email, "@")
	maskedEmail := u.Email
	if len(parts) == 2 {
		maskedName := string(parts[0][0]) + "***" + string(parts[0][len(parts[0])-1])
		maskedEmail = maskedName + "@" + parts[1][:1] + "****" + parts[1][strings.LastIndex(parts[1], "."):]
	}

	return token, maskedEmail, nil
}

func (s *Service) ValidateResetToken(token string) error {
	ctx := context.Background()
	rdb := database.GetRedis()
	
	exists, err := rdb.Exists(ctx, "reset_token:"+token).Result()
	if err != nil || exists == 0 {
		return errors.New("invalid or expired token")
	}

	return nil
}

func (s *Service) ResetPassword(token, newPassword string) error {
	ctx := context.Background()
	rdb := database.GetRedis()

	email, err := rdb.Get(ctx, "reset_token:"+token).Result()
	if err != nil {
		return errors.New("invalid or expired token")
	}

	u, err := s.repo.GetByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(newPassword))
	if err == nil {
		return errors.New("new password cannot be the same as the old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.repo.UpdatePassword(email, string(hashedPassword))
	if err != nil {
		return err
	}

	s.auditLog(u.ID, "PASSWORD_CHANGED_VIA_RESET", "", "", nil)
	rdb.Del(ctx, "reset_token:"+token)

	return nil
}

func (s *Service) LoginWithGoogleToken(idToken, ip, userAgent string) (*user.TokenResponse, error) {
	// 1. Verify token with Google API
	// Documentation: https://developers.google.com/identity/gsi/web/guides/verify-google-id-token
	verifyURL := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", idToken)
	
	resp, err := http.Get(verifyURL)
	if err != nil {
		return nil, fmt.Errorf("failed to verify google token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid google token")
	}

	var profile struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		Sub   string `json:"sub"` // Google unique ID
	}

	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("failed to decode google profile: %v", err)
	}

	// 2. Map to social user
	socialUser := &socialauth.SocialUser{
		ID:       profile.Sub,
		Email:    profile.Email,
		FullName: profile.Name,
	}

	// 3. Process social login (create if not exists, then return tokens)
	return s.SocialLogin(socialUser, ip, userAgent)
}
