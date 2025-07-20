# Payroll Service

This is a payroll service written in Go.

## Prerequisites

- Go (version 1.24 or later)
- Docker and Docker Compose
- Make (optional, for using the Makefile)
- PostgreSQL client for database access
- [golangci-lint](https://golangci-lint.run/usage/install/) (for linting)

## Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/restuu/payroll.git
    cd payroll
    ```

2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Set up environment variables:**
    Copy the `.env.example` file to `.env` and update the values for your local environment.
    ```bash
    cp .env.example .env
    ```

## Running the Application

### Using Docker (Recommended)

This is the easiest way to get the application and the database running.

```bash
make docker-up
```

This will build the Docker image and start the application and PostgreSQL database containers. The application will be available at `http://localhost:8080`.

To stop the containers, run:
```bash
make docker-down
```

### Running Locally

If you prefer to run the application directly on your machine:

1.  **Start a PostgreSQL database.** You can use Docker or a local installation. Make sure the connection details in your `.env` file are correct.

2.  **Run database migrations:**
    ```bash
    make migrate-up
    ```

3.  **Run the application:**
    ```bash
    make run
    ```

    Alternatively, you can build the binary first and then run it:
    ```bash
    make build
    ./out/payroll-server
    ```

## Kafka Integration

This service uses Apache Kafka for asynchronous event handling. The Docker Compose setup includes the necessary services to run Kafka locally.

*   **Kafka:** The message broker itself.
*   **Zookeeper:** Required for Kafka coordination.
*   **Kafka UI (Redpanda Console):** A web-based user interface for managing and monitoring Kafka. You can access it at [http://localhost:8081](http://localhost:8081).

When the application starts, it subscribes to the following topics:
*   `payroll-generate`

Messages sent to these topics will be processed by the corresponding handlers in the `internal/presentation/messaging/handlers` package.

## Available `make` Commands

-   `build`: Builds the Go binary into the `out/` directory.
-   `run`: Runs the application in development mode.
-   `test`: Runs the unit tests.
-   `test-coverage`: Runs tests and generates an HTML coverage report.
-   `docker-up`: Starts the application and database using Docker Compose.
-   `docker-down`: Stops the Docker Compose containers.
-   `lint`: Lints the codebase using `golangci-lint`.
-   `wire`: Runs `wire` to generate dependency injection code.
-   `sqlc-generate`: Generates Go code from SQL queries using `sqlc`.
-   `migrate-create name=<migration_name>`: Creates a new SQL migration file.
-   `migrate-up`: Applies all pending database migrations.
-   `migrate-down`: Rolls back the last database migration.
-   `clean`: Cleans up build artifacts.
