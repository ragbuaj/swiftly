package user

import "time"

type User struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Username    string    `json:"username,omitempty"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	FullName    string    `json:"full_name"`
	Password    string    `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Email        string `json:"email"`
	Username     string `json:"username"`
	PhoneNumber  string `json:"phone_number"`
	FullName     string `json:"full_name"`
	Password     string `json:"password"`
	CaptchaToken string `json:"captcha_token"`
}

type LoginRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	CaptchaToken string `json:"captcha_token"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
