# Supply Chain Microservice

This microservice is responsible for managing the hotel's supply chain, including suppliers, inventory items, and purchase orders.

## Features

- Supplier management (CRUD operations)
- Inventory management (CRUD operations)
- Purchase order management (CRUD operations)
- Low stock alerts
- Inventory transaction tracking

## API Endpoints

### Suppliers

- `POST /api/v1/suppliers` - Create a new supplier (Admin only)
- `GET /api/v1/suppliers/{id}` - Get a supplier by ID
- `PUT /api/v1/suppliers/{id}` - Update a supplier (Admin only)
- `DELETE /api/v1/suppliers/{id}` - Delete a supplier (Admin only)
- `GET /api/v1/suppliers` - List all suppliers with pagination

### Inventory

- `POST /api/v1/inventory` - Create a new inventory item (Admin only)
- `GET /api/v1/inventory/{id}` - Get an inventory item by ID
- `PUT /api/v1/inventory/{id}` - Update an inventory item (Admin only)
- `PUT /api/v1/inventory/{id}/quantity` - Update an inventory item's quantity (Admin only)
- `DELETE /api/v1/inventory/{id}` - Delete an inventory item (Admin only)
- `GET /api/v1/inventory` - List all inventory items with pagination
- `GET /api/v1/inventory/category/{category}` - List inventory items by category
- `GET /api/v1/inventory/low-stock` - List items with quantities below the minimum threshold
- `GET /api/v1/inventory/categories` - List all inventory categories

### Purchase Orders

- `POST /api/v1/purchase-orders` - Create a new purchase order (Admin only)
- `GET /api/v1/purchase-orders/{id}` - Get a purchase order by ID (Admin only)
- `PUT /api/v1/purchase-orders/{id}/status` - Update a purchase order's status (Admin only)
- `GET /api/v1/purchase-orders` - List all purchase orders with pagination (Admin only)
- `GET /api/v1/purchase-orders/status/{status}` - List purchase orders by status (Admin only)

## Environment Variables

- `PORT` - The port the service will run on (default: 8083)
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - Secret key for JWT authentication

## Running the Service

### Using Docker

```bash
docker build -t hotel-supply-service .
docker run -p 8083:8083 --env-file .env hotel-supply-service
```

### Using Docker Compose

```bash
docker-compose up -d supply-service
```

### Locally

```bash
go run cmd/main.go
```

## Database Migrations

The service uses SQL migrations to set up the database schema. Migrations are located in the `migrations` directory.

To run migrations:

```bash
migrate -path ./migrations -database "${DATABASE_URL}" up
```

To rollback migrations:

```bash
migrate -path ./migrations -database "${DATABASE_URL}" down
``` 