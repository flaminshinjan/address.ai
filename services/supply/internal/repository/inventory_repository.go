package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/flaminshinjan/address.ai/services/supply/internal/model"
)

// InventoryRepository handles database operations for inventory items
type InventoryRepository struct {
	db *sql.DB
}

// NewInventoryRepository creates a new InventoryRepository
func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

// Create creates a new inventory item
func (r *InventoryRepository) Create(item model.InventoryItem) (model.InventoryItem, error) {
	query := `
		INSERT INTO inventory_items (id, name, category, description, quantity, unit, min_quantity, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, name, category, description, quantity, unit, min_quantity, price, created_at, updated_at
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
		item.Name,
		item.Category,
		item.Description,
		item.Quantity,
		item.Unit,
		item.MinQuantity,
		item.Price,
		item.CreatedAt,
		item.UpdatedAt,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Category,
		&item.Description,
		&item.Quantity,
		&item.Unit,
		&item.MinQuantity,
		&item.Price,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return model.InventoryItem{}, err
	}

	return item, nil
}

// GetByID gets an inventory item by ID
func (r *InventoryRepository) GetByID(id string) (model.InventoryItem, error) {
	query := `
		SELECT id, name, category, description, quantity, unit, min_quantity, price, created_at, updated_at
		FROM inventory_items
		WHERE id = $1
	`

	var item model.InventoryItem
	err := r.db.QueryRow(query, id).Scan(
		&item.ID,
		&item.Name,
		&item.Category,
		&item.Description,
		&item.Quantity,
		&item.Unit,
		&item.MinQuantity,
		&item.Price,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.InventoryItem{}, errors.New("inventory item not found")
		}
		return model.InventoryItem{}, err
	}

	return item, nil
}

// Update updates an inventory item
func (r *InventoryRepository) Update(item model.InventoryItem) (model.InventoryItem, error) {
	query := `
		UPDATE inventory_items
		SET name = $1, category = $2, description = $3, quantity = $4, unit = $5, min_quantity = $6, price = $7, updated_at = $8
		WHERE id = $9
		RETURNING id, name, category, description, quantity, unit, min_quantity, price, created_at, updated_at
	`

	item.UpdatedAt = time.Now()

	err := r.db.QueryRow(
		query,
		item.Name,
		item.Category,
		item.Description,
		item.Quantity,
		item.Unit,
		item.MinQuantity,
		item.Price,
		item.UpdatedAt,
		item.ID,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Category,
		&item.Description,
		&item.Quantity,
		&item.Unit,
		&item.MinQuantity,
		&item.Price,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.InventoryItem{}, errors.New("inventory item not found")
		}
		return model.InventoryItem{}, err
	}

	return item, nil
}

// UpdateQuantity updates an inventory item's quantity
func (r *InventoryRepository) UpdateQuantity(id string, quantity int) error {
	query := `
		UPDATE inventory_items
		SET quantity = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(query, quantity, time.Now(), id)
	return err
}

// Delete deletes an inventory item
func (r *InventoryRepository) Delete(id string) error {
	query := `
		DELETE FROM inventory_items
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
		return errors.New("inventory item not found")
	}

	return nil
}

// List lists all inventory items
func (r *InventoryRepository) List(limit, offset int) ([]model.InventoryItem, error) {
	query := `
		SELECT id, name, category, description, quantity, unit, min_quantity, price, created_at, updated_at
		FROM inventory_items
		ORDER BY category, name
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.InventoryItem
	for rows.Next() {
		var item model.InventoryItem
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Category,
			&item.Description,
			&item.Quantity,
			&item.Unit,
			&item.MinQuantity,
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

// ListByCategory lists inventory items by category
func (r *InventoryRepository) ListByCategory(category string, limit, offset int) ([]model.InventoryItem, error) {
	query := `
		SELECT id, name, category, description, quantity, unit, min_quantity, price, created_at, updated_at
		FROM inventory_items
		WHERE category = $1
		ORDER BY name
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, category, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.InventoryItem
	for rows.Next() {
		var item model.InventoryItem
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Category,
			&item.Description,
			&item.Quantity,
			&item.Unit,
			&item.MinQuantity,
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

// ListLowStock lists inventory items with quantity below min_quantity
func (r *InventoryRepository) ListLowStock(limit, offset int) ([]model.InventoryItem, error) {
	query := `
		SELECT id, name, category, description, quantity, unit, min_quantity, price, created_at, updated_at
		FROM inventory_items
		WHERE quantity < min_quantity
		ORDER BY (quantity::float / min_quantity::float) ASC, category, name
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.InventoryItem
	for rows.Next() {
		var item model.InventoryItem
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Category,
			&item.Description,
			&item.Quantity,
			&item.Unit,
			&item.MinQuantity,
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

// ListCategories lists all inventory categories
func (r *InventoryRepository) ListCategories() ([]string, error) {
	query := `
		SELECT DISTINCT category
		FROM inventory_items
		ORDER BY category
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// CreateTransaction creates a new inventory transaction
func (r *InventoryRepository) CreateTransaction(transaction model.InventoryTransaction) (model.InventoryTransaction, error) {
	query := `
		INSERT INTO inventory_transactions (id, inventory_item_id, quantity, type, source, source_id, notes, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, inventory_item_id, quantity, type, source, source_id, notes, created_by, created_at
	`

	// Generate UUID if not provided
	if transaction.ID == "" {
		transaction.ID = uuid.New().String()
	}

	// Set timestamp
	transaction.CreatedAt = time.Now()

	err := r.db.QueryRow(
		query,
		transaction.ID,
		transaction.InventoryItemID,
		transaction.Quantity,
		transaction.Type,
		transaction.Source,
		transaction.SourceID,
		transaction.Notes,
		transaction.CreatedBy,
		transaction.CreatedAt,
	).Scan(
		&transaction.ID,
		&transaction.InventoryItemID,
		&transaction.Quantity,
		&transaction.Type,
		&transaction.Source,
		&transaction.SourceID,
		&transaction.Notes,
		&transaction.CreatedBy,
		&transaction.CreatedAt,
	)

	if err != nil {
		return model.InventoryTransaction{}, err
	}

	return transaction, nil
}
