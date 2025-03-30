# Build stage
FROM golang:1.24-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git ca-certificates tzdata && \
    update-ca-certificates

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/stablemcp ./cmd/server.go

# Create the directory structure for the final image and set permissions
RUN mkdir -p /app/configs && \
    chown -R 65534:65534 /app/bin/stablemcp /app/configs && \
    chmod -R 755 /app/bin/stablemcp /app/configs

# Final stage - using distroless
FROM gcr.io/distroless/static-debian12

# Set user to 65534 (nobody)
USER 65534:65534

# Set working directory
WORKDIR /app

# Copy binary from builder stage with correct ownership
COPY --from=builder /app/bin/stablemcp /app/stablemcp

# Copy configuration files with correct ownership
COPY --from=builder /app/configs /app/configs

# Copy CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Expose the application port
EXPOSE 8080

# Set environment variables
ENV GO_ENV=production

# Command to run the application
CMD ["/app/stablemcp", "--config", "/app/configs/default.yaml"]