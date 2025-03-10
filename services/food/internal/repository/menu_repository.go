package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/hotel-management/services/food/internal/model"
)

// MenuRepository handles database operations for menu items
type MenuRepository struct {
	db *sql.DB
}

// NewMenuRepository creates a new MenuRepository
func NewMenuRepository(db *sql.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

// CreateMenuItem creates a new menu item
func (r *MenuRepository) CreateMenuItem(item model.MenuItem) (model.MenuItem, error) {
	query := `
		INSERT INTO menu_items (id, name, description, category, price, is_available, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, name, description, category, price, is_available, created_at, updated_at
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
		item.Description,
		item.Category,
		item.Price,
		item.IsAvailable,
		item.CreatedAt,
		item.UpdatedAt,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Category,
		&item.Price,
		&item.IsAvailable,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return model.MenuItem{}, err
	}

	return item, nil
}

// GetMenuItemByID gets a menu item by ID
func (r *MenuRepository) GetMenuItemByID(id string) (model.MenuItem, error) {
	query := `
		SELECT id, name, description, category, price, is_available, created_at, updated_at
		FROM menu_items
		WHERE id = $1
	`

	var item model.MenuItem
	err := r.db.QueryRow(query, id).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Category,
		&item.Price,
		&item.IsAvailable,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.MenuItem{}, errors.New("menu item not found")
		}
		return model.MenuItem{}, err
	}

	return item, nil
}

// UpdateMenuItem updates a menu item
func (r *MenuRepository) UpdateMenuItem(item model.MenuItem) (model.MenuItem, error) {
	query := `
		UPDATE menu_items
		SET name = $1, description = $2, category = $3, price = $4, is_available = $5, updated_at = $6
		WHERE id = $7
		RETURNING id, name, description, category, price, is_available, created_at, updated_at
	`

	item.UpdatedAt = time.Now()

	err := r.db.QueryRow(
		query,
		item.Name,
		item.Description,
		item.Category,
		item.Price,
		item.IsAvailable,
		item.UpdatedAt,
		item.ID,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Category,
		&item.Price,
		&item.IsAvailable,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.MenuItem{}, errors.New("menu item not found")
		}
		return model.MenuItem{}, err
	}

	return item, nil
}

// DeleteMenuItem deletes a menu item
func (r *MenuRepository) DeleteMenuItem(id string) error {
	query := `
		DELETE FROM menu_items
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
		return errors.New("menu item not found")
	}

	return nil
}

// ListMenuItems lists all menu items
func (r *MenuRepository) ListMenuItems(limit, offset int) ([]model.MenuItem, error) {
	query := `
		SELECT id, name, description, category, price, is_available, created_at, updated_at
		FROM menu_items
		ORDER BY category, name
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.MenuItem
	for rows.Next() {
		var item model.MenuItem
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Category,
			&item.Price,
			&item.IsAvailable,
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

// ListMenuItemsByCategory lists menu items by category
func (r *MenuRepository) ListMenuItemsByCategory(category string, limit, offset int) ([]model.MenuItem, error) {
	query := `
		SELECT id, name, description, category, price, is_available, created_at, updated_at
		FROM menu_items
		WHERE category = $1
		ORDER BY name
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, category, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.MenuItem
	for rows.Next() {
		var item model.MenuItem
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Category,
			&item.Price,
			&item.IsAvailable,
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

// ListAvailableMenuItems lists all available menu items
func (r *MenuRepository) ListAvailableMenuItems(limit, offset int) ([]model.MenuItem, error) {
	query := `
		SELECT id, name, description, category, price, is_available, created_at, updated_at
		FROM menu_items
		WHERE is_available = true
		ORDER BY category, name
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.MenuItem
	for rows.Next() {
		var item model.MenuItem
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Category,
			&item.Price,
			&item.IsAvailable,
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

// ListCategories lists all menu categories
func (r *MenuRepository) ListCategories() ([]string, error) {
	query := `
		SELECT DISTINCT category
		FROM menu_items
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
