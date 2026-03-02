# Build stage
FROM golang:1.25.6-alpine AS builder

WORKDIR /app

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application (for production stage)
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# Run stage
FROM alpine:latest

# Install ca-certificates to allow HTTPS requests (needed for Supabase)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the API port
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
