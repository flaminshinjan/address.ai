# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o supply-service ./services/supply/cmd

# Final stage
FROM alpine:latest

# Set working directory
WORKDIR /app

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from the builder stage
COPY --from=builder /app/supply-service .

# Copy any additional files needed
COPY --from=builder /app/services/supply/.env.example ./.env

# Expose the port
EXPOSE 8083

# Set the entry point
CMD ["./supply-service"] 