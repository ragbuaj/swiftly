package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"swiftly/backend/internal/database"
	"swiftly/backend/internal/pkg/auth"
	"swiftly/backend/internal/pkg/captcha"
	"swiftly/backend/internal/pkg/sanitizer"
	"swiftly/backend/internal/pkg/socialauth"
	"swiftly/backend/internal/user"
	"swiftly/backend/internal/user/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo         *repository.Repository
	activityRepo *repository.ActivityRepository
}

func NewService(repo *repository.Repository, activityRepo *repository.ActivityRepository) *Service {
	return &Service{
		repo:         repo,
		activityRepo: activityRepo,
	}
}

func (s *Service) auditLog(userID, activityType, ip, userAgent string, metadata map[string]interface{}) {
	_ = s.activityRepo.Log(userID, activityType, ip, userAgent, metadata)
}

func (s *Service) CreateUser(req user.CreateUserRequest) (*user.TokenResponse, error) {
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
	}

	err = s.repo.Create(u)
	if err != nil {
		return nil, err
	}

	s.GenerateAndStoreOTP(cleanEmail)
	s.auditLog(u.ID, "USER_REGISTERED", "", "", nil)

	accessToken, refreshToken, err := auth.GenerateTokens(u.ID)
	if err != nil {
		return nil, err
	}

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

func (s *Service) Login(req user.LoginRequest) (*user.TokenResponse, error) {
	valid, err := captcha.VerifyToken(req.CaptchaToken)
	if err != nil || !valid {
		return nil, errors.New("bot detection failed. please try again")
	}

	cleanEmail := sanitizer.Email(req.Email)

	u, err := s.repo.GetByEmail(cleanEmail)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
	if err != nil {
		s.auditLog(u.ID, "LOGIN_FAILED", "", "", map[string]interface{}{"reason": "invalid_password"})
		return nil, errors.New("invalid email or password")
	}

	accessToken, refreshToken, err := auth.GenerateTokens(u.ID)
	if err != nil {
		return nil, err
	}

	s.auditLog(u.ID, "LOGIN_SUCCESS", "", "", nil)

	return &user.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshToken(refreshToken string) (*user.TokenResponse, error) {
	claims, err := auth.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	accessToken, newRefreshToken, err := auth.GenerateTokens(claims.UserID)
	if err != nil {
		return nil, err
	}

	return &user.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *Service) SocialLogin(socialUser *socialauth.SocialUser) (*user.TokenResponse, error) {
	cleanEmail := sanitizer.Email(socialUser.Email)
	cleanFullName := sanitizer.Text(socialUser.FullName)

	u, err := s.repo.GetByEmail(cleanEmail)
	if err != nil {
		if err.Error() == "user not found" {
			u = &user.User{
				Email:    cleanEmail,
				FullName: cleanFullName,
				Password: "SOCIAL_AUTH_NO_PASSWORD",
			}
			err = s.repo.Create(u)
			if err != nil {
				return nil, err
			}
			s.auditLog(u.ID, "USER_REGISTERED_SOCIAL", "", "", nil)
		} else {
			return nil, err
		}
	}

	accessToken, refreshToken, err := auth.GenerateTokens(u.ID)
	if err != nil {
		return nil, err
	}

	s.auditLog(u.ID, "LOGIN_SUCCESS_SOCIAL", "", "", nil)

	return &user.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) GetUserByID(id string) (*user.User, error) {
	return s.repo.GetByID(id)
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

func (s *Service) LoginWithGoogleToken(idToken string) (*user.TokenResponse, error) {
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
	return s.SocialLogin(socialUser)
}
