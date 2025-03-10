package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/flaminshinjan/address.ai/services/user/internal/model"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user model.User) (model.User, error) {
	query := `
		INSERT INTO users (id, username, email, password, first_name, last_name, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, username, email, password, first_name, last_name, role, created_at, updated_at
	`

	// Generate UUID if not provided
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Set default role if not provided
	if user.Role == "" {
		user.Role = "guest"
	}

	err := r.db.QueryRow(
		query,
		user.ID,
		user.Username,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

// GetByID gets a user by ID
func (r *UserRepository) GetByID(id string) (model.User, error) {
	query := `
		SELECT id, username, email, password, first_name, last_name, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user model.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, errors.New("user not found")
		}
		return model.User{}, err
	}

	return user, nil
}

// GetByUsername gets a user by username
func (r *UserRepository) GetByUsername(username string) (model.User, error) {
	query := `
		SELECT id, username, email, password, first_name, last_name, role, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	var user model.User
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, errors.New("user not found")
		}
		return model.User{}, err
	}

	return user, nil
}

// GetByEmail gets a user by email
func (r *UserRepository) GetByEmail(email string) (model.User, error) {
	query := `
		SELECT id, username, email, password, first_name, last_name, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user model.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, errors.New("user not found")
		}
		return model.User{}, err
	}

	return user, nil
}

// Update updates a user
func (r *UserRepository) Update(user model.User) (model.User, error) {
	query := `
		UPDATE users
		SET username = $1, email = $2, first_name = $3, last_name = $4, role = $5, updated_at = $6
		WHERE id = $7
		RETURNING id, username, email, password, first_name, last_name, role, created_at, updated_at
	`

	user.UpdatedAt = time.Now()

	err := r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Role,
		user.UpdatedAt,
		user.ID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, errors.New("user not found")
		}
		return model.User{}, err
	}

	return user, nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(id, password string) error {
	query := `
		UPDATE users
		SET password = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(query, password, time.Now(), id)
	return err
}

// Delete deletes a user
func (r *UserRepository) Delete(id string) error {
	query := `
		DELETE FROM users
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
		return errors.New("user not found")
	}

	return nil
}

// List lists all users
func (r *UserRepository) List(limit, offset int) ([]model.User, error) {
	query := `
		SELECT id, username, email, password, first_name, last_name, role, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.FirstName,
			&user.LastName,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
