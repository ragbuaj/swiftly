package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PasswordResetRepository struct {
	pool *pgxpool.Pool
}

func NewPasswordResetRepository(pool *pgxpool.Pool) *PasswordResetRepository {
	return &PasswordResetRepository{pool: pool}
}

func (r *PasswordResetRepository) StoreToken(email, token string, expiresAt time.Time) error {
	query := `INSERT INTO password_resets (email, token, expires_at, created_at) 
	          VALUES ($1, $2, $3, NOW())
	          ON CONFLICT (email) DO UPDATE 
	          SET token = EXCLUDED.token, expires_at = EXCLUDED.expires_at, created_at = NOW()`
	
	_, err := r.pool.Exec(context.Background(), query, email, token, expiresAt)
	return err
}

func (r *PasswordResetRepository) GetToken(token string) (string, time.Time, error) {
	var email string
	var expiresAt time.Time
	query := `SELECT email, expires_at FROM password_resets WHERE token = $1`
	
	err := r.pool.QueryRow(context.Background(), query, token).Scan(&email, &expiresAt)
	return email, expiresAt, err
}

func (r *PasswordResetRepository) DeleteToken(email string) error {
	query := `DELETE FROM password_resets WHERE email = $1`
	_, err := r.pool.Exec(context.Background(), query, email)
	return err
}
