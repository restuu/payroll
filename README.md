# Payroll Service

This is a payroll service written in Go.

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
  - [Running the application](#running-the-application)
  - [API Endpoints](#api-endpoints)
- [Database Migrations](#database-migrations)
- [Code Generation](#code-generation)
- [Running Tests](#running-tests)
- [Built With](#built-with)

## Features

- User authentication
- Attendance tracking

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- [Go](https://golang.org/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

### Installation

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/restuu/payroll.git
    cd payroll
    ```

2.  **Set up environment variables:**

    Copy the `.env.example` file to `.env` and update the variables for your local environment.

    ```bash
    cp .env.example .env
    ```

3.  **Install dependencies:**

    ```bash
    go mod tidy
    ```

## Usage

### Running the application

#### Using Docker (recommended)

This is the easiest way to get the application and the database running.

```bash
docker-compose up --build
```

The application will be available at `http://localhost:8080`.

#### Running locally

1.  **Build the application:**

    ```bash
    go build -o out/payroll-server ./cmd/webserver
    ```

2.  **Run the application:**

    ```bash
    ./out/payroll-server
    ```

    Or, for development:

    ```bash
    go run ./cmd/webserver/main.go
    ```

### API Endpoints

The available API endpoints are:

-   `GET /ping`: Health check endpoint.
-   Auth endpoints (see `internal/app/auth/delivery/http/auth_http_handler.go`)
-   Attendance endpoints (see `internal/app/attendance/delivery/http/attendance_http_handler.go`)

## Database Migrations

The project uses `goose` for database migrations.

-   **Create a new migration:**

    ```bash
    go run github.com/pressly/goose/v3/cmd/goose -dir ./internal/infrastructure/database/postgres/migration postgres "user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} host=${POSTGRES_HOST} port=${POSTGRES_PORT} sslmode=disable" create <migration_name> sql
    ```

-   **Apply all migrations:**

    ```bash
    go run github.com/pressly/goose/v3/cmd/goose -dir ./internal/infrastructure/database/postgres/migration postgres "user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} host=${POSTGRES_HOST} port=${POSTGRES_PORT} sslmode=disable" up
    ```

    *(Note: You need to have the environment variables from your `.env` file loaded in your shell for the above commands to work)*

## Code Generation

The project uses `sqlc` to generate Go code from SQL queries.

-   **Generate code:**

    ```bash
    go run github.com/sqlc-dev/sqlc/cmd/sqlc generate
    ```

## Running Tests

```bash
go test -race -p 2 ./...
```

## Built With

-   [Go](https://golang.org/)
-   [PostgreSQL](https://www.postgresql.org/)
-   [Docker](https://www.docker.com/)
-   [chi](https://github.com/go-chi/chi)
-   [sqlc](https://github.com/sqlc-dev/sqlc)
-   [goose](https://github.com/pressly/goose)
-   [wire](https://github.com/google/wire)
