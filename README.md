# AI Hotel Management System

This is a microservice-based hotel management system built with Go, Docker, and PostgreSQL.

## Architecture

The system consists of four microservices:

1. **User Service**: Handles user authentication, registration, and profile management
2. **Room Management Service**: Manages room inventory, bookings, and availability
3. **Food Management Service**: Handles food orders, menu management, and kitchen operations
4. **Supply Chain Service**: Manages inventory, suppliers, and procurement

## Prerequisites

- Go 1.20+
- Docker and Docker Compose
- PostgreSQL (handled via Docker)

## Setup and Running

1. Clone the repository:
```bash
git clone <repository-url>
cd ai-hotel-management
```

2. Start the services using Docker Compose:
```bash
docker-compose up -d
```

3. The services will be available at:
   - User Service: http://localhost:8081
   - Room Management Service: http://localhost:8082
   - Food Management Service: http://localhost:8083
   - Supply Chain Service: http://localhost:8084

## API Documentation

Each service has its own API documentation available at `/swagger/index.html` when the service is running.

## Development

To run individual services during development:

```bash
# User Service
cd services/user
go run cmd/main.go

# Room Management Service
cd services/room
go run cmd/main.go

# Food Management Service
cd services/food
go run cmd/main.go

# Supply Chain Service
cd services/supply
go run cmd/main.go
```

## Testing

Run tests for all services:

```bash
go test ./...
```

Or for a specific service:

```bash
cd services/user
go test ./...
``` 