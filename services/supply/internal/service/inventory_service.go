package service

import (
	"errors"
	"time"

	"github.com/flaminshinjan/address.ai/services/supply/internal/model"
	"github.com/flaminshinjan/address.ai/services/supply/internal/repository"
)

// InventoryService handles business logic for inventory items
type InventoryService struct {
	inventoryRepo *repository.InventoryRepository
}

// NewInventoryService creates a new InventoryService
func NewInventoryService(inventoryRepo *repository.InventoryRepository) *InventoryService {
	return &InventoryService{
		inventoryRepo: inventoryRepo,
	}
}

// CreateInventoryItem creates a new inventory item
func (s *InventoryService) CreateInventoryItem(item model.InventoryItem) (model.InventoryItemResponse, error) {
	// Validate item
	if item.Name == "" {
		return model.InventoryItemResponse{}, errors.New("name is required")
	}

	if item.Category == "" {
		return model.InventoryItemResponse{}, errors.New("category is required")
	}

	if item.Unit == "" {
		return model.InventoryItemResponse{}, errors.New("unit is required")
	}

	if item.MinQuantity < 0 {
		return model.InventoryItemResponse{}, errors.New("min quantity must be non-negative")
	}

	if item.Price < 0 {
		return model.InventoryItemResponse{}, errors.New("price must be non-negative")
	}

	// Create item
	createdItem, err := s.inventoryRepo.Create(item)
	if err != nil {
		return model.InventoryItemResponse{}, err
	}

	return createdItem.ToResponse(), nil
}

// GetInventoryItemByID gets an inventory item by ID
func (s *InventoryService) GetInventoryItemByID(id string) (model.InventoryItemResponse, error) {
	item, err := s.inventoryRepo.GetByID(id)
	if err != nil {
		return model.InventoryItemResponse{}, err
	}

	return item.ToResponse(), nil
}

// UpdateInventoryItem updates an inventory item
func (s *InventoryService) UpdateInventoryItem(id string, item model.InventoryItem) (model.InventoryItemResponse, error) {
	// Check if item exists
	existingItem, err := s.inventoryRepo.GetByID(id)
	if err != nil {
		return model.InventoryItemResponse{}, err
	}

	// Update fields
	existingItem.Name = item.Name
	existingItem.Category = item.Category
	existingItem.Description = item.Description
	existingItem.Unit = item.Unit
	existingItem.MinQuantity = item.MinQuantity
	existingItem.Price = item.Price
	existingItem.UpdatedAt = time.Now()

	// Save item
	updatedItem, err := s.inventoryRepo.Update(existingItem)
	if err != nil {
		return model.InventoryItemResponse{}, err
	}

	return updatedItem.ToResponse(), nil
}

// UpdateInventoryQuantity updates an inventory item's quantity
func (s *InventoryService) UpdateInventoryQuantity(id string, quantity int, userID, notes string) error {
	// Check if item exists
	existingItem, err := s.inventoryRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Calculate quantity change
	change := quantity - existingItem.Quantity

	// Create transaction
	transactionType := "in"
	if change < 0 {
		transactionType = "out"
		change = -change // Make positive for the transaction
	}

	transaction := model.InventoryTransaction{
		InventoryItemID: id,
		Quantity:        change,
		Type:            transactionType,
		Source:          "adjustment",
		Notes:           notes,
		CreatedBy:       userID,
	}

	// Create transaction
	_, err = s.inventoryRepo.CreateTransaction(transaction)
	if err != nil {
		return err
	}

	// Update quantity
	return s.inventoryRepo.UpdateQuantity(id, quantity)
}

// DeleteInventoryItem deletes an inventory item
func (s *InventoryService) DeleteInventoryItem(id string) error {
	return s.inventoryRepo.Delete(id)
}

// ListInventoryItems lists all inventory items
func (s *InventoryService) ListInventoryItems(limit, offset int) ([]model.InventoryItemResponse, error) {
	items, err := s.inventoryRepo.List(limit, offset)
	if err != nil {
		return nil, err
	}

	var itemResponses []model.InventoryItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, item.ToResponse())
	}

	return itemResponses, nil
}

// ListInventoryItemsByCategory lists inventory items by category
func (s *InventoryService) ListInventoryItemsByCategory(category string, limit, offset int) ([]model.InventoryItemResponse, error) {
	items, err := s.inventoryRepo.ListByCategory(category, limit, offset)
	if err != nil {
		return nil, err
	}

	var itemResponses []model.InventoryItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, item.ToResponse())
	}

	return itemResponses, nil
}

// ListLowStockItems lists inventory items with quantity below min_quantity
func (s *InventoryService) ListLowStockItems(limit, offset int) ([]model.InventoryItemResponse, error) {
	items, err := s.inventoryRepo.ListLowStock(limit, offset)
	if err != nil {
		return nil, err
	}

	var itemResponses []model.InventoryItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, item.ToResponse())
	}

	return itemResponses, nil
}

// ListCategories lists all inventory categories
func (s *InventoryService) ListCategories() ([]string, error) {
	return s.inventoryRepo.ListCategories()
}
