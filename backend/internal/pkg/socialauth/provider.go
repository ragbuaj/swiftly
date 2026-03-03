package socialauth

import (
	"context"
	"golang.org/x/oauth2"
)

type SocialUser struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Avatar   string `json:"avatar"`
}

type Provider interface {
	GetAuthURL(state string) string
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	GetUser(ctx context.Context, token *oauth2.Token) (*SocialUser, error)
}
