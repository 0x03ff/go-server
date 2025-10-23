# Use the official Golang image for ARM
FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the application to bin/main (matching your project structure)
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/main ./cmd/api

# Final stage (smaller image)
FROM alpine:latest

# Create app directory structure
RUN mkdir -p /app/bin
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/bin/main /app/bin/main

# Make sure it's executable
RUN chmod +x /app/bin/main

# Create a non-root user (security best practice)
RUN adduser -D -u 1000 appuser
USER appuser

# Expose port
EXPOSE 8080

# Run the application using ABSOLUTE PATH (critical fix)
CMD ["/app/bin/main"]

