# Stage 1: Build the Go binary in a dedicated build environment
FROM golang:1.24-alpine AS builder

# Install tzdata for time zone support, ca-certificates for SSL, and git.
# These will be copied to the final image as needed.
RUN apk add --no-cache ca-certificates git tzdata

WORKDIR /app

# Copy go.mod and go.sum to leverage Docker layer caching for dependencies
COPY go.mod go.sum ./

# Use Docker's build cache for Go modules to speed up subsequent builds
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy the rest of the application source code
COPY . .

# Build the application binary.
# -ldflags="-w -s" strips debug information, reducing the binary size.
# CGO_ENABLED=0 creates a statically linked binary suitable for scratch images.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/payroll-server ./cmd/webserver

# Stage 2: Create the final, minimal production image
FROM scratch

# Copy timezone and SSL certificate data from the builder stage
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/payroll-server /payroll-server

# Set the ZONEINFO environment variable so Go can find timezone data
ENV ZONEINFO=/usr/share/zoneinfo

# Expose the port the server runs on (from development.config.yaml)
EXPOSE 8080

# Set the entrypoint for the container to run the server
ENTRYPOINT ["/payroll-server"]