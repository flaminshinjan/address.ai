FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum from root
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY services/user/ ./services/user/
COPY pkg/ ./pkg/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-service ./services/user/cmd

# Use a minimal alpine image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/user-service .

# Copy migrations
COPY --from=builder /app/services/user/migrations ./migrations

# Expose port
EXPOSE 8081

# Run the application
CMD ["./user-service"] 