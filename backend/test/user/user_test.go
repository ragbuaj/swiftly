package user_test

import (
	"context"
	"net/http/httptest"
	"testing"
	"swiftly/backend/internal/pkg/auth"
	"swiftly/backend/internal/user/service"
)

func TestAuthMiddlewareIntegration(t *testing.T) {
	token, _, err := auth.GenerateTokens("user-123", "session-456")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/users/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)

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

func TestGetLocationFromIP(t *testing.T) {
	s := service.NewService(nil, nil, nil)
	ctx := context.Background()

	tests := []struct {
		name     string
		ip       string
		expected string
	}{
		{
			name:     "Localhost IPv4",
			ip:       "127.0.0.1",
			expected: "Localhost (Dev)",
		},
		{
			name:     "Private IP (Docker)",
			ip:       "172.18.0.1",
			expected: "Internal Network (Docker/VPN)",
		},
		{
			name:     "Invalid IP Format",
			ip:       "not-an-ip",
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.GetLocationFromIP(ctx, tt.ip)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
