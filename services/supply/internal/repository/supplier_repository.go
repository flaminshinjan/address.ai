package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/hotel-management/services/supply/internal/model"
)

// SupplierRepository handles database operations for suppliers
type SupplierRepository struct {
	db *sql.DB
}

// NewSupplierRepository creates a new SupplierRepository
func NewSupplierRepository(db *sql.DB) *SupplierRepository {
	return &SupplierRepository{db: db}
}

// Create creates a new supplier
func (r *SupplierRepository) Create(supplier model.Supplier) (model.Supplier, error) {
	query := `
		INSERT INTO suppliers (id, name, email, phone, address, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, name, email, phone, address, description, is_active, created_at, updated_at
	`

	// Generate UUID if not provided
	if supplier.ID == "" {
		supplier.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	supplier.CreatedAt = now
	supplier.UpdatedAt = now

	// Set default active status if not provided
	if !supplier.IsActive {
		supplier.IsActive = true
	}

	err := r.db.QueryRow(
		query,
		supplier.ID,
		supplier.Name,
		supplier.Email,
		supplier.Phone,
		supplier.Address,
		supplier.Description,
		supplier.IsActive,
		supplier.CreatedAt,
		supplier.UpdatedAt,
	).Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.Email,
		&supplier.Phone,
		&supplier.Address,
		&supplier.Description,
		&supplier.IsActive,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
	)

	if err != nil {
		return model.Supplier{}, err
	}

	return supplier, nil
}

// GetByID gets a supplier by ID
func (r *SupplierRepository) GetByID(id string) (model.Supplier, error) {
	query := `
		SELECT id, name, email, phone, address, description, is_active, created_at, updated_at
		FROM suppliers
		WHERE id = $1
	`

	var supplier model.Supplier
	err := r.db.QueryRow(query, id).Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.Email,
		&supplier.Phone,
		&supplier.Address,
		&supplier.Description,
		&supplier.IsActive,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Supplier{}, errors.New("supplier not found")
		}
		return model.Supplier{}, err
	}

	return supplier, nil
}

// Update updates a supplier
func (r *SupplierRepository) Update(supplier model.Supplier) (model.Supplier, error) {
	query := `
		UPDATE suppliers
		SET name = $1, email = $2, phone = $3, address = $4, description = $5, is_active = $6, updated_at = $7
		WHERE id = $8
		RETURNING id, name, email, phone, address, description, is_active, created_at, updated_at
	`

	supplier.UpdatedAt = time.Now()

	err := r.db.QueryRow(
		query,
		supplier.Name,
		supplier.Email,
		supplier.Phone,
		supplier.Address,
		supplier.Description,
		supplier.IsActive,
		supplier.UpdatedAt,
		supplier.ID,
	).Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.Email,
		&supplier.Phone,
		&supplier.Address,
		&supplier.Description,
		&supplier.IsActive,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Supplier{}, errors.New("supplier not found")
		}
		return model.Supplier{}, err
	}

	return supplier, nil
}

// Delete deletes a supplier
func (r *SupplierRepository) Delete(id string) error {
	query := `
		DELETE FROM suppliers
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
		return errors.New("supplier not found")
	}

	return nil
}

// List lists all suppliers
func (r *SupplierRepository) List(limit, offset int) ([]model.Supplier, error) {
	query := `
		SELECT id, name, email, phone, address, description, is_active, created_at, updated_at
		FROM suppliers
		ORDER BY name
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []model.Supplier
	for rows.Next() {
		var supplier model.Supplier
		err := rows.Scan(
			&supplier.ID,
			&supplier.Name,
			&supplier.Email,
			&supplier.Phone,
			&supplier.Address,
			&supplier.Description,
			&supplier.IsActive,
			&supplier.CreatedAt,
			&supplier.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		suppliers = append(suppliers, supplier)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return suppliers, nil
}

// ListActive lists all active suppliers
func (r *SupplierRepository) ListActive(limit, offset int) ([]model.Supplier, error) {
	query := `
		SELECT id, name, email, phone, address, description, is_active, created_at, updated_at
		FROM suppliers
		WHERE is_active = true
		ORDER BY name
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []model.Supplier
	for rows.Next() {
		var supplier model.Supplier
		err := rows.Scan(
			&supplier.ID,
			&supplier.Name,
			&supplier.Email,
			&supplier.Phone,
			&supplier.Address,
			&supplier.Description,
			&supplier.IsActive,
			&supplier.CreatedAt,
			&supplier.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		suppliers = append(suppliers, supplier)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return suppliers, nil
}
