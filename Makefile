# Makefile for the payroll service

# Define the output binary name
BINARY_NAME=payroll-server
BINARY_PATH=out/$(BINARY_NAME)
FILENAME=

.PHONY: all test build service-up run clean wire

all: build

# test: Run all tests with race detector
test:
	@echo "Running tests..."
	go test -race -p 2 ./...

# build: Compile the application binary into the ./out directory
build:
	@echo "Building binary..."
	mkdir -p out
	go build -o $(BINARY_PATH) ./cmd/webserver

# service-up: Start the application and database using Docker Compose in detached mode
service-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up --build -d

# run: Run the compiled binary locally. Depends on 'build' to ensure the binary is up-to-date.
run: build
	@echo "Running binary from $(BINARY_PATH)..."
	./$(BINARY_PATH)

wire:
	go tool wire ./cmd/webserver/...

migration-create:
	./script/migration_create ${FILENAME}

migration-up:
	./script/migration_up

migration-down:
	./script/migration_down

sql-generate:
	go tool sqlc generate