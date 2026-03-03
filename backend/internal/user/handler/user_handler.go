package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"swiftly/backend/internal/middleware"
	"swiftly/backend/internal/pkg/response"
	"swiftly/backend/internal/pkg/socialauth"
	"swiftly/backend/internal/user"
	"swiftly/backend/internal/user/service"
)

type UserHandler struct {
	service        *service.Service
	socialRegistry *socialauth.Registry
}

func NewUserHandler(s *service.Service, r *socialauth.Registry) *UserHandler {
	return &UserHandler{
		service:        s,
		socialRegistry: r,
	}
}

func (h *UserHandler) Register(mux *http.ServeMux) {
	// Auth Core
	mux.HandleFunc("POST /api/auth/register", h.CreateUser)
	mux.HandleFunc("POST /api/auth/login", h.Login)
	mux.HandleFunc("POST /api/auth/refresh", h.RefreshToken)
	mux.HandleFunc("POST /api/auth/logout", h.Logout)
	
	// Password Recovery (MUST be before social wildcard to avoid conflict)
	mux.HandleFunc("POST /api/auth/forgot-password", h.ForgotPassword)
	mux.HandleFunc("POST /api/auth/reset-password", h.ResetPassword)
	mux.HandleFunc("GET /api/auth/validate-reset-token", h.ValidateResetToken)

	// Phone Verification
	mux.HandleFunc("POST /api/auth/verify-otp", h.VerifyOTP)
	mux.HandleFunc("POST /api/auth/resend-otp", h.ResendOTP)
	
	// Utilities
	mux.Handle("GET /api/users/check-email", middleware.RateLimitMiddleware(3, 5)(http.HandlerFunc(h.CheckEmail)))
	
	// Social Auth
	mux.HandleFunc("POST /api/auth/google/token", h.GoogleTokenLogin)
	mux.HandleFunc("GET /api/auth/{provider}/callback", h.SocialCallback)
	mux.HandleFunc("GET /api/auth/{provider}", h.SocialLogin)
	
	// Protected routes
	mux.Handle("GET /api/users/profile", middleware.AuthMiddleware(http.HandlerFunc(h.GetUser)))
	mux.Handle("PUT /api/users/profile", middleware.AuthMiddleware(http.HandlerFunc(h.UpdateProfile)))
	mux.Handle("POST /api/users/profile/avatar", middleware.AuthMiddleware(http.HandlerFunc(h.UploadAvatar)))
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req user.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	err := h.service.UpdateProfile(userID, req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(w, http.StatusOK, "Profile updated successfully", nil)
}

func (h *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	// Limit request size to 5MB to prevent memory exhaustion
	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)

	// Parse multipart form
	err := r.ParseMultipartForm(5 << 20) 
	if err != nil {
		response.Error(w, http.StatusBadRequest, "File too large. Maximum size is 5MB", nil)
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Failed to get file from request", nil)
		return
	}
	defer file.Close()

	// Get file content type
	contentType := header.Header.Get("Content-Type")
	size := header.Size

	fileURL, err := h.service.UploadAvatar(r.Context(), userID, file, contentType, size)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(w, http.StatusOK, "Avatar uploaded successfully", map[string]string{
		"avatar_url": fileURL,
	})
}


func (h *UserHandler) ValidateResetToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		response.Error(w, http.StatusBadRequest, "Token is required", nil)
		return
	}

	err := h.service.ValidateResetToken(token)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(w, http.StatusOK, "Token is valid", nil)
}

func (h *UserHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Identifier string `json:"identifier"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	token, email, err := h.service.ForgotPassword(req.Identifier)
	if err != nil {
		response.Success(w, http.StatusOK, "Instructions sent if account exists", nil)
		return
	}

	response.Success(w, http.StatusOK, "Instructions sent", map[string]string{
		"token":      token,
		"email_hint": email,
	})
}

func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	err := h.service.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(w, http.StatusOK, "Password reset successful", nil)
}

func (h *UserHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	err := h.service.VerifyOTP(req.Email, req.OTP)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(w, http.StatusOK, "Phone number verified successfully", nil)
}

func (h *UserHandler) ResendOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	otp, err := h.service.GenerateAndStoreOTP(req.Email)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to send OTP", nil)
		return
	}

	response.Success(w, http.StatusOK, "OTP sent successfully", map[string]string{
		"otp": otp,
	})
}

func (h *UserHandler) CheckEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		response.Error(w, http.StatusBadRequest, "Email query parameter is required", nil)
		return
	}

	available, err := h.service.IsEmailAvailable(email)
	if err != nil {
		fmt.Printf("Error checking email availability: %v\n", err)
		response.Error(w, http.StatusInternalServerError, "Failed to check email availability", nil)
		return
	}

	response.Success(w, http.StatusOK, "Email availability checked", map[string]bool{
		"available": available,
	})
}

func (h *UserHandler) SocialLogin(w http.ResponseWriter, r *http.Request) {
	providerName := r.PathValue("provider")
	provider, err := h.socialRegistry.GetProvider(providerName)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	state := "random-state" 
	url := provider.GetAuthURL(state)
	
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *UserHandler) SocialCallback(w http.ResponseWriter, r *http.Request) {
	providerName := r.PathValue("provider")
	provider, _ := h.socialRegistry.GetProvider(providerName)
	code := r.URL.Query().Get("code")
	token, _ := provider.Exchange(r.Context(), code)
	socialUser, _ := provider.GetUser(r.Context(), token)
	tokens, _ := h.service.SocialLogin(socialUser)

	h.setAuthCookies(w, tokens)
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" { frontendURL = "http://localhost:5173" }
	http.Redirect(w, r, frontendURL, http.StatusTemporaryRedirect)
}

func (h *UserHandler) setAuthCookies(w http.ResponseWriter, tokens *user.TokenResponse) {
	http.SetCookie(w, &http.Cookie{
		Name: "access_token", Value: tokens.AccessToken, Path: "/", HttpOnly: true, Secure: false, SameSite: http.SameSiteLaxMode, MaxAge: 3600 * 24,
	})
	http.SetCookie(w, &http.Cookie{
		Name: "refresh_token", Value: tokens.RefreshToken, Path: "/", HttpOnly: true, Secure: false, SameSite: http.SameSiteLaxMode, MaxAge: 3600 * 24 * 7,
	})
}

func (h *UserHandler) clearAuthCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{ Name: "access_token", Value: "", Path: "/", HttpOnly: true, MaxAge: -1 })
	http.SetCookie(w, &http.Cookie{ Name: "refresh_token", Value: "", Path: "/", HttpOnly: true, MaxAge: -1 })
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req user.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	tokens, err := h.service.CreateUser(req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.setAuthCookies(w, tokens)
	response.Success(w, http.StatusCreated, "User created successfully", nil)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req user.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	tokens, err := h.service.Login(req)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	h.setAuthCookies(w, tokens)
	response.Success(w, http.StatusOK, "Login successful", nil)
}

func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshStr := ""
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		refreshStr = cookie.Value
	} else {
		var req struct{ RefreshToken string `json:"refresh_token"` }
		json.NewDecoder(r.Body).Decode(&req)
		refreshStr = req.RefreshToken
	}

	if refreshStr == "" {
		response.Error(w, http.StatusUnauthorized, "Refresh token required", nil)
		return
	}

	tokens, err := h.service.RefreshToken(refreshStr)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Invalid refresh token", nil)
		return
	}

	h.setAuthCookies(w, tokens)
	response.Success(w, http.StatusOK, "Token refreshed successfully", nil)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := ""
	cookie, err := r.Cookie("access_token")
	if err == nil { token = cookie.Value }
	if token != "" { h.service.Logout(token) }
	h.clearAuthCookies(w)
	response.Success(w, http.StatusOK, "Logout successful", nil)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	u, err := h.service.GetUserByID(userID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "User not found", nil)
		return
	}
	response.Success(w, http.StatusOK, "User retrieved successfully", u)
}

func (h *UserHandler) GoogleTokenLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IDToken string `json:"id_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// For now, we will use the existing SocialLogin logic in service
	// Note: In a real app, you should validate the id_token using Google's library first.
	// But let's assume we are calling service to handle it.
	
	// Since we need to validate IDToken, we need to implement it in Service
	tokens, err := h.service.LoginWithGoogleToken(req.IDToken)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	h.setAuthCookies(w, tokens)
	response.Success(w, http.StatusOK, "Login successful", nil)
}
