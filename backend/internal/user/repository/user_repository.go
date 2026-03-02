package repository

import (
	"context"
	"errors"
	"swiftly/backend/internal/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(u *user.User) error {
	if r.pool == nil {
		return errors.New("database pool not initialized")
	}

	if u.ID == "" {
		u.ID = uuid.New().String()
	}

	query := `INSERT INTO users (id, email, full_name, created_at, updated_at) 
	          VALUES ($1, $2, $3, NOW(), NOW())`
	
	_, err := r.pool.Exec(context.Background(), query, u.ID, u.Email, u.FullName)
	return err
}

func (r *UserRepository) GetByID(id string) (*user.User, error) {
	if r.pool == nil {
		return nil, errors.New("database pool not initialized")
	}

	var u user.User
	query := `SELECT id, email, full_name, created_at, updated_at FROM users WHERE id = $1`
	
	err := r.pool.QueryRow(context.Background(), query, id).Scan(
		&u.ID, &u.Email, &u.FullName, &u.CreatedAt, &u.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) GetByEmail(email string) (*user.User, error) {
	if r.pool == nil {
		return nil, errors.New("database pool not initialized")
	}

	var u user.User
	query := `SELECT id, email, full_name, created_at, updated_at FROM users WHERE email = $1`
	
	err := r.pool.QueryRow(context.Background(), query, email).Scan(
		&u.ID, &u.Email, &u.FullName, &u.CreatedAt, &u.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &u, nil
}
