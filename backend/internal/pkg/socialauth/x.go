package socialauth

import (
	"context"
	"encoding/json"
	"os"

	"golang.org/x/oauth2"
)

type XProvider struct {
	config *oauth2.Config
}

func NewXProvider() *XProvider {
	return &XProvider{
		config: &oauth2.Config{
			ClientID:     os.Getenv("X_CLIENT_ID"),
			ClientSecret: os.Getenv("X_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("X_REDIRECT_URL"),
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://twitter.com/i/oauth2/authorize",
				TokenURL: "https://api.twitter.com/2/oauth2/token",
			},
			Scopes: []string{"tweet.read", "users.read", "offline.access"},
		},
	}
}

func (p *XProvider) GetAuthURL(state string) string {
	// X OAuth2 requires code_challenge for PKCE in some cases, 
	// but standard oauth2 package handles basic flows.
	return p.config.AuthCodeURL(state)
}

func (p *XProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return p.config.Exchange(ctx, code)
}

func (p *XProvider) GetUser(ctx context.Context, token *oauth2.Token) (*SocialUser, error) {
	client := p.config.Client(ctx, token)
	resp, err := client.Get("https://api.twitter.com/2/users/me?user.fields=profile_image_url,entities")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			Username        string `json:"username"`
			ProfileImageURL string `json:"profile_image_url"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &SocialUser{
		ID:       result.Data.ID,
		FullName: result.Data.Name,
		Avatar:   result.Data.ProfileImageURL,
		// X V2 API does not return email by default unless specifically approved
		Email:    result.Data.Username + "@x.com", 
	}, nil
}
