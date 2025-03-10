package model

import (
	"time"
)

// Room represents a hotel room
type Room struct {
	ID          string    `json:"id"`
	Number      string    `json:"number"`
	Type        string    `json:"type"`
	Floor       int       `json:"floor"`
	Description string    `json:"description"`
	Capacity    int       `json:"capacity"`
	PricePerDay float64   `json:"price_per_day"`
	Status      string    `json:"status"` // available, occupied, maintenance
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RoomType represents a type of room
type RoomType struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BasePrice   float64   `json:"base_price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Booking represents a room booking
type Booking struct {
	ID         string    `json:"id"`
	RoomID     string    `json:"room_id"`
	UserID     string    `json:"user_id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"` // confirmed, cancelled, completed
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// BookingRequest represents a request to book a room
type BookingRequest struct {
	RoomID    string    `json:"room_id" validate:"required"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`
}

// RoomResponse represents the room data returned in responses
type RoomResponse struct {
	ID          string    `json:"id"`
	Number      string    `json:"number"`
	Type        string    `json:"type"`
	Floor       int       `json:"floor"`
	Description string    `json:"description"`
	Capacity    int       `json:"capacity"`
	PricePerDay float64   `json:"price_per_day"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// BookingResponse represents the booking data returned in responses
type BookingResponse struct {
	ID         string       `json:"id"`
	Room       RoomResponse `json:"room"`
	UserID     string       `json:"user_id"`
	StartDate  time.Time    `json:"start_date"`
	EndDate    time.Time    `json:"end_date"`
	TotalPrice float64      `json:"total_price"`
	Status     string       `json:"status"`
	CreatedAt  time.Time    `json:"created_at"`
}

// ToResponse converts a Room to a RoomResponse
func (r *Room) ToResponse() RoomResponse {
	return RoomResponse{
		ID:          r.ID,
		Number:      r.Number,
		Type:        r.Type,
		Floor:       r.Floor,
		Description: r.Description,
		Capacity:    r.Capacity,
		PricePerDay: r.PricePerDay,
		Status:      r.Status,
		CreatedAt:   r.CreatedAt,
	}
}
