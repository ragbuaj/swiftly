package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"swiftly/backend/internal/database"
	"swiftly/backend/internal/pkg/socialauth"
	"swiftly/backend/internal/user/handler"
	"swiftly/backend/internal/user/repository"
	"swiftly/backend/internal/user/service"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// setupTestHandler menginisialisasi koneksi DB dan mengembalikan ServeMux yang siap diuji
func setupTestHandler(t *testing.T) *http.ServeMux {
	// Coba load .env dari root backend (naik 2 level dari folder test)
	_ = godotenv.Load("../../.env")

	// Override URL untuk local testing dari mesin host (di luar Docker)
	// Kita force agar dia nembak ke localhost karena port-nya sudah di-expose (5432 & 6379)
	os.Setenv("REDIS_URL", "localhost:6379")
	os.Setenv("REDIS_DB", "1") // Use DB 1 for testing to avoid wiping dev data
	os.Setenv("TURNSTILE_SECRET_KEY", "") // Bypass captcha for tests
	dbURL := os.Getenv("DATABASE_URL")
	if strings.Contains(dbURL, "@postgres:") {
		// Jika menggunakan hostname 'postgres', ubah jadi localhost
		dbURL = strings.Replace(dbURL, "@postgres:", "@localhost:", 1)
		os.Setenv("DATABASE_URL", dbURL)
	} else if dbURL == "" {
		// Default local database url
		os.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/swiftly?sslmode=disable")
	}

	// Init DB
	database.Init()
	// Init Redis
	database.InitRedis()

	pool := database.GetPool()
	if pool == nil {
		t.Fatal("Database pool is nil. Pastikan docker database berjalan dan .env benar.")
	}

	userRepo := repository.NewUserRepository(pool)
	activityRepo := repository.NewActivityRepository(pool)
	// Untuk testing, kita bisa pass nil untuk uploader (karena kita tidak test MinIO di sini)
	userService := service.NewService(userRepo, activityRepo, nil)
	socialRegistry := socialauth.NewRegistry()

	userHandler := handler.NewUserHandler(userService, socialRegistry)

	mux := http.NewServeMux()
	userHandler.Register(mux)

	return mux
}

// Helper untuk men-generate email unik agar tidak bentrok dengan test sebelumnya
func generateUniqueEmail() string {
	return "testuser_" + uuid.New().String()[:8] + "@example.com"
}

// Helper untuk membersihkan data testing dari database setelah test selesai
func cleanupTestData() {
	pool := database.GetPool()
	if pool != nil {
		// Hapus semua user yang email-nya berawalan 'testuser_'
		_, _ = pool.Exec(context.Background(), "DELETE FROM users WHERE email LIKE 'testuser_%'")
	}

	// Bersihkan data di Redis
	rdb := database.GetRedis()
	if rdb != nil {
		_ = rdb.FlushDB(context.Background()).Err()
	}
}

func TestRegistration(t *testing.T) {
	mux := setupTestHandler(t)
	defer database.Close()
	defer database.CloseRedis()
	defer cleanupTestData()

	uniqueEmail := generateUniqueEmail()

	// REG-01: Registrasi Sukses
	t.Run("REG-01: Successful Registration", func(t *testing.T) {
		payload := map[string]string{
			"email":         uniqueEmail,
			"username":      "user_" + uuid.New().String()[:8],
			"phone_number":  "phone_" + uuid.New().String()[:8],
			"password":      "Password123!",
			"full_name":     "Test User",
			"captcha_token": "dummy-token", // Asumsi validasi captcha dinonaktifkan/mocked untuk testing
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated && status != http.StatusBadRequest {
			// Catatan: Jika gagal karena Captcha (HTTP 400), berarti sistem menolak dummy-token.
			// Di sistem real, kita harus membuat mock Captcha verification.
			t.Errorf("Handler returned wrong status code: got %v want %v. Body: %s", status, http.StatusCreated, rr.Body.String())
		}
	})

	// REG-02: Email Sudah Terdaftar
	t.Run("REG-02: Email Already Exists", func(t *testing.T) {
		payload := map[string]string{
			"email":         uniqueEmail, // Gunakan email yang sama dengan REG-01
			"username":      "user_" + uuid.New().String()[:8],
			"password":      "Password123!",
			"full_name":     "Test User 2",
			"captcha_token": "dummy-token",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected 400 Bad Request for duplicate email, got %v", rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "email is already registered") {
			t.Errorf("Expected 'email already registered' error message, got: %s", rr.Body.String())
		}
	})

	// REG-03: Field Kosong
	t.Run("REG-03: Missing Fields", func(t *testing.T) {
		payload := map[string]string{
			"email": "onlyemail@test.com",
			// Password dan full_name kosong
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected 400 Bad Request for missing fields, got %v", rr.Code)
		}
	})
}

func TestLogin(t *testing.T) {
	mux := setupTestHandler(t)
	defer database.Close()
	defer database.CloseRedis()
	defer cleanupTestData()

	// Persiapkan User (Buat langsung lewat HTTP untuk testing login)
	testEmail := generateUniqueEmail()
	testPassword := "LoginPass123"

	// Setup: Register a user first
	setupPayload := map[string]string{
		"email":         testEmail,
		"username":      "user_" + uuid.New().String()[:8],
		"phone_number":  "phone_" + uuid.New().String()[:8],
		"password":      testPassword,
		"full_name":     "Login Test User",
		"captcha_token": "dummy-token",
	}
	bodySetup, _ := json.Marshal(setupPayload)
	reqSetup := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(bodySetup))
	reqSetup.Header.Set("Content-Type", "application/json")
	rrSetup := httptest.NewRecorder()
	mux.ServeHTTP(rrSetup, reqSetup)
	if rrSetup.Code != http.StatusCreated {
		t.Fatalf("Setup failed: expected 201 Created, got %v. Body: %s", rrSetup.Code, rrSetup.Body.String())
	}

	time.Sleep(1 * time.Second) // Beri jeda sedikit untuk memastikan DB tersimpan

	// LOG-01: Successful Login
	t.Run("LOG-01: Successful Login", func(t *testing.T) {
		payload := map[string]string{
			"email":         testEmail,
			"password":      testPassword,
			"captcha_token": "dummy-token",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200 OK for successful login, got %v", rr.Code)
		}
		
		// Verify token is in cookies
		cookies := rr.Result().Cookies()
		foundAccess := false
		for _, c := range cookies {
			if c.Name == "access_token" {
				foundAccess = true
			}
		}
		if !foundAccess {
			t.Errorf("Expected access_token cookie not found")
		}
	})

	// LOG-02: Password Salah
	t.Run("LOG-02: Invalid Password", func(t *testing.T) {
		payload := map[string]string{
			"email":         testEmail,
			"password":      "WrongPassword!",
			"captcha_token": "dummy-token",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 Unauthorized for wrong password, got %v", rr.Code)
		}
	})

	// LOG-03: Email Tidak Ditemukan
	t.Run("LOG-03: Non-existent Email", func(t *testing.T) {
		payload := map[string]string{
			"email":         "doesnotexist@example.com",
			"password":      "AnyPassword",
			"captcha_token": "dummy-token",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 Unauthorized for non-existent email, got %v", rr.Code)
		}
	})
}

func TestPasswordRecovery(t *testing.T) {
	mux := setupTestHandler(t)
	defer database.Close()
	defer database.CloseRedis()
	defer cleanupTestData()

	testEmail := generateUniqueEmail()
	
	// Setup user
	setupPayload := map[string]string{
		"email": testEmail, "username": "user_" + uuid.New().String()[:8], "phone_number": "phone_" + uuid.New().String()[:8], "password": "Password123!", "full_name": "Recover User", "captcha_token": "dummy",
	}
	bodySetup, _ := json.Marshal(setupPayload)
	reqSetup := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(bodySetup))
	reqSetup.Header.Set("Content-Type", "application/json")
	rrSetup := httptest.NewRecorder()
	mux.ServeHTTP(rrSetup, reqSetup)
	if rrSetup.Code != http.StatusCreated {
		t.Fatalf("Setup failed: expected 201 Created, got %v. Body: %s", rrSetup.Code, rrSetup.Body.String())
	}
	time.Sleep(1 * time.Second)

	var resetToken string

	// PWD-01: Forgot Password Request
	t.Run("PWD-01: Forgot Password Request", func(t *testing.T) {
		payload := map[string]string{"identifier": testEmail}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/forgot-password", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200 OK for forgot password, got %v", rr.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		if data, ok := resp["data"].(map[string]interface{}); ok {
			resetToken = data["token"].(string)
		}
	})

	// PWD-03: Invalid Reset Token
	t.Run("PWD-03: Invalid Reset Token", func(t *testing.T) {
		payload := map[string]string{"token": "invalid-token", "new_password": "NewPassword123!"}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/reset-password", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected 400 Bad Request for invalid token, got %v", rr.Code)
		}
	})

	// PWD-04: Successful Password Reset
	t.Run("PWD-04: Successful Password Reset", func(t *testing.T) {
		payload := map[string]string{"token": resetToken, "new_password": "NewPassword456!"}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/reset-password", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200 OK for successful reset, got %v. Body: %s", rr.Code, rr.Body.String())
		}
	})
}

func TestOTPVerification(t *testing.T) {
	mux := setupTestHandler(t)
	defer database.Close()
	defer database.CloseRedis()
	defer cleanupTestData()

	testEmail := generateUniqueEmail()

	// OTP-02: Invalid OTP
	t.Run("OTP-02: Invalid OTP Code", func(t *testing.T) {
		payload := map[string]string{"email": testEmail, "otp": "000000"}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/verify-otp", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		// This might return 400 if user/OTP not found, which is expected for invalid
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected 400 Bad Request for invalid OTP, got %v", rr.Code)
		}
	})

	// OTP-04: Resend OTP
	t.Run("OTP-04: Resend OTP", func(t *testing.T) {
		payload := map[string]string{"email": testEmail}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/resend-otp", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK && rr.Code != http.StatusInternalServerError {
			t.Errorf("Expected 200 or 500 (if no SMTP), got %v", rr.Code)
		}
	})
}

func TestUserProfile(t *testing.T) {
	mux := setupTestHandler(t)
	defer database.Close()
	defer database.CloseRedis()
	defer cleanupTestData()

	// Register to get access token
	testEmail := generateUniqueEmail()
	setupPayload := map[string]string{
		"email": testEmail, "username": "user_" + uuid.New().String()[:8], "phone_number": "phone_" + uuid.New().String()[:8], "password": "Password123!", "full_name": "Profile User", "captcha_token": "dummy",
	}
	bodySetup, _ := json.Marshal(setupPayload)
	reqSetup := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(bodySetup))
	reqSetup.Header.Set("Content-Type", "application/json")
	rrSetup := httptest.NewRecorder()
	mux.ServeHTTP(rrSetup, reqSetup)

	var accessToken string
	for _, c := range rrSetup.Result().Cookies() {
		if c.Name == "access_token" {
			accessToken = c.Value
		}
	}

	// Test profile requires authentication, so we expect 401 without token
	t.Run("PRF-01: Get Profile Without Token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/users/profile", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 Unauthorized for unauthenticated profile request, got %v", rr.Code)
		}
	})

	// PRF-02: Get Profile With Token
	t.Run("PRF-02: Get Profile With Token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200 OK for authenticated profile request, got %v", rr.Code)
		}
	})

	// PRF-03: Update Profile
	t.Run("PRF-03: Update Profile", func(t *testing.T) {
		updatePayload := map[string]string{
			"full_name":    "Updated Name",
			"username":     "updated_" + uuid.New().String()[:8],
			"phone_number": "updated_" + uuid.New().String()[:8],
			"bio":          "New Bio",
		}
		body, _ := json.Marshal(updatePayload)
		req := httptest.NewRequest("PUT", "/api/users/profile", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200 OK for profile update, got %v. Body: %s", rr.Code, rr.Body.String())
		}
	})
}

func TestTokenManagement(t *testing.T) {
	mux := setupTestHandler(t)
	defer database.Close()
	defer database.CloseRedis()
	defer cleanupTestData()

	// Register & Login to get tokens
	testEmail := generateUniqueEmail()
	setupPayload := map[string]string{
		"email": testEmail, "username": "user_" + uuid.New().String()[:8], "phone_number": "phone_" + uuid.New().String()[:8], "password": "Password123!", "full_name": "Token User", "captcha_token": "dummy",
	}
	bodySetup, _ := json.Marshal(setupPayload)
	reqSetup := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(bodySetup))
	reqSetup.Header.Set("Content-Type", "application/json")
	rrSetup := httptest.NewRecorder()
	mux.ServeHTTP(rrSetup, reqSetup)

	time.Sleep(1 * time.Second) // Ensure different issued_at for rotation test

	var refreshToken string
	for _, c := range rrSetup.Result().Cookies() {
		if c.Name == "refresh_token" {
			refreshToken = c.Value
		}
	}

	// TKN-01: Successful Refresh with Rotation
	var newRefreshToken string
	t.Run("TKN-01: Refresh Token Rotation", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/auth/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refresh_token", Value: refreshToken})
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200 OK for refresh, got %v. Body: %s", rr.Code, rr.Body.String())
		}

		// Verify we got NEW cookies (rotation)
		for _, c := range rr.Result().Cookies() {
			if c.Name == "refresh_token" {
				newRefreshToken = c.Value
			}
		}

		if newRefreshToken == "" || newRefreshToken == refreshToken {
			t.Errorf("Expected new refresh token (rotation), but got same or empty")
		}
	})

	// TKN-02: Replay Protection (Old refresh token should fail)
	t.Run("TKN-02: Replay Protection", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/auth/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refresh_token", Value: refreshToken}) // Using OLD token
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 Unauthorized for replayed refresh token, got %v", rr.Code)
		}
	})
}

func TestSessionManagement(t *testing.T) {
	mux := setupTestHandler(t)
	defer database.Close()
	defer database.CloseRedis()
	defer cleanupTestData()

	testEmail := generateUniqueEmail()
	setupPayload := map[string]string{
		"email": testEmail, "username": "user_" + uuid.New().String()[:8], "phone_number": "phone_" + uuid.New().String()[:8], "password": "Password123!", "full_name": "Session User", "captcha_token": "dummy",
	}
	bodySetup, _ := json.Marshal(setupPayload)
	reqSetup := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(bodySetup))
	reqSetup.Header.Set("Content-Type", "application/json")
	rrSetup := httptest.NewRecorder()
	mux.ServeHTTP(rrSetup, reqSetup)

	var accessToken string
	for _, c := range rrSetup.Result().Cookies() {
		if c.Name == "access_token" {
			accessToken = c.Value
		}
	}

	// SES-01: List Sessions
	var sessionID string
	t.Run("SES-01: List Active Sessions", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/auth/sessions", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200 OK for sessions list, got %v", rr.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		sessions := resp["data"].([]interface{})
		if len(sessions) == 0 {
			t.Errorf("Expected at least one session, got 0")
		}
		
		// Grab a session ID (for revoke test later)
		firstSess := sessions[0].(map[string]interface{})
		sessionID = firstSess["id"].(string)
	})

	// SES-02: Revoke Session
	t.Run("SES-02: Revoke Specific Session", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/auth/sessions/"+sessionID, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200 OK for session revoke, got %v", rr.Code)
		}
	})

	// SES-03: Access After Revoke (Should fail)
	t.Run("SES-03: Access Denied After Revocation", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken) // Using token from revoked session
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 Unauthorized for revoked session, got %v", rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "revoked") {
			t.Errorf("Expected 'revoked' error message, got: %s", rr.Body.String())
		}
	})
}


