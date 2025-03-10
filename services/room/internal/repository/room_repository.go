package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/flaminshinjan/address.ai/services/room/internal/model"
)

// RoomRepository handles database operations for rooms
type RoomRepository struct {
	db *sql.DB
}

// NewRoomRepository creates a new RoomRepository
func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

// Create creates a new room
func (r *RoomRepository) Create(room model.Room) (model.Room, error) {
	query := `
		INSERT INTO rooms (id, number, type, floor, description, capacity, price_per_day, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, number, type, floor, description, capacity, price_per_day, status, created_at, updated_at
	`

	// Generate UUID if not provided
	if room.ID == "" {
		room.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	room.CreatedAt = now
	room.UpdatedAt = now

	// Set default status if not provided
	if room.Status == "" {
		room.Status = "available"
	}

	err := r.db.QueryRow(
		query,
		room.ID,
		room.Number,
		room.Type,
		room.Floor,
		room.Description,
		room.Capacity,
		room.PricePerDay,
		room.Status,
		room.CreatedAt,
		room.UpdatedAt,
	).Scan(
		&room.ID,
		&room.Number,
		&room.Type,
		&room.Floor,
		&room.Description,
		&room.Capacity,
		&room.PricePerDay,
		&room.Status,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return model.Room{}, err
	}

	return room, nil
}

// GetByID gets a room by ID
func (r *RoomRepository) GetByID(id string) (model.Room, error) {
	query := `
		SELECT id, number, type, floor, description, capacity, price_per_day, status, created_at, updated_at
		FROM rooms
		WHERE id = $1
	`

	var room model.Room
	err := r.db.QueryRow(query, id).Scan(
		&room.ID,
		&room.Number,
		&room.Type,
		&room.Floor,
		&room.Description,
		&room.Capacity,
		&room.PricePerDay,
		&room.Status,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Room{}, errors.New("room not found")
		}
		return model.Room{}, err
	}

	return room, nil
}

// GetByNumber gets a room by number
func (r *RoomRepository) GetByNumber(number string) (model.Room, error) {
	query := `
		SELECT id, number, type, floor, description, capacity, price_per_day, status, created_at, updated_at
		FROM rooms
		WHERE number = $1
	`

	var room model.Room
	err := r.db.QueryRow(query, number).Scan(
		&room.ID,
		&room.Number,
		&room.Type,
		&room.Floor,
		&room.Description,
		&room.Capacity,
		&room.PricePerDay,
		&room.Status,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Room{}, errors.New("room not found")
		}
		return model.Room{}, err
	}

	return room, nil
}

// Update updates a room
func (r *RoomRepository) Update(room model.Room) (model.Room, error) {
	query := `
		UPDATE rooms
		SET number = $1, type = $2, floor = $3, description = $4, capacity = $5, price_per_day = $6, status = $7, updated_at = $8
		WHERE id = $9
		RETURNING id, number, type, floor, description, capacity, price_per_day, status, created_at, updated_at
	`

	room.UpdatedAt = time.Now()

	err := r.db.QueryRow(
		query,
		room.Number,
		room.Type,
		room.Floor,
		room.Description,
		room.Capacity,
		room.PricePerDay,
		room.Status,
		room.UpdatedAt,
		room.ID,
	).Scan(
		&room.ID,
		&room.Number,
		&room.Type,
		&room.Floor,
		&room.Description,
		&room.Capacity,
		&room.PricePerDay,
		&room.Status,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Room{}, errors.New("room not found")
		}
		return model.Room{}, err
	}

	return room, nil
}

// Delete deletes a room
func (r *RoomRepository) Delete(id string) error {
	query := `
		DELETE FROM rooms
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
		return errors.New("room not found")
	}

	return nil
}

// List lists all rooms
func (r *RoomRepository) List(limit, offset int) ([]model.Room, error) {
	query := `
		SELECT id, number, type, floor, description, capacity, price_per_day, status, created_at, updated_at
		FROM rooms
		ORDER BY number ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []model.Room
	for rows.Next() {
		var room model.Room
		err := rows.Scan(
			&room.ID,
			&room.Number,
			&room.Type,
			&room.Floor,
			&room.Description,
			&room.Capacity,
			&room.PricePerDay,
			&room.Status,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

// ListAvailable lists all available rooms
func (r *RoomRepository) ListAvailable(startDate, endDate time.Time, limit, offset int) ([]model.Room, error) {
	query := `
		SELECT r.id, r.number, r.type, r.floor, r.description, r.capacity, r.price_per_day, r.status, r.created_at, r.updated_at
		FROM rooms r
		WHERE r.status = 'available'
		AND r.id NOT IN (
			SELECT b.room_id
			FROM bookings b
			WHERE b.status = 'confirmed'
			AND (
				(b.start_date <= $1 AND b.end_date >= $1)
				OR (b.start_date <= $2 AND b.end_date >= $2)
				OR (b.start_date >= $1 AND b.end_date <= $2)
			)
		)
		ORDER BY r.number ASC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(query, startDate, endDate, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []model.Room
	for rows.Next() {
		var room model.Room
		err := rows.Scan(
			&room.ID,
			&room.Number,
			&room.Type,
			&room.Floor,
			&room.Description,
			&room.Capacity,
			&room.PricePerDay,
			&room.Status,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

// UpdateStatus updates a room's status
func (r *RoomRepository) UpdateStatus(id, status string) error {
	query := `
		UPDATE rooms
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}
