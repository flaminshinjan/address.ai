package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/shinjan/address.ai/services/room/config"
)

type Room struct {
	ID            string    `json:"id"`
	RoomNumber    string    `json:"room_number"`
	RoomType      string    `json:"room_type"`
	Description   string    `json:"description"`
	PricePerNight float64   `json:"price_per_night"`
	Capacity      int       `json:"capacity"`
	Amenities     []string  `json:"amenities"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Booking struct {
	ID              string    `json:"id"`
	RoomID          string    `json:"room_id"`
	UserID          string    `json:"user_id"`
	CheckInDate     string    `json:"check_in_date"`
	CheckOutDate    string    `json:"check_out_date"`
	TotalPrice      float64   `json:"total_price"`
	Status          string    `json:"status"`
	SpecialRequests string    `json:"special_requests"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func GetRooms(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query(`
		SELECT id, room_number, room_type, description, price_per_night, 
		       capacity, amenities, status, created_at, updated_at 
		FROM rooms
	`)
	if err != nil {
		log.Printf("Error querying rooms: %v", err)
		http.Error(w, "Failed to get rooms", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var rooms []Room
	for rows.Next() {
		var room Room
		var amenitiesStr string
		err := rows.Scan(
			&room.ID, &room.RoomNumber, &room.RoomType, &room.Description,
			&room.PricePerNight, &room.Capacity, &amenitiesStr, &room.Status,
			&room.CreatedAt, &room.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning room row: %v", err)
			continue
		}

		// Parse amenities string array
		err = json.Unmarshal([]byte(amenitiesStr), &room.Amenities)
		if err != nil {
			log.Printf("Error parsing amenities: %v", err)
			room.Amenities = []string{}
		}

		rooms = append(rooms, room)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}

func CreateBooking(w http.ResponseWriter, r *http.Request) {
	var booking Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Start a transaction
	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Check if room is available
	var status string
	err = tx.QueryRow(`
		SELECT status FROM rooms WHERE id = $1
	`, booking.RoomID).Scan(&status)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}
	if status != "available" {
		http.Error(w, "Room is not available", http.StatusBadRequest)
		return
	}

	// Create booking
	err = tx.QueryRow(`
		INSERT INTO bookings (
			room_id, user_id, check_in_date, check_out_date,
			total_price, status, special_requests
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`,
		booking.RoomID, booking.UserID, booking.CheckInDate,
		booking.CheckOutDate, booking.TotalPrice, "confirmed",
		booking.SpecialRequests,
	).Scan(&booking.ID, &booking.CreatedAt, &booking.UpdatedAt)
	if err != nil {
		http.Error(w, "Failed to create booking", http.StatusInternalServerError)
		return
	}

	// Update room status
	_, err = tx.Exec(`
		UPDATE rooms SET status = 'booked' WHERE id = $1
	`, booking.RoomID)
	if err != nil {
		http.Error(w, "Failed to update room status", http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

func GetBookings(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	rows, err := config.DB.Query(`
		SELECT id, room_id, user_id, check_in_date, check_out_date,
		       total_price, status, special_requests, created_at, updated_at
		FROM bookings WHERE user_id = $1
	`, userID)
	if err != nil {
		http.Error(w, "Failed to get bookings", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var booking Booking
		err := rows.Scan(
			&booking.ID, &booking.RoomID, &booking.UserID,
			&booking.CheckInDate, &booking.CheckOutDate,
			&booking.TotalPrice, &booking.Status,
			&booking.SpecialRequests, &booking.CreatedAt,
			&booking.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning booking row: %v", err)
			continue
		}
		bookings = append(bookings, booking)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}
