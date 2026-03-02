package socialauth

import (
	"context"
	"encoding/json"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

type FacebookProvider struct {
	config *oauth2.Config
}

func NewFacebookProvider() *FacebookProvider {
	return &FacebookProvider{
		config: &oauth2.Config{
			ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
			ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("FACEBOOK_REDIRECT_URL"),
			Endpoint:     facebook.Endpoint,
			Scopes:       []string{"email", "public_profile"},
		},
	}
}

func (p *FacebookProvider) GetAuthURL(state string) string {
	return p.config.AuthCodeURL(state)
}

func (p *FacebookProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return p.config.Exchange(ctx, code)
}

func (p *FacebookProvider) GetUser(ctx context.Context, token *oauth2.Token) (*SocialUser, error) {
	client := p.config.Client(ctx, token)
	// Facebook Graph API
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email,picture")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var profile struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Picture struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	return &SocialUser{
		ID:       profile.ID,
		Email:    profile.Email,
		FullName: profile.Name,
		Avatar:   profile.Picture.Data.URL,
	}, nil
}
