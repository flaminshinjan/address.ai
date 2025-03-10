FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o room-service ./services/room/cmd

# Use a minimal alpine image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/room-service .

# Copy migrations
COPY --from=builder /app/services/room/migrations ./migrations

# Expose port
EXPOSE 8082

# Run the application
CMD ["./room-service"] 