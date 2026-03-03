package socialauth

import (
	"context"
	"encoding/json"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleProvider struct {
	config *oauth2.Config
}

func NewGoogleProvider() *GoogleProvider {
	return &GoogleProvider{
		config: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Endpoint:     google.Endpoint,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
		},
	}
}

func (p *GoogleProvider) GetAuthURL(state string) string {
	return p.config.AuthCodeURL(state)
}

func (p *GoogleProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return p.config.Exchange(ctx, code)
}

func (p *GoogleProvider) GetUser(ctx context.Context, token *oauth2.Token) (*SocialUser, error) {
	client := p.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var profile struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	return &SocialUser{
		ID:       profile.ID,
		Email:    profile.Email,
		FullName: profile.Name,
		Avatar:   profile.Picture,
	}, nil
}
