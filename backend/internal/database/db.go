package database

import (
	"context"
	"log"
	"os"

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
