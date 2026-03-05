package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

// Init initializes the pgx connection pool
func Init() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Println("Warning: DATABASE_URL not set")
		return
	}

	var err error
	Pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	log.Println("Successfully connected to database with pgx")
}

// GetPool returns the database connection pool
func GetPool() *pgxpool.Pool {
	return Pool
}

// Close closes the database connection pool
func Close() {
	if Pool != nil {
		Pool.Close()
	}
}

// WithTransaction provides a helper for executing multiple database operations within a transaction.
// If the callback returns an error, the transaction is automatically rolled back.
// If the callback returns nil, the transaction is automatically committed.
func WithTransaction(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p) // re-throw panic after rollback
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
