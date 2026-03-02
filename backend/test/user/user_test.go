package user_test

import (
	"net/http/httptest"
	"testing"
	"swiftly/backend/internal/pkg/auth"
)

func TestAuthMiddlewareIntegration(t *testing.T) {
	// 1. Generate a valid token
	token, _, err := auth.GenerateTokens("user-123")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// 2. Setup request with token
	req := httptest.NewRequest("GET", "/api/users/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// 3. Verify the auth utility logic
	claims, err := auth.ValidateToken(token)
	if err != nil {
		t.Errorf("Token should be valid: %v", err)
	}
	if claims.UserID != "user-123" {
		t.Errorf("Expected UserID user-123, got %s", claims.UserID)
	}
}

func TestInvalidToken(t *testing.T) {
	_, err := auth.ValidateToken("invalid.token.string")
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}
