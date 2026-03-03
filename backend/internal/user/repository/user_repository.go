package repository

import (
	"context"
	"errors"
	"swiftly/backend/internal/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(u *user.User) error {
	if r.pool == nil {
		return errors.New("database pool not initialized")
	}

	if u.ID == "" {
		u.ID = uuid.New().String()
	}

	query := `INSERT INTO users (id, email, username, phone_number, full_name, password, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`
	
	_, err := r.pool.Exec(context.Background(), query, u.ID, u.Email, u.Username, u.PhoneNumber, u.FullName, u.Password)
	return err
}

func (r *Repository) GetByID(id string) (*user.User, error) {
	if r.pool == nil {
		return nil, errors.New("database pool not initialized")
	}

	var u user.User
	query := `SELECT id, email, username, phone_number, full_name, role, status, email_verified_at, deleted_at, date_of_birth, gender, newsletter_subscribed, avatar_url, bio, default_shipping_address_id, default_billing_address_id, created_at, updated_at FROM users WHERE id = $1`
	
	var defaultShippingAddressID, defaultBillingAddressID *string
	var username, phoneNumber, role, status, gender, avatarURL, bio *string

	err := r.pool.QueryRow(context.Background(), query, id).Scan(
		&u.ID, &u.Email, &username, &phoneNumber, &u.FullName, &role, &status, &u.EmailVerifiedAt, &u.DeletedAt, &u.DateOfBirth, &gender, &u.NewsletterSubscribed, &avatarURL, &bio, &defaultShippingAddressID, &defaultBillingAddressID, &u.CreatedAt, &u.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if username != nil { u.Username = *username }
	if phoneNumber != nil { u.PhoneNumber = *phoneNumber }
	if role != nil { u.Role = *role }
	if status != nil { u.Status = *status }
	if gender != nil { u.Gender = *gender }
	if avatarURL != nil { u.AvatarURL = *avatarURL }
	if bio != nil { u.Bio = *bio }
	if defaultShippingAddressID != nil { u.DefaultShippingAddressID = *defaultShippingAddressID }
	if defaultBillingAddressID != nil { u.DefaultBillingAddressID = *defaultBillingAddressID }

	return &u, nil
}

func (r *Repository) GetByEmail(email string) (*user.User, error) {
	if r.pool == nil {
		return nil, errors.New("database pool not initialized")
	}

	var u user.User
	query := `SELECT id, email, username, phone_number, full_name, password, role, status, newsletter_subscribed, created_at, updated_at FROM users WHERE email = $1`

	var username, phoneNumber, role, status *string
	var newsletterSubscribed *bool

	err := r.pool.QueryRow(context.Background(), query, email).Scan(
		&u.ID, &u.Email, &username, &phoneNumber, &u.FullName, &u.Password, &role, &status, &newsletterSubscribed, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if username != nil { u.Username = *username }
	if phoneNumber != nil { u.PhoneNumber = *phoneNumber }
	if role != nil { u.Role = *role }
	if status != nil { u.Status = *status }
	if newsletterSubscribed != nil { u.NewsletterSubscribed = *newsletterSubscribed }

	return &u, nil
}

func (r *Repository) GetByIdentifier(identifier string) (*user.User, error) {
	if r.pool == nil {
		return nil, errors.New("database pool not initialized")
	}

	var u user.User
	// Logic: Match Email OR (Match Phone AND Phone must be verified)
	// We removed Username from this lookup per your request.
	query := `SELECT id, email, username, phone_number, full_name, password, role, status, newsletter_subscribed, created_at, updated_at
	          FROM users
	          WHERE email = $1 OR (phone_number = $1 AND phone_verified_at IS NOT NULL)`

	var username, phoneNumber, role, status *string
	var newsletterSubscribed *bool

	err := r.pool.QueryRow(context.Background(), query, identifier).Scan(
		&u.ID, &u.Email, &username, &phoneNumber, &u.FullName, &u.Password, &role, &status, &newsletterSubscribed, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if username != nil { u.Username = *username }
	if phoneNumber != nil { u.PhoneNumber = *phoneNumber }
	if role != nil { u.Role = *role }
	if status != nil { u.Status = *status }
	if newsletterSubscribed != nil { u.NewsletterSubscribed = *newsletterSubscribed }

	return &u, nil
}
func (r *Repository) UpdatePassword(email, newHashedPassword string) error {
	if r.pool == nil {
		return errors.New("database pool not initialized")
	}

	query := `UPDATE users SET password = $1, updated_at = NOW() WHERE email = $2`
	_, err := r.pool.Exec(context.Background(), query, newHashedPassword, email)
	return err
}

func (r *Repository) VerifyPhone(email string) error {
	if r.pool == nil {
		return errors.New("database pool not initialized")
	}

	query := `UPDATE users SET phone_verified_at = NOW(), updated_at = NOW() WHERE email = $1`
	_, err := r.pool.Exec(context.Background(), query, email)
	return err
}

func (r *Repository) UpdateProfile(id string, req *user.UpdateProfileRequest) error {
	if r.pool == nil {
		return errors.New("database pool not initialized")
	}

	query := `UPDATE users 
			  SET full_name = $1, username = $2, phone_number = $3, bio = $4, gender = $5, date_of_birth = $6, newsletter_subscribed = $7, updated_at = NOW() 
			  WHERE id = $8`
	
	_, err := r.pool.Exec(context.Background(), query, 
		req.FullName, req.Username, req.PhoneNumber, req.Bio, req.Gender, req.DateOfBirth, req.NewsletterSubscribed, id)
	return err
}

func (r *Repository) UpdateAvatar(id string, avatarURL string) error {
	if r.pool == nil {
		return errors.New("database pool not initialized")
	}

	query := `UPDATE users SET avatar_url = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(context.Background(), query, avatarURL, id)
	return err
}

