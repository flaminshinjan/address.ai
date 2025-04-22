# Build stage
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum from root
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY services/supply/ ./services/supply/
COPY pkg/ ./pkg/

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

# Copy migrations
COPY --from=builder /app/services/supply/migrations ./migrations

# Expose the port
EXPOSE 8084

# Set the entry point
CMD ["./supply-service"] 