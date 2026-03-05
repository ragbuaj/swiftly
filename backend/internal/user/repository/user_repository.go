package repository

import (
	"context"
	"errors"
	"swiftly/backend/internal/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DBTX is an interface that matches both *pgxpool.Pool and pgx.Tx
type DBTX interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

type Repository struct {
	db DBTX
}

func NewUserRepository(db DBTX) *Repository {
	return &Repository{db: db}
}

// WithTx returns a new repository instance that uses the provided transaction
func (r *Repository) WithTx(tx pgx.Tx) *Repository {
	return &Repository{db: tx}
}

func (r *Repository) Create(ctx context.Context, u *user.User) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}

	query := `INSERT INTO users (id, email, username, phone_number, full_name, password, role, status, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
			  RETURNING id`
	
	err := r.db.QueryRow(ctx, query, u.ID, u.Email, u.Username, u.PhoneNumber, u.FullName, u.Password, u.Role, u.Status).Scan(&u.ID)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id string) (*user.User, error) {
	var u user.User
	query := `SELECT id, email, username, phone_number, full_name, role, status, email_verified_at, deleted_at, date_of_birth, gender, newsletter_subscribed, avatar_url, bio, default_shipping_address_id, default_billing_address_id, created_at, updated_at FROM users WHERE id = $1`
	
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.Username, &u.PhoneNumber, &u.FullName, &u.Role, &u.Status, &u.EmailVerifiedAt, &u.DeletedAt, &u.DateOfBirth, &u.Gender, &u.NewsletterSubscribed, &u.AvatarURL, &u.Bio, &u.DefaultShippingAddressID, &u.DefaultBillingAddressID, &u.CreatedAt, &u.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &u, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	query := `SELECT id, email, username, phone_number, full_name, password, role, status, email_verified_at, deleted_at, date_of_birth, gender, newsletter_subscribed, avatar_url, bio, default_shipping_address_id, default_billing_address_id, created_at, updated_at FROM users WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.Username, &u.PhoneNumber, &u.FullName, &u.Password, &u.Role, &u.Status, &u.EmailVerifiedAt, &u.DeletedAt, &u.DateOfBirth, &u.Gender, &u.NewsletterSubscribed, &u.AvatarURL, &u.Bio, &u.DefaultShippingAddressID, &u.DefaultBillingAddressID, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &u, nil
}

func (r *Repository) GetByIdentifier(ctx context.Context, identifier string) (*user.User, error) {
	var u user.User
	query := `SELECT id, email, username, phone_number, full_name, password, role, status, email_verified_at, deleted_at, date_of_birth, gender, newsletter_subscribed, avatar_url, bio, default_shipping_address_id, default_billing_address_id, created_at, updated_at
	          FROM users
	          WHERE email = $1 OR (phone_number = $1 AND phone_verified_at IS NOT NULL)`

	err := r.db.QueryRow(ctx, query, identifier).Scan(
		&u.ID, &u.Email, &u.Username, &u.PhoneNumber, &u.FullName, &u.Password, &u.Role, &u.Status, &u.EmailVerifiedAt, &u.DeletedAt, &u.DateOfBirth, &u.Gender, &u.NewsletterSubscribed, &u.AvatarURL, &u.Bio, &u.DefaultShippingAddressID, &u.DefaultBillingAddressID, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &u, nil
}

func (r *Repository) UpdatePassword(ctx context.Context, email, newHashedPassword string) error {
	query := `UPDATE users SET password = $1, updated_at = NOW() WHERE email = $2`
	_, err := r.db.Exec(ctx, query, newHashedPassword, email)
	return err
}

func (r *Repository) VerifyPhone(ctx context.Context, email string) error {
	query := `UPDATE users SET phone_verified_at = NOW(), updated_at = NOW() WHERE email = $1`
	_, err := r.db.Exec(ctx, query, email)
	return err
}

func (r *Repository) UpdateProfile(ctx context.Context, id string, req *user.UpdateProfileRequest) error {
	query := `UPDATE users 
			  SET full_name = $1, username = $2, phone_number = $3, bio = $4, gender = $5, date_of_birth = $6, newsletter_subscribed = $7, updated_at = NOW() 
			  WHERE id = $8`
	
	_, err := r.db.Exec(ctx, query, 
		req.FullName, req.Username, req.PhoneNumber, req.Bio, req.Gender, req.DateOfBirth, req.NewsletterSubscribed, id)
	return err
}

func (r *Repository) UpdateAvatar(ctx context.Context, id string, avatarURL string) error {
	query := `UPDATE users SET avatar_url = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, avatarURL, id)
	return err
}
