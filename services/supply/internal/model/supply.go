package model

import (
	"time"
)

// Supplier represents a supplier
type Supplier struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// InventoryItem represents an inventory item
type InventoryItem struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Description string    `json:"description,omitempty"`
	Quantity    int       `json:"quantity"`
	Unit        string    `json:"unit"`
	MinQuantity int       `json:"min_quantity"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PurchaseOrder represents a purchase order
type PurchaseOrder struct {
	ID           string    `json:"id"`
	SupplierID   string    `json:"supplier_id"`
	Status       string    `json:"status"` // pending, approved, received, cancelled
	TotalPrice   float64   `json:"total_price"`
	Notes        string    `json:"notes,omitempty"`
	OrderDate    time.Time `json:"order_date"`
	DeliveryDate time.Time `json:"delivery_date,omitempty"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// OrderItem represents an item in a purchase order
type OrderItem struct {
	ID              string    `json:"id"`
	PurchaseOrderID string    `json:"purchase_order_id"`
	InventoryItemID string    `json:"inventory_item_id"`
	Quantity        int       `json:"quantity"`
	Price           float64   `json:"price"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// InventoryTransaction represents an inventory transaction
type InventoryTransaction struct {
	ID              string    `json:"id"`
	InventoryItemID string    `json:"inventory_item_id"`
	Quantity        int       `json:"quantity"`
	Type            string    `json:"type"`             // in, out
	Source          string    `json:"source,omitempty"` // purchase_order, consumption, adjustment
	SourceID        string    `json:"source_id,omitempty"`
	Notes           string    `json:"notes,omitempty"`
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
}

// SupplierResponse represents the supplier data returned in responses
type SupplierResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

// InventoryItemResponse represents the inventory item data returned in responses
type InventoryItemResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Description string    `json:"description,omitempty"`
	Quantity    int       `json:"quantity"`
	Unit        string    `json:"unit"`
	MinQuantity int       `json:"min_quantity"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

// OrderItemResponse represents the order item data returned in responses
type OrderItemResponse struct {
	ID            string                `json:"id"`
	InventoryItem InventoryItemResponse `json:"inventory_item"`
	Quantity      int                   `json:"quantity"`
	Price         float64               `json:"price"`
}

// PurchaseOrderResponse represents the purchase order data returned in responses
type PurchaseOrderResponse struct {
	ID           string              `json:"id"`
	Supplier     SupplierResponse    `json:"supplier"`
	Status       string              `json:"status"`
	TotalPrice   float64             `json:"total_price"`
	Notes        string              `json:"notes,omitempty"`
	OrderDate    time.Time           `json:"order_date"`
	DeliveryDate time.Time           `json:"delivery_date,omitempty"`
	CreatedBy    string              `json:"created_by"`
	Items        []OrderItemResponse `json:"items"`
	CreatedAt    time.Time           `json:"created_at"`
}

// CreatePurchaseOrderRequest represents a request to create a purchase order
type CreatePurchaseOrderRequest struct {
	SupplierID string             `json:"supplier_id" validate:"required"`
	Items      []OrderItemRequest `json:"items" validate:"required,min=1"`
	Notes      string             `json:"notes,omitempty"`
}

// OrderItemRequest represents an item in a create purchase order request
type OrderItemRequest struct {
	InventoryItemID string `json:"inventory_item_id" validate:"required"`
	Quantity        int    `json:"quantity" validate:"required,min=1"`
}

// ToResponse converts a Supplier to a SupplierResponse
func (s *Supplier) ToResponse() SupplierResponse {
	return SupplierResponse{
		ID:          s.ID,
		Name:        s.Name,
		Email:       s.Email,
		Phone:       s.Phone,
		Address:     s.Address,
		Description: s.Description,
		IsActive:    s.IsActive,
		CreatedAt:   s.CreatedAt,
	}
}

// ToResponse converts an InventoryItem to an InventoryItemResponse
func (i *InventoryItem) ToResponse() InventoryItemResponse {
	return InventoryItemResponse{
		ID:          i.ID,
		Name:        i.Name,
		Category:    i.Category,
		Description: i.Description,
		Quantity:    i.Quantity,
		Unit:        i.Unit,
		MinQuantity: i.MinQuantity,
		Price:       i.Price,
		CreatedAt:   i.CreatedAt,
	}
}
