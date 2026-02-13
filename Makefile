# GoTasker Makefile

# Database configuration
DB_URL ?= postgres://gotasker:password@localhost:5432/gotasker_db?sslmode=disable
MIGRATIONS_PATH = migrations

# Install migrate CLI tool
.PHONY: migrate-install
migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create a new migration
# Usage: make migrate-create name=create_users_table
.PHONY: migrate-create
migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

# Run all pending migrations
.PHONY: migrate-up
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

# Rollback the last migration
.PHONY: migrate-down
migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

# Rollback all migrations
.PHONY: migrate-down-all
migrate-down-all:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down -all

# Force set migration version (use with caution)
# Usage: make migrate-force version=1
.PHONY: migrate-force
migrate-force:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force $(version)

# Show current migration version
.PHONY: migrate-version
migrate-version:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version

# Run the application
.PHONY: run
run:
	go run src/cmd/main.go

# Build the application
.PHONY: build
build:
	go build -o bin/gotasker src/cmd/main.go

# Start docker-compose services
.PHONY: docker-up
docker-up:
	docker-compose up -d

# Stop docker-compose services
.PHONY: docker-down
docker-down:
	docker-compose down

# Show help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  migrate-install   - Install migrate CLI tool"
	@echo "  migrate-create    - Create new migration (name=migration_name)"
	@echo "  migrate-up        - Run all pending migrations"
	@echo "  migrate-down      - Rollback the last migration"
	@echo "  migrate-down-all  - Rollback all migrations"
	@echo "  migrate-force     - Force set migration version (version=N)"
	@echo "  migrate-version   - Show current migration version"
	@echo "  run               - Run the application"
	@echo "  build             - Build the application"
	@echo "  docker-up         - Start docker-compose services"
	@echo "  docker-down       - Stop docker-compose services"
