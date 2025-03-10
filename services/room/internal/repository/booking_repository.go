package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/flaminshinjan/address.ai/services/room/internal/model"
)

// BookingRepository handles database operations for bookings
type BookingRepository struct {
	db *sql.DB
}

// NewBookingRepository creates a new BookingRepository
func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

// Create creates a new booking
func (r *BookingRepository) Create(booking model.Booking) (model.Booking, error) {
	query := `
		INSERT INTO bookings (id, room_id, user_id, start_date, end_date, total_price, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, room_id, user_id, start_date, end_date, total_price, status, created_at, updated_at
	`

	// Generate UUID if not provided
	if booking.ID == "" {
		booking.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	booking.CreatedAt = now
	booking.UpdatedAt = now

	// Set default status if not provided
	if booking.Status == "" {
		booking.Status = "confirmed"
	}

	err := r.db.QueryRow(
		query,
		booking.ID,
		booking.RoomID,
		booking.UserID,
		booking.StartDate,
		booking.EndDate,
		booking.TotalPrice,
		booking.Status,
		booking.CreatedAt,
		booking.UpdatedAt,
	).Scan(
		&booking.ID,
		&booking.RoomID,
		&booking.UserID,
		&booking.StartDate,
		&booking.EndDate,
		&booking.TotalPrice,
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err != nil {
		return model.Booking{}, err
	}

	return booking, nil
}

// GetByID gets a booking by ID
func (r *BookingRepository) GetByID(id string) (model.Booking, error) {
	query := `
		SELECT id, room_id, user_id, start_date, end_date, total_price, status, created_at, updated_at
		FROM bookings
		WHERE id = $1
	`

	var booking model.Booking
	err := r.db.QueryRow(query, id).Scan(
		&booking.ID,
		&booking.RoomID,
		&booking.UserID,
		&booking.StartDate,
		&booking.EndDate,
		&booking.TotalPrice,
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Booking{}, errors.New("booking not found")
		}
		return model.Booking{}, err
	}

	return booking, nil
}

// GetByUserID gets bookings by user ID
func (r *BookingRepository) GetByUserID(userID string, limit, offset int) ([]model.Booking, error) {
	query := `
		SELECT id, room_id, user_id, start_date, end_date, total_price, status, created_at, updated_at
		FROM bookings
		WHERE user_id = $1
		ORDER BY start_date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []model.Booking
	for rows.Next() {
		var booking model.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.RoomID,
			&booking.UserID,
			&booking.StartDate,
			&booking.EndDate,
			&booking.TotalPrice,
			&booking.Status,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

// GetByRoomID gets bookings by room ID
func (r *BookingRepository) GetByRoomID(roomID string, limit, offset int) ([]model.Booking, error) {
	query := `
		SELECT id, room_id, user_id, start_date, end_date, total_price, status, created_at, updated_at
		FROM bookings
		WHERE room_id = $1
		ORDER BY start_date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, roomID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []model.Booking
	for rows.Next() {
		var booking model.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.RoomID,
			&booking.UserID,
			&booking.StartDate,
			&booking.EndDate,
			&booking.TotalPrice,
			&booking.Status,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

// Update updates a booking
func (r *BookingRepository) Update(booking model.Booking) (model.Booking, error) {
	query := `
		UPDATE bookings
		SET room_id = $1, user_id = $2, start_date = $3, end_date = $4, total_price = $5, status = $6, updated_at = $7
		WHERE id = $8
		RETURNING id, room_id, user_id, start_date, end_date, total_price, status, created_at, updated_at
	`

	booking.UpdatedAt = time.Now()

	err := r.db.QueryRow(
		query,
		booking.RoomID,
		booking.UserID,
		booking.StartDate,
		booking.EndDate,
		booking.TotalPrice,
		booking.Status,
		booking.UpdatedAt,
		booking.ID,
	).Scan(
		&booking.ID,
		&booking.RoomID,
		&booking.UserID,
		&booking.StartDate,
		&booking.EndDate,
		&booking.TotalPrice,
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Booking{}, errors.New("booking not found")
		}
		return model.Booking{}, err
	}

	return booking, nil
}

// UpdateStatus updates a booking's status
func (r *BookingRepository) UpdateStatus(id, status string) error {
	query := `
		UPDATE bookings
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

// Delete deletes a booking
func (r *BookingRepository) Delete(id string) error {
	query := `
		DELETE FROM bookings
		WHERE id = $1
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("booking not found")
	}

	return nil
}

// List lists all bookings
func (r *BookingRepository) List(limit, offset int) ([]model.Booking, error) {
	query := `
		SELECT id, room_id, user_id, start_date, end_date, total_price, status, created_at, updated_at
		FROM bookings
		ORDER BY start_date DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []model.Booking
	for rows.Next() {
		var booking model.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.RoomID,
			&booking.UserID,
			&booking.StartDate,
			&booking.EndDate,
			&booking.TotalPrice,
			&booking.Status,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

// CheckRoomAvailability checks if a room is available for the given dates
func (r *BookingRepository) CheckRoomAvailability(roomID string, startDate, endDate time.Time) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM bookings
		WHERE room_id = $1
		AND status = 'confirmed'
		AND (
			(start_date <= $2 AND end_date >= $2)
			OR (start_date <= $3 AND end_date >= $3)
			OR (start_date >= $2 AND end_date <= $3)
		)
	`

	var count int
	err := r.db.QueryRow(query, roomID, startDate, endDate).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}
