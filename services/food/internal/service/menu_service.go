package service

import (
	"errors"
	"time"

	"github.com/flaminshinjan/address.ai/services/food/internal/model"
	"github.com/flaminshinjan/address.ai/services/food/internal/repository"
)

// MenuService handles business logic for menu items
type MenuService struct {
	menuRepo *repository.MenuRepository
}

// NewMenuService creates a new MenuService
func NewMenuService(menuRepo *repository.MenuRepository) *MenuService {
	return &MenuService{
		menuRepo: menuRepo,
	}
}

// CreateMenuItem creates a new menu item
func (s *MenuService) CreateMenuItem(item model.MenuItem) (model.MenuItemResponse, error) {
	// Validate item
	if item.Name == "" {
		return model.MenuItemResponse{}, errors.New("name is required")
	}

	if item.Category == "" {
		return model.MenuItemResponse{}, errors.New("category is required")
	}

	if item.Price <= 0 {
		return model.MenuItemResponse{}, errors.New("price must be greater than zero")
	}

	// Create item
	createdItem, err := s.menuRepo.CreateMenuItem(item)
	if err != nil {
		return model.MenuItemResponse{}, err
	}

	return createdItem.ToResponse(), nil
}

// GetMenuItemByID gets a menu item by ID
func (s *MenuService) GetMenuItemByID(id string) (model.MenuItemResponse, error) {
	item, err := s.menuRepo.GetMenuItemByID(id)
	if err != nil {
		return model.MenuItemResponse{}, err
	}

	return item.ToResponse(), nil
}

// UpdateMenuItem updates a menu item
func (s *MenuService) UpdateMenuItem(id string, item model.MenuItem) (model.MenuItemResponse, error) {
	// Check if item exists
	existingItem, err := s.menuRepo.GetMenuItemByID(id)
	if err != nil {
		return model.MenuItemResponse{}, err
	}

	// Update fields
	existingItem.Name = item.Name
	existingItem.Description = item.Description
	existingItem.Category = item.Category
	existingItem.Price = item.Price
	existingItem.IsAvailable = item.IsAvailable
	existingItem.UpdatedAt = time.Now()

	// Save item
	updatedItem, err := s.menuRepo.UpdateMenuItem(existingItem)
	if err != nil {
		return model.MenuItemResponse{}, err
	}

	return updatedItem.ToResponse(), nil
}

// DeleteMenuItem deletes a menu item
func (s *MenuService) DeleteMenuItem(id string) error {
	return s.menuRepo.DeleteMenuItem(id)
}

// ListMenuItems lists all menu items
func (s *MenuService) ListMenuItems(limit, offset int) ([]model.MenuItemResponse, error) {
	items, err := s.menuRepo.ListMenuItems(limit, offset)
	if err != nil {
		return nil, err
	}

	var itemResponses []model.MenuItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, item.ToResponse())
	}

	return itemResponses, nil
}

// ListMenuItemsByCategory lists menu items by category
func (s *MenuService) ListMenuItemsByCategory(category string, limit, offset int) ([]model.MenuItemResponse, error) {
	items, err := s.menuRepo.ListMenuItemsByCategory(category, limit, offset)
	if err != nil {
		return nil, err
	}

	var itemResponses []model.MenuItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, item.ToResponse())
	}

	return itemResponses, nil
}

// ListAvailableMenuItems lists all available menu items
func (s *MenuService) ListAvailableMenuItems(limit, offset int) ([]model.MenuItemResponse, error) {
	items, err := s.menuRepo.ListAvailableMenuItems(limit, offset)
	if err != nil {
		return nil, err
	}

	var itemResponses []model.MenuItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, item.ToResponse())
	}

	return itemResponses, nil
}

// ListCategories lists all menu categories
func (s *MenuService) ListCategories() ([]string, error) {
	return s.menuRepo.ListCategories()
}
