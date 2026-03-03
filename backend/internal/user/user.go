package user

import "time"

type User struct {
	ID                       string     `json:"id"`
	Email                    string     `json:"email"`
	Username                 string     `json:"username,omitempty"`
	PhoneNumber              string     `json:"phone_number,omitempty"`
	FullName                 string     `json:"full_name"`
	Password                 string     `json:"-"`
	Role                     string     `json:"role"`
	Status                   string     `json:"status"`
	EmailVerifiedAt          *time.Time `json:"email_verified_at,omitempty"`
	DeletedAt                *time.Time `json:"deleted_at,omitempty"`
	DateOfBirth              *time.Time `json:"date_of_birth,omitempty"`
	Gender                   string     `json:"gender,omitempty"`
	NewsletterSubscribed     bool       `json:"newsletter_subscribed"`
	AvatarURL                string     `json:"avatar_url,omitempty"`
	Bio                      string     `json:"bio,omitempty"`
	DefaultShippingAddressID string     `json:"default_shipping_address_id,omitempty"`
	DefaultBillingAddressID  string     `json:"default_billing_address_id,omitempty"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
}

type CreateUserRequest struct {
	Email        string `json:"email"`
	Username     string `json:"username"`
	PhoneNumber  string `json:"phone_number"`
	FullName     string `json:"full_name"`
	Password     string `json:"password"`
	CaptchaToken string `json:"captcha_token"`
}

type UpdateProfileRequest struct {
	FullName             string     `json:"full_name"`
	Username             string     `json:"username"`
	PhoneNumber          string     `json:"phone_number"`
	Bio                  string     `json:"bio"`
	Gender               string     `json:"gender"`
	DateOfBirth          *time.Time `json:"date_of_birth,omitempty"`
	NewsletterSubscribed bool       `json:"newsletter_subscribed"`
}

type LoginRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	CaptchaToken string `json:"captcha_token"`
}

// TokenResponse encapsulates the security tokens issued upon successful authentication.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Session represents an active user session tracked in Redis.
// It contains metadata about the device, location, and activity to help users manage their account security.
type Session struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	DeviceType   string    `json:"device_type"` // Categorized as 'Mobile' or 'Desktop'
	Location     string    `json:"location"`    // Geographic location (e.g., 'Jakarta, Indonesia')
	LastActiveAt time.Time `json:"last_active_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	IsCurrent    bool      `json:"is_current"` // True if this session matches the one used for the current request
}
