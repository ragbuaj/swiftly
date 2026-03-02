package captcha

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type TurnstileResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

func VerifyToken(token string) (bool, error) {
	secretKey := os.Getenv("TURNSTILE_SECRET_KEY")
	if secretKey == "" {
		return true, nil // Skip validation if key is not set (dev mode)
	}

	siteVerifyURL := "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	
	formData := url.Values{
		"secret":   {secretKey},
		"response": {token},
	}

	resp, err := http.PostForm(siteVerifyURL, formData)
	if err != nil {
		return false, fmt.Errorf("failed to call turnstile API: %v", err)
	}
	defer resp.Body.Close()

	var turnstileResp TurnstileResponse
	if err := json.NewDecoder(resp.Body).Decode(&turnstileResp); err != nil {
		return false, fmt.Errorf("failed to decode turnstile response: %v", err)
	}

	return turnstileResp.Success, nil
}
