package model

import (
	"time"
)

// MenuItem represents a food menu item
type MenuItem struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	IsAvailable bool      `json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MenuCategory represents a food menu category
type MenuCategory struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FoodOrder represents a food order
type FoodOrder struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	RoomID     string    `json:"room_id,omitempty"`
	Status     string    `json:"status"` // pending, preparing, delivered, cancelled
	TotalPrice float64   `json:"total_price"`
	Notes      string    `json:"notes,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// OrderItem represents an item in a food order
type OrderItem struct {
	ID         string    `json:"id"`
	OrderID    string    `json:"order_id"`
	MenuItemID string    `json:"menu_item_id"`
	Quantity   int       `json:"quantity"`
	Price      float64   `json:"price"`
	Notes      string    `json:"notes,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// MenuItemResponse represents the menu item data returned in responses
type MenuItemResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	IsAvailable bool      `json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
}

// OrderItemResponse represents the order item data returned in responses
type OrderItemResponse struct {
	ID       string           `json:"id"`
	MenuItem MenuItemResponse `json:"menu_item"`
	Quantity int              `json:"quantity"`
	Price    float64          `json:"price"`
	Notes    string           `json:"notes,omitempty"`
}

// FoodOrderResponse represents the food order data returned in responses
type FoodOrderResponse struct {
	ID         string              `json:"id"`
	UserID     string              `json:"user_id"`
	RoomID     string              `json:"room_id,omitempty"`
	Status     string              `json:"status"`
	TotalPrice float64             `json:"total_price"`
	Notes      string              `json:"notes,omitempty"`
	Items      []OrderItemResponse `json:"items"`
	CreatedAt  time.Time           `json:"created_at"`
}

// CreateOrderRequest represents a request to create a food order
type CreateOrderRequest struct {
	RoomID string             `json:"room_id,omitempty"`
	Items  []OrderItemRequest `json:"items" validate:"required,min=1"`
	Notes  string             `json:"notes,omitempty"`
}

// OrderItemRequest represents an item in a create order request
type OrderItemRequest struct {
	MenuItemID string `json:"menu_item_id" validate:"required"`
	Quantity   int    `json:"quantity" validate:"required,min=1"`
	Notes      string `json:"notes,omitempty"`
}

// ToResponse converts a MenuItem to a MenuItemResponse
func (m *MenuItem) ToResponse() MenuItemResponse {
	return MenuItemResponse{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Category:    m.Category,
		Price:       m.Price,
		IsAvailable: m.IsAvailable,
		CreatedAt:   m.CreatedAt,
	}
}
