package user_test

import (
	"net/http/httptest"
	"strings"
	"testing"
	"swiftly/backend/internal/pkg/auth"
	"swiftly/backend/internal/user/service"
)

func TestAuthMiddlewareIntegration(t *testing.T) {
	// 1. Generate a valid token
	token, _, err := auth.GenerateTokens("user-123", "session-456")
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

func TestGetLocationFromIP(t *testing.T) {
	// Initialize service with nil repos as we only test a utility method
	s := service.NewService(nil, nil, nil)

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
			name:     "Localhost IPv6",
			ip:       "::1",
			expected: "Localhost (Dev)",
		},
		{
			name:     "Empty IP",
			ip:       "",
			expected: "Localhost (Dev)",
		},
		{
			name:     "Public IP (Google DNS)",
			ip:       "8.8.8.8",
			expected: "Ashburn, United States", // Expected result from ip-api.com for 8.8.8.8
		},
		{
			name:     "Private IP (Docker)",
			ip:       "172.18.0.1",
			expected: "Internal Network (Docker/VPN)",
		},
		{
			name:     "Invalid IP Format",
			ip:       "not-an-ip",
			expected: "Unknown Location",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: Public IP test might fail if there's no internet connection
			// or if the API provider changes their database result slightly.
			result := s.GetLocationFromIP(tt.ip)
			
			// For public IP, we check if it contains the country to be less brittle
			if tt.ip == "8.8.8.8" {
				if !strings.Contains(result, "United States") {
					t.Errorf("Expected location for 8.8.8.8 to contain 'United States', got: %s", result)
				}
			} else {
				if result != tt.expected {
					t.Errorf("Expected %s, got %s", tt.expected, result)
				}
			}
		})
	}
}
