# Build stage
FROM golang:1.23-alpine AS builder

# Install git for fetching dependencies
RUN apk add --no-cache git
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate Swagger docs
RUN swag init

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/hackers-auth

# Final stage
FROM alpine:3.19

# Add CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/hackers-auth .

# Copy generated docs
COPY --from=builder /app/docs ./docs

# Expose port
EXPOSE 8080

# Set environment variables
ENV GIN_MODE=release

# Run the application
CMD ["./hackers-auth"] 
