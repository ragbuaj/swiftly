package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"swiftly/backend/internal/database"
	"swiftly/backend/internal/pkg/apperror"
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

func (s *Service) auditLog(ctx context.Context, userID, activityType, ip, userAgent string, metadata map[string]interface{}) {
	_ = s.activityRepo.Log(ctx, userID, activityType, ip, userAgent, metadata)
}

func (s *Service) CreateUser(ctx context.Context, req user.CreateUserRequest, ip, userAgent string) (*user.TokenResponse, error) {
	valid, err := captcha.VerifyToken(req.CaptchaToken)
	if err != nil || !valid {
		return nil, apperror.New(http.StatusBadRequest, "CAPTCHA_FAILED", "Bot detection failed. Please try again")
	}

	cleanEmail := sanitizer.Email(req.Email)
	cleanFullName := sanitizer.Text(req.FullName)
	cleanUsername := sanitizer.Username(req.Username)
	cleanPhone := sanitizer.Phone(req.PhoneNumber)

	exists, _ := s.IsEmailAvailable(ctx, cleanEmail)
	if !exists {
		return nil, apperror.ErrConflict
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.ErrInternalServer
	}

	u := &user.User{
		Email:       cleanEmail,
		Username:    &cleanUsername,
		PhoneNumber: &cleanPhone,
		FullName:    cleanFullName,
		Password:    string(hashedPassword),
		Role:        user.RoleUser,
		Status:      user.StatusActive,
	}

	err = s.repo.Create(ctx, u)
	if err != nil {
		return nil, apperror.ErrInternalServer
	}

	s.GenerateAndStoreOTP(ctx, cleanEmail)
	s.auditLog(ctx, u.ID, "USER_REGISTERED", ip, userAgent, nil)

	sessionID := uuid.New().String()
	accessToken, refreshToken, err := auth.GenerateTokens(u.ID, sessionID)
	if err != nil {
		return nil, apperror.ErrInternalServer
	}

	s.storeSession(ctx, u.ID, sessionID, refreshToken, ip, userAgent)

	return &user.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Login(ctx context.Context, req user.LoginRequest, ip, userAgent string) (*user.TokenResponse, error) {
	valid, err := captcha.VerifyToken(req.CaptchaToken)
	if err != nil || !valid {
		return nil, apperror.New(http.StatusBadRequest, "CAPTCHA_FAILED", "Bot detection failed. Please try again")
	}

	cleanEmail := sanitizer.Email(req.Email)

	u, err := s.repo.GetByEmail(ctx, cleanEmail)
	if err != nil {
		return nil, apperror.New(http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
	}

	if u.Status != user.StatusActive {
		return nil, apperror.New(http.StatusForbidden, "ACCOUNT_SUSPENDED", "Account is "+u.Status)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
	if err != nil {
		s.auditLog(ctx, u.ID, "LOGIN_FAILED", ip, userAgent, map[string]interface{}{"reason": "invalid_password"})
		return nil, apperror.New(http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
	}

	sessionID := uuid.New().String()
	accessToken, refreshToken, err := auth.GenerateTokens(u.ID, sessionID)
	if err != nil {
		return nil, apperror.ErrInternalServer
	}

	s.storeSession(ctx, u.ID, sessionID, refreshToken, ip, userAgent)
	s.auditLog(ctx, u.ID, "LOGIN_SUCCESS", ip, userAgent, nil)

	return &user.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, oldRefreshToken, ip, userAgent string) (*user.TokenResponse, error) {
	claims, err := auth.ValidateToken(oldRefreshToken)
	if err != nil {
		return nil, apperror.ErrUnauthorized
	}

	rdb := database.GetRedis()
	blacklisted, _ := rdb.Exists(ctx, "blacklist:"+oldRefreshToken).Result()
	if blacklisted > 0 {
		s.RevokeAllSessions(ctx, claims.UserID)
		return nil, apperror.New(http.StatusUnauthorized, "SECURITY_ALERT", "Session compromised. Please login again")
	}

	sessionKey := fmt.Sprintf("session:%s:%s", claims.UserID, claims.SessionID)
	exists, _ := rdb.Exists(ctx, sessionKey).Result()
	if exists == 0 {
		return nil, apperror.New(http.StatusUnauthorized, "SESSION_EXPIRED", "Session expired or revoked")
	}

	accessToken, newRefreshToken, err := auth.GenerateTokens(claims.UserID, claims.SessionID)
	if err != nil {
		return nil, apperror.ErrInternalServer
	}

	s.storeSession(ctx, claims.UserID, claims.SessionID, newRefreshToken, ip, userAgent)
	rdb.Set(ctx, "blacklist:"+oldRefreshToken, "true", 7*24*time.Hour)

	return &user.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *Service) UpdateProfile(ctx context.Context, id string, req user.UpdateProfileRequest) error {
	req.FullName = sanitizer.Text(req.FullName)
	req.Username = sanitizer.Username(req.Username)
	req.PhoneNumber = sanitizer.Phone(req.PhoneNumber)
	req.Bio = sanitizer.Text(req.Bio)

	err := s.repo.UpdateProfile(ctx, id, &req)
	if err != nil {
		if strings.Contains(err.Error(), "users_username_key") {
			return apperror.New(http.StatusConflict, "USERNAME_TAKEN", "Username is already taken")
		}
		if strings.Contains(err.Error(), "users_phone_number_key") {
			return apperror.New(http.StatusConflict, "PHONE_TAKEN", "Phone number is already taken")
		}
		return apperror.ErrInternalServer
	}

	s.auditLog(ctx, id, "PROFILE_UPDATED", "", "", nil)
	return nil
}

func (s *Service) UploadAvatar(ctx context.Context, userID string, file io.Reader, contentType string, size int64) (string, error) {
	if s.uploader == nil {
		return "", apperror.New(http.StatusInternalServerError, "STORAGE_ERROR", "Storage uploader is not configured")
	}

	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
		return "", apperror.New(http.StatusBadRequest, "INVALID_FILE_TYPE", "Only JPEG, PNG and WebP are allowed")
	}

	ext := "jpg"
	if contentType == "image/png" { ext = "png" } else if contentType == "image/webp" { ext = "webp" }

	objectName := fmt.Sprintf("avatars/%s-%s.%s", userID, uuid.New().String()[:8], ext)
	fileURL, err := s.uploader.UploadFile(ctx, file, objectName, contentType, size)
	if err != nil {
		return "", apperror.ErrInternalServer
	}

	err = s.repo.UpdateAvatar(ctx, userID, fileURL)
	if err != nil {
		_ = s.uploader.DeleteFile(ctx, objectName)
		return "", apperror.ErrInternalServer
	}

	s.auditLog(ctx, userID, "AVATAR_UPDATED", "", "", nil)
	return fmt.Sprintf("%s/%s", s.uploader.GetBaseURL(), fileURL), nil
}

func (s *Service) SocialLogin(ctx context.Context, socialUser *socialauth.SocialUser, ip, userAgent string) (*user.TokenResponse, error) {
	cleanEmail := sanitizer.Email(socialUser.Email)
	cleanFullName := sanitizer.Text(socialUser.FullName)

	var u *user.User
	var err error

	u, err = s.repo.GetByEmail(ctx, cleanEmail)
	if err != nil {
		if err.Error() == "user not found" {
			now := time.Now()
			
			baseUsername := strings.Split(cleanEmail, "@")[0]
			if len(baseUsername) < 3 { baseUsername = "user" }
			uniqueUsername := fmt.Sprintf("%s%d", baseUsername, time.Now().Unix()%1000)

			newUser := &user.User{
				Email:    cleanEmail,
				Username: &uniqueUsername,
				FullName: cleanFullName,
				Password: "SOCIAL_AUTH_NO_PASSWORD",
				Role:     user.RoleUser,
				Status:   user.StatusActive,
				EmailVerifiedAt: &now,
			}
			err = s.repo.Create(ctx, newUser)
			if err != nil { return nil, apperror.ErrInternalServer }
			
			u = newUser
			s.auditLog(ctx, u.ID, "USER_REGISTERED_SOCIAL", ip, userAgent, nil)
		} else {
			return nil, apperror.ErrInternalServer
		}
	}

	if u == nil { return nil, apperror.ErrInternalServer }

	if u.Status != user.StatusActive {
		return nil, apperror.New(http.StatusForbidden, "ACCOUNT_SUSPENDED", "Account is "+u.Status)
	}

	sessionID := uuid.New().String()
	accessToken, refreshToken, err := auth.GenerateTokens(u.ID, sessionID)
	if err != nil { return nil, apperror.ErrInternalServer }

	s.storeSession(ctx, u.ID, sessionID, refreshToken, ip, userAgent)
	s.auditLog(ctx, u.ID, "LOGIN_SUCCESS_SOCIAL", ip, userAgent, nil)
	s.ensureFullURL(u)

	return &user.TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *Service) ForgotPassword(ctx context.Context, identifier string) (string, string, error) {
	u, err := s.repo.GetByIdentifier(ctx, identifier)
	if err != nil { return "", "", apperror.ErrNotFound }

	token := uuid.New().String()
	rdb := database.GetRedis()
	rdb.Set(ctx, "reset_token:"+token, u.Email, 15*time.Minute)

	maskedEmail := u.Email
	parts := strings.Split(u.Email, "@")
	if len(parts) == 2 { maskedEmail = string(parts[0][0]) + "***@" + parts[1] }

	return token, maskedEmail, nil
}

func (s *Service) ValidateResetToken(ctx context.Context, token string) error {
	rdb := database.GetRedis()
	exists, err := rdb.Exists(ctx, "reset_token:"+token).Result()
	if err != nil || exists == 0 { return apperror.New(http.StatusBadRequest, "INVALID_TOKEN", "Invalid or expired token") }
	return nil
}

func (s *Service) ResetPassword(ctx context.Context, token, newPassword string) error {
	rdb := database.GetRedis()
	email, err := rdb.Get(ctx, "reset_token:"+token).Result()
	if err != nil { return apperror.New(http.StatusBadRequest, "INVALID_TOKEN", "Invalid or expired token") }

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	err = s.repo.UpdatePassword(ctx, email, string(hashedPassword))
	if err != nil { return apperror.ErrInternalServer }

	rdb.Del(ctx, "reset_token:"+token)
	return nil
}

func (s *Service) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		fmt.Printf("DEBUG: GetUserByID failed for ID %s: %v\n", id, err)
		return nil, apperror.ErrNotFound
	}
	s.ensureFullURL(u)
	return u, nil
}

func (s *Service) VerifyOTP(ctx context.Context, email, otp string) error {
	rdb := database.GetRedis()
	storedOTP, err := rdb.Get(ctx, "otp:"+email).Result()
	if err != nil || storedOTP != otp { return apperror.New(http.StatusBadRequest, "INVALID_OTP", "Invalid or expired OTP") }

	err = s.repo.VerifyPhone(ctx, email)
	if err != nil { return apperror.ErrInternalServer }

	rdb.Del(ctx, "otp:"+email)
	return nil
}

func (s *Service) GenerateAndStoreOTP(ctx context.Context, email string) (string, error) {
	b := make([]byte, 3)
	rand.Read(b)
	otp := fmt.Sprintf("%06d", (uint32(b[0])<<16|uint32(b[1])<<8|uint32(b[2]))%1000000)
	rdb := database.GetRedis()
	rdb.Set(ctx, "otp:"+email, otp, 5*time.Minute)
	return otp, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	rdb := database.GetRedis()
	return rdb.Set(ctx, "blacklist:"+token, "true", 24*time.Hour).Err()
}

func (s *Service) RevokeSession(ctx context.Context, userID, sessionID string) error {
	rdb := database.GetRedis()
	return rdb.Del(ctx, fmt.Sprintf("session:%s:%s", userID, sessionID)).Err()
}

func (s *Service) RevokeAllSessions(ctx context.Context, userID string) error {
	rdb := database.GetRedis()
	keys, _ := rdb.Keys(ctx, fmt.Sprintf("session:%s:*", userID)).Result()
	if len(keys) > 0 { return rdb.Del(ctx, keys...).Err() }
	return nil
}

func (s *Service) GetActiveSessions(ctx context.Context, userID, currentSessionID string) ([]user.Session, error) {
	rdb := database.GetRedis()
	keys, _ := rdb.Keys(ctx, fmt.Sprintf("session:%s:*", userID)).Result()
	sessions := make([]user.Session, 0)
	for _, key := range keys {
		data, _ := rdb.Get(ctx, key).Result()
		var sess user.Session
		json.Unmarshal([]byte(data), &sess)
		if sess.ID == currentSessionID { sess.IsCurrent = true }
		sessions = append(sessions, sess)
	}
	return sessions, nil
}

func (s *Service) IsEmailAvailable(ctx context.Context, email string) (bool, error) {
	_, err := s.repo.GetByEmail(ctx, sanitizer.Email(email))
	if err != nil && err.Error() == "user not found" { return true, nil }
	return false, nil
}

func (s *Service) GetLocationFromIP(ctx context.Context, ip string) string {
	cleanIP := strings.Split(ip, ":")[0]
	if cleanIP == "127.0.0.1" || cleanIP == "::1" || cleanIP == "" { return "Localhost (Dev)" }
	if strings.HasPrefix(cleanIP, "10.") || strings.HasPrefix(cleanIP, "192.168.") || (strings.HasPrefix(cleanIP, "172.") && s.isPrivate172(cleanIP)) {
		return "Internal Network (Docker/VPN)"
	}
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=city,country", cleanIP)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := (&http.Client{Timeout: 2 * time.Second}).Do(req)
	if err != nil { return "Unknown" }
	defer resp.Body.Close()
	var result struct{ City, Country string }
	json.NewDecoder(resp.Body).Decode(&result)
	if result.City == "" { return "Unknown" }
	return fmt.Sprintf("%s, %s", result.City, result.Country)
}

func (s *Service) isPrivate172(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) < 2 { return false }
	val, _ := strconv.Atoi(parts[1])
	return val >= 16 && val <= 31
}

func (s *Service) ensureFullURL(u *user.User) {
	if u.AvatarURL != nil && *u.AvatarURL != "" && !strings.HasPrefix(*u.AvatarURL, "http") && s.uploader != nil {
		fullURL := fmt.Sprintf("%s/%s", s.uploader.GetBaseURL(), *u.AvatarURL)
		u.AvatarURL = &fullURL
	}
}

func (s *Service) storeSession(ctx context.Context, userID, sessionID, refreshToken, ip, userAgent string) {
	rdb := database.GetRedis()
	session := user.Session{
		ID: sessionID, UserID: userID, IPAddress: ip, UserAgent: userAgent,
		LastActiveAt: time.Now(), ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	session.Location = s.GetLocationFromIP(ctx, ip)
	if strings.Contains(strings.ToLower(userAgent), "mobile") { session.DeviceType = "Mobile" } else { session.DeviceType = "Desktop" }
	data, _ := json.Marshal(session)
	rdb.Set(ctx, fmt.Sprintf("session:%s:%s", userID, sessionID), data, 7*24*time.Hour)
}

func (s *Service) LoginWithGoogleToken(ctx context.Context, idToken, ip, userAgent string) (*user.TokenResponse, error) {
	verifyURL := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", idToken)
	resp, err := http.Get(verifyURL)
	if err != nil || resp.StatusCode != http.StatusOK { return nil, apperror.New(http.StatusUnauthorized, "GOOGLE_AUTH_FAILED", "Invalid google token") }
	defer resp.Body.Close()
	var profile struct{ Email, Name, Sub string }
	json.NewDecoder(resp.Body).Decode(&profile)
	return s.SocialLogin(ctx, &socialauth.SocialUser{ID: profile.Sub, Email: profile.Email, FullName: profile.Name}, ip, userAgent)
}
