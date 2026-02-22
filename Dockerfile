# Build stage
FROM golang:1.22.6-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set destination for COPY
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main .

# Production stage
FROM scratch

# Import ca-certificates from builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import timezone data from builder stage
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Create app user
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy built application
COPY --from=builder /app/main /app/main

# Copy UI (templates + static)
COPY --from=builder /app/ui /app/ui

# Set working directory
WORKDIR /app

# Create data directory for CSV files
RUN mkdir -p /app/data

# Expose port
EXPOSE 8080

# Add health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/main", "-health-check"]

# Set environment variables for production
ENV ENVIRONMENT=production
ENV PORT=8080
ENV CSV_PATH=/app/data/products.csv
ENV GIN_MODE=release

# Run the application
CMD ["/app/main"]