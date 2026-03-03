package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found, using system environment variables")
	}

	// Debug: list files in migrations directory
	files, _ := filepath.Glob("migrations/*.sql")
	log.Printf("Found %d migration files in folder", len(files))
	for _, f := range files {
		log.Printf(" - %s", f)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		log.Fatal(err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatal(err)
	}
	log.Printf("Database version before: %d, Dirty: %v", version, dirty)

	cmd := "up"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("Migrations applied successfully!")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("Migrations rolled back successfully!")
	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Version number required for force command")
		}
		v, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("Invalid version number")
		}
		if err := m.Force(v); err != nil {
			log.Fatal(err)
		}
		log.Printf("Forced version to %d", v)
	default:
		log.Fatalf("Unknown command: %s", cmd)
	}

	newVersion, newDirty, _ := m.Version()
	log.Printf("Database version after: %d, Dirty: %v", newVersion, newDirty)
}
