package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/hotel-management/services/food/internal/model"
)

// OrderRepository handles database operations for food orders
type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository creates a new OrderRepository
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// CreateOrder creates a new food order
func (r *OrderRepository) CreateOrder(order model.FoodOrder) (model.FoodOrder, error) {
	query := `
		INSERT INTO food_orders (id, user_id, room_id, status, total_price, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, room_id, status, total_price, notes, created_at, updated_at
	`

	// Generate UUID if not provided
	if order.ID == "" {
		order.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now

	// Set default status if not provided
	if order.Status == "" {
		order.Status = "pending"
	}

	err := r.db.QueryRow(
		query,
		order.ID,
		order.UserID,
		order.RoomID,
		order.Status,
		order.TotalPrice,
		order.Notes,
		order.CreatedAt,
		order.UpdatedAt,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.RoomID,
		&order.Status,
		&order.TotalPrice,
		&order.Notes,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return model.FoodOrder{}, err
	}

	return order, nil
}

// CreateOrderItem creates a new order item
func (r *OrderRepository) CreateOrderItem(item model.OrderItem) (model.OrderItem, error) {
	query := `
		INSERT INTO order_items (id, order_id, menu_item_id, quantity, price, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, order_id, menu_item_id, quantity, price, notes, created_at, updated_at
	`

	// Generate UUID if not provided
	if item.ID == "" {
		item.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	err := r.db.QueryRow(
		query,
		item.ID,
		item.OrderID,
		item.MenuItemID,
		item.Quantity,
		item.Price,
		item.Notes,
		item.CreatedAt,
		item.UpdatedAt,
	).Scan(
		&item.ID,
		&item.OrderID,
		&item.MenuItemID,
		&item.Quantity,
		&item.Price,
		&item.Notes,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return model.OrderItem{}, err
	}

	return item, nil
}

// GetOrderByID gets a food order by ID
func (r *OrderRepository) GetOrderByID(id string) (model.FoodOrder, error) {
	query := `
		SELECT id, user_id, room_id, status, total_price, notes, created_at, updated_at
		FROM food_orders
		WHERE id = $1
	`

	var order model.FoodOrder
	err := r.db.QueryRow(query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.RoomID,
		&order.Status,
		&order.TotalPrice,
		&order.Notes,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.FoodOrder{}, errors.New("order not found")
		}
		return model.FoodOrder{}, err
	}

	return order, nil
}

// GetOrderItemsByOrderID gets order items by order ID
func (r *OrderRepository) GetOrderItemsByOrderID(orderID string) ([]model.OrderItem, error) {
	query := `
		SELECT id, order_id, menu_item_id, quantity, price, notes, created_at, updated_at
		FROM order_items
		WHERE order_id = $1
		ORDER BY created_at
	`

	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.OrderItem
	for rows.Next() {
		var item model.OrderItem
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.MenuItemID,
			&item.Quantity,
			&item.Price,
			&item.Notes,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// UpdateOrderStatus updates a food order's status
func (r *OrderRepository) UpdateOrderStatus(id, status string) error {
	query := `
		UPDATE food_orders
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

// ListOrdersByUserID lists food orders by user ID
func (r *OrderRepository) ListOrdersByUserID(userID string, limit, offset int) ([]model.FoodOrder, error) {
	query := `
		SELECT id, user_id, room_id, status, total_price, notes, created_at, updated_at
		FROM food_orders
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.FoodOrder
	for rows.Next() {
		var order model.FoodOrder
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.RoomID,
			&order.Status,
			&order.TotalPrice,
			&order.Notes,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

// ListOrders lists all food orders
func (r *OrderRepository) ListOrders(limit, offset int) ([]model.FoodOrder, error) {
	query := `
		SELECT id, user_id, room_id, status, total_price, notes, created_at, updated_at
		FROM food_orders
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.FoodOrder
	for rows.Next() {
		var order model.FoodOrder
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.RoomID,
			&order.Status,
			&order.TotalPrice,
			&order.Notes,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

// ListOrdersByStatus lists food orders by status
func (r *OrderRepository) ListOrdersByStatus(status string, limit, offset int) ([]model.FoodOrder, error) {
	query := `
		SELECT id, user_id, room_id, status, total_price, notes, created_at, updated_at
		FROM food_orders
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.FoodOrder
	for rows.Next() {
		var order model.FoodOrder
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.RoomID,
			&order.Status,
			&order.TotalPrice,
			&order.Notes,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
