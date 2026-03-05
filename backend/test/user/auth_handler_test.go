package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"swiftly/backend/internal/api"
	"swiftly/backend/internal/app"
	"swiftly/backend/internal/config"
	"swiftly/backend/internal/database"
)

func setupTestHandler(t *testing.T) http.Handler {
	// Gunakan kredensial yang terbukti berhasil di debug tadi
	os.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/swiftly?sslmode=disable")
	os.Setenv("JWT_SECRET", "test-secret-key-that-is-long-enough-12345")
	os.Setenv("REDIS_URL", "redis://localhost:6379/1")
	os.Setenv("APP_ENV", "test")

	cfg := config.Load()
	database.Init()
	database.InitRedis()

	pool := database.GetPool()
	if pool == nil { t.Fatal("Database pool is nil") }

	application := app.Init(cfg)
	return api.NewRouter(application)
}

func cleanupTestData() {
	pool := database.GetPool()
	if pool != nil {
		// Hapus semua data test agar tidak mengganggu test berikutnya
		_, _ = pool.Exec(context.Background(), "DELETE FROM users WHERE email LIKE 'testuser_%' OR username LIKE 'testuser_%'")
	}
	rdb := database.GetRedis()
	if rdb != nil { _ = rdb.FlushDB(context.Background()).Err() }
}

func TestArchitectureValidation(t *testing.T) {
	handler := setupTestHandler(t)
	defer database.Close()
	defer database.CloseRedis()
	defer cleanupTestData()

	t.Run("ARCH-01: Validation - Input Too Short", func(t *testing.T) {
		payload := map[string]string{"email": "valid@email.com", "password": "123", "username": "me"}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest { t.Errorf("Expected 400, got %v", rr.Code) }
	})

	t.Run("ARCH-02: Error Mapping - Invalid Login", func(t *testing.T) {
		payload := map[string]string{"email": "testuser_nonexistent@example.com", "password": "wrongpassword", "captcha_token": "dummy"}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusUnauthorized { t.Errorf("Expected 401, got %v", rr.Code) }
	})
}

func TestRegistration(t *testing.T) {
	handler := setupTestHandler(t)
	defer database.Close()
	defer database.CloseRedis()
	// Kita tidak panggil cleanupTestData() di awal agar data benar-benar bersih dari cycle sebelumnya
	cleanupTestData()
	defer cleanupTestData()

	t.Run("Successful Registration", func(t *testing.T) {
		ts := time.Now().UnixNano()
		email := fmt.Sprintf("testuser_%d@example.com", ts)
		username := fmt.Sprintf("testuser%d", ts % 100000) // Hindari karakter aneh, pakai numerik saja
		phone := fmt.Sprintf("08%d", ts % 100000000)

		payload := map[string]string{
			"email": email, "username": username, "phone_number": phone,
			"password": "Password123!", "full_name": "Test User", "captcha_token": "dummy",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		
		if rr.Code != http.StatusCreated { 
			t.Errorf("Expected 201, got %v. Body: %s", rr.Code, rr.Body.String()) 
		}
	})
}
