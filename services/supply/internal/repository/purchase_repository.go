package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/hotel-management/services/supply/internal/model"
)

// PurchaseRepository handles database operations for purchase orders
type PurchaseRepository struct {
	db *sql.DB
}

// NewPurchaseRepository creates a new PurchaseRepository
func NewPurchaseRepository(db *sql.DB) *PurchaseRepository {
	return &PurchaseRepository{db: db}
}

// CreateOrder creates a new purchase order
func (r *PurchaseRepository) CreateOrder(order model.PurchaseOrder) (model.PurchaseOrder, error) {
	query := `
		INSERT INTO purchase_orders (id, supplier_id, status, total_price, notes, order_date, delivery_date, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, supplier_id, status, total_price, notes, order_date, delivery_date, created_by, created_at, updated_at
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

	// Set order date if not provided
	if order.OrderDate.IsZero() {
		order.OrderDate = now
	}

	err := r.db.QueryRow(
		query,
		order.ID,
		order.SupplierID,
		order.Status,
		order.TotalPrice,
		order.Notes,
		order.OrderDate,
		order.DeliveryDate,
		order.CreatedBy,
		order.CreatedAt,
		order.UpdatedAt,
	).Scan(
		&order.ID,
		&order.SupplierID,
		&order.Status,
		&order.TotalPrice,
		&order.Notes,
		&order.OrderDate,
		&order.DeliveryDate,
		&order.CreatedBy,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return model.PurchaseOrder{}, err
	}

	return order, nil
}

// CreateOrderItem creates a new order item
func (r *PurchaseRepository) CreateOrderItem(item model.OrderItem) (model.OrderItem, error) {
	query := `
		INSERT INTO order_items (id, purchase_order_id, inventory_item_id, quantity, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, purchase_order_id, inventory_item_id, quantity, price, created_at, updated_at
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
		item.PurchaseOrderID,
		item.InventoryItemID,
		item.Quantity,
		item.Price,
		item.CreatedAt,
		item.UpdatedAt,
	).Scan(
		&item.ID,
		&item.PurchaseOrderID,
		&item.InventoryItemID,
		&item.Quantity,
		&item.Price,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return model.OrderItem{}, err
	}

	return item, nil
}

// GetOrderByID gets a purchase order by ID
func (r *PurchaseRepository) GetOrderByID(id string) (model.PurchaseOrder, error) {
	query := `
		SELECT id, supplier_id, status, total_price, notes, order_date, delivery_date, created_by, created_at, updated_at
		FROM purchase_orders
		WHERE id = $1
	`

	var order model.PurchaseOrder
	err := r.db.QueryRow(query, id).Scan(
		&order.ID,
		&order.SupplierID,
		&order.Status,
		&order.TotalPrice,
		&order.Notes,
		&order.OrderDate,
		&order.DeliveryDate,
		&order.CreatedBy,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.PurchaseOrder{}, errors.New("purchase order not found")
		}
		return model.PurchaseOrder{}, err
	}

	return order, nil
}

// GetOrderItemsByOrderID gets order items by order ID
func (r *PurchaseRepository) GetOrderItemsByOrderID(orderID string) ([]model.OrderItem, error) {
	query := `
		SELECT id, purchase_order_id, inventory_item_id, quantity, price, created_at, updated_at
		FROM order_items
		WHERE purchase_order_id = $1
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
			&item.PurchaseOrderID,
			&item.InventoryItemID,
			&item.Quantity,
			&item.Price,
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

// UpdateOrderStatus updates a purchase order's status
func (r *PurchaseRepository) UpdateOrderStatus(id, status string) error {
	query := `
		UPDATE purchase_orders
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

// UpdateDeliveryDate updates a purchase order's delivery date
func (r *PurchaseRepository) UpdateDeliveryDate(id string, deliveryDate time.Time) error {
	query := `
		UPDATE purchase_orders
		SET delivery_date = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(query, deliveryDate, time.Now(), id)
	return err
}

// ListOrders lists all purchase orders
func (r *PurchaseRepository) ListOrders(limit, offset int) ([]model.PurchaseOrder, error) {
	query := `
		SELECT id, supplier_id, status, total_price, notes, order_date, delivery_date, created_by, created_at, updated_at
		FROM purchase_orders
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.PurchaseOrder
	for rows.Next() {
		var order model.PurchaseOrder
		err := rows.Scan(
			&order.ID,
			&order.SupplierID,
			&order.Status,
			&order.TotalPrice,
			&order.Notes,
			&order.OrderDate,
			&order.DeliveryDate,
			&order.CreatedBy,
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

// ListOrdersByStatus lists purchase orders by status
func (r *PurchaseRepository) ListOrdersByStatus(status string, limit, offset int) ([]model.PurchaseOrder, error) {
	query := `
		SELECT id, supplier_id, status, total_price, notes, order_date, delivery_date, created_by, created_at, updated_at
		FROM purchase_orders
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.PurchaseOrder
	for rows.Next() {
		var order model.PurchaseOrder
		err := rows.Scan(
			&order.ID,
			&order.SupplierID,
			&order.Status,
			&order.TotalPrice,
			&order.Notes,
			&order.OrderDate,
			&order.DeliveryDate,
			&order.CreatedBy,
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

// ListOrdersBySupplier lists purchase orders by supplier
func (r *PurchaseRepository) ListOrdersBySupplier(supplierID string, limit, offset int) ([]model.PurchaseOrder, error) {
	query := `
		SELECT id, supplier_id, status, total_price, notes, order_date, delivery_date, created_by, created_at, updated_at
		FROM purchase_orders
		WHERE supplier_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, supplierID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.PurchaseOrder
	for rows.Next() {
		var order model.PurchaseOrder
		err := rows.Scan(
			&order.ID,
			&order.SupplierID,
			&order.Status,
			&order.TotalPrice,
			&order.Notes,
			&order.OrderDate,
			&order.DeliveryDate,
			&order.CreatedBy,
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
