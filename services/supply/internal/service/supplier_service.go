package service

import (
	"errors"
	"time"

	"github.com/flaminshinjan/address.ai/services/supply/internal/model"
	"github.com/flaminshinjan/address.ai/services/supply/internal/repository"
)

// SupplierService handles business logic for suppliers
type SupplierService struct {
	supplierRepo *repository.SupplierRepository
}

// NewSupplierService creates a new SupplierService
func NewSupplierService(supplierRepo *repository.SupplierRepository) *SupplierService {
	return &SupplierService{
		supplierRepo: supplierRepo,
	}
}

// CreateSupplier creates a new supplier
func (s *SupplierService) CreateSupplier(supplier model.Supplier) (model.SupplierResponse, error) {
	// Validate supplier
	if supplier.Name == "" {
		return model.SupplierResponse{}, errors.New("name is required")
	}

	if supplier.Email == "" {
		return model.SupplierResponse{}, errors.New("email is required")
	}

	if supplier.Phone == "" {
		return model.SupplierResponse{}, errors.New("phone is required")
	}

	if supplier.Address == "" {
		return model.SupplierResponse{}, errors.New("address is required")
	}

	// Create supplier
	createdSupplier, err := s.supplierRepo.Create(supplier)
	if err != nil {
		return model.SupplierResponse{}, err
	}

	return createdSupplier.ToResponse(), nil
}

// GetSupplierByID gets a supplier by ID
func (s *SupplierService) GetSupplierByID(id string) (model.SupplierResponse, error) {
	supplier, err := s.supplierRepo.GetByID(id)
	if err != nil {
		return model.SupplierResponse{}, err
	}

	return supplier.ToResponse(), nil
}

// UpdateSupplier updates a supplier
func (s *SupplierService) UpdateSupplier(id string, supplier model.Supplier) (model.SupplierResponse, error) {
	// Check if supplier exists
	existingSupplier, err := s.supplierRepo.GetByID(id)
	if err != nil {
		return model.SupplierResponse{}, err
	}

	// Update fields
	existingSupplier.Name = supplier.Name
	existingSupplier.Email = supplier.Email
	existingSupplier.Phone = supplier.Phone
	existingSupplier.Address = supplier.Address
	existingSupplier.Description = supplier.Description
	existingSupplier.IsActive = supplier.IsActive
	existingSupplier.UpdatedAt = time.Now()

	// Save supplier
	updatedSupplier, err := s.supplierRepo.Update(existingSupplier)
	if err != nil {
		return model.SupplierResponse{}, err
	}

	return updatedSupplier.ToResponse(), nil
}

// DeleteSupplier deletes a supplier
func (s *SupplierService) DeleteSupplier(id string) error {
	return s.supplierRepo.Delete(id)
}

// ListSuppliers lists all suppliers
func (s *SupplierService) ListSuppliers(limit, offset int) ([]model.SupplierResponse, error) {
	suppliers, err := s.supplierRepo.List(limit, offset)
	if err != nil {
		return nil, err
	}

	var supplierResponses []model.SupplierResponse
	for _, supplier := range suppliers {
		supplierResponses = append(supplierResponses, supplier.ToResponse())
	}

	return supplierResponses, nil
}

// ListActiveSuppliers lists all active suppliers
func (s *SupplierService) ListActiveSuppliers(limit, offset int) ([]model.SupplierResponse, error) {
	suppliers, err := s.supplierRepo.ListActive(limit, offset)
	if err != nil {
		return nil, err
	}

	var supplierResponses []model.SupplierResponse
	for _, supplier := range suppliers {
		supplierResponses = append(supplierResponses, supplier.ToResponse())
	}

	return supplierResponses, nil
}
