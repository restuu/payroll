# Makefile for the payroll service

# Define the output binary name
BINARY_NAME=payroll-server
BINARY_PATH=out/$(BINARY_NAME)

# This line allows to load variables from .env file
-include .env

.PHONY: all build run test test-coverage lint wire sqlc-generate migrate-create migrate-up migrate-down docker-up docker-down clean

all: build

# build: Compile the application binary into the ./out directory
build:
	@echo "Building binary..."
	@mkdir -p out
	go build -o $(BINARY_PATH) ./cmd/webserver

# run: Run the application in development mode
run:
	@echo "Running application in development mode..."
	go run ./cmd/webserver/main.go

# test: Run all tests with race detector
test:
	@echo "Running tests..."
	go test -race -p 2 ./...

# test-coverage: Run tests and generate an HTML coverage report
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# lint: Lint the codebase using golangci-lint
# Note: Requires golangci-lint to be installed (https://golangci-lint.run/usage/install/)
lint:
	@echo "Linting codebase..."
	golangci-lint run

# wire: Generate dependency injection code using Wire
wire:
	@echo "Generating wire code..."
	go tool wire ./cmd/webserver/...

# sqlc-generate: Generate Go code from SQL queries using sqlc
sqlc-generate:
	@echo "Generating sqlc code..."
	go tool sqlc generate

# DB_CONN_STRING is required for migration commands. It is loaded from the .env file.
DB_CONN_STRING=user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} host=${POSTGRES_HOST} port=${POSTGRES_PORT} sslmode=disable

# migrate-create: Create a new SQL migration file. Usage: make migrate-create name=<migration_name>
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: migration name is not set. Use 'make migrate-create name=<migration_name>'"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)..."
	go tool goose -dir ./internal/infrastructure/database/postgres/migration postgres "$(DB_CONN_STRING)" create $(name) sql

# migrate-up: Apply all pending database migrations
migrate-up:
	@echo "Applying migrations..."
	go tool goose -dir ./internal/infrastructure/database/postgres/migration postgres "$(DB_CONN_STRING)" up

# migrate-down: Roll back the last database migration
migrate-down:
	@echo "Rolling back last migration..."
	go tool goose -dir ./internal/infrastructure/database/postgres/migration postgres "$(DB_CONN_STRING)" down

# docker-up: Start the application and database using Docker Compose
docker-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up --build -d

# docker-down: Stop the Docker Compose containers
docker-down:
	@echo "Stopping services..."
	docker-compose down

# clean: Clean up build artifacts and binaries
clean:
	@echo "Cleaning up..."
	rm -rf ./out
	rm -f coverage.out
