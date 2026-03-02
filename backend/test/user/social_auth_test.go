package user_test

import (
	"context"
	"testing"

	"golang.org/x/oauth2"
	"swiftly/backend/internal/pkg/socialauth"
)

// MockProvider is a fake social auth provider for testing
type MockProvider struct {
	UserToReturn *socialauth.SocialUser
	ShouldFail   bool
}

func (m *MockProvider) GetAuthURL(state string) string {
	return "https://mock-auth.com/authorize?state=" + state
}

func (m *MockProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	if m.ShouldFail {
		return nil, context.DeadlineExceeded
	}
	return &oauth2.Token{AccessToken: "mock-token"}, nil
}

func (m *MockProvider) GetUser(ctx context.Context, token *oauth2.Token) (*socialauth.SocialUser, error) {
	if m.ShouldFail {
		return nil, context.DeadlineExceeded
	}
	return m.UserToReturn, nil
}

func TestSocialAuthLogic(t *testing.T) {
	mockUser := &socialauth.SocialUser{
		ID:       "12345",
		Email:    "social@example.com",
		FullName: "Social User",
	}

	provider := &MockProvider{UserToReturn: mockUser}

	t.Run("Get Auth URL", func(t *testing.T) {
		url := provider.GetAuthURL("test-state")
		expected := "https://mock-auth.com/authorize?state=test-state"
		if url != expected {
			t.Errorf("Expected URL %s, got %s", expected, url)
		}
	})

	t.Run("Successful Exchange and GetUser", func(t *testing.T) {
		token, err := provider.Exchange(context.Background(), "mock-code")
		if err != nil {
			t.Fatalf("Exchange failed: %v", err)
		}

		user, err := provider.GetUser(context.Background(), token)
		if err != nil {
			t.Fatalf("GetUser failed: %v", err)
		}

		if user.Email != mockUser.Email {
			t.Errorf("Expected email %s, got %s", mockUser.Email, user.Email)
		}
	})

	t.Run("Failed Provider", func(t *testing.T) {
		failingProvider := &MockProvider{ShouldFail: true}
		_, err := failingProvider.Exchange(context.Background(), "code")
		if err == nil {
			t.Error("Expected error for failing provider, got nil")
		}
	})
}

func TestSocialAuthRegistry(t *testing.T) {
	// Note: This tests the actual registry implementation
	registry := socialauth.NewRegistry()

	t.Run("Get Valid Provider", func(t *testing.T) {
		p, err := registry.GetProvider("google")
		if err != nil {
			t.Errorf("Expected google provider, got error: %v", err)
		}
		if p == nil {
			t.Error("Provider should not be nil")
		}
	})

	t.Run("Get Invalid Provider", func(t *testing.T) {
		_, err := registry.GetProvider("invalid-provider")
		if err == nil {
			t.Error("Expected error for invalid provider, got nil")
		}
	})
}
