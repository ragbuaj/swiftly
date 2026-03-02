# Load environment variables from .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: help run build test migrate-up migrate-down migrate-status migrate-create seed

help:
	@echo "Swiftly Backend Management Commands:"
	@echo "  make run             - Run the API server"
	@echo "  make build           - Build the API binary"
	@echo "  make test            - Run all backend tests"
	@echo "  make migrate-up      - Run database migrations up"
	@echo "  make migrate-down    - Rollback the last migration"
	@echo "  make migrate-status  - Check migration status"
	@echo "  make migrate-create NAME=name_here - Create a new migration file"
	@echo "  make seed            - Seed the database with initial data"

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test ./... -v

migrate-up:
	go run cmd/migrate/main.go

migrate-down:
	@echo "Rollback is not natively supported by GORM AutoMigrate. Use SQL or migrations packages if needed."

migrate-status:
	@echo "GORM AutoMigrate is declarative. Status check is not applicable."

# Usage: make migrate-create NAME=create_users_table
migrate-create:
	goose -dir migrations create $(NAME) sql

seed:
	go run cmd/seed/main.go
