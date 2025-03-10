package service

import (
	"errors"
	"time"

	"github.com/flaminshinjan/address.ai/services/supply/internal/model"
	"github.com/flaminshinjan/address.ai/services/supply/internal/repository"
)

// PurchaseService handles business logic for purchase orders
type PurchaseService struct {
	purchaseRepo  *repository.PurchaseRepository
	supplierRepo  *repository.SupplierRepository
	inventoryRepo *repository.InventoryRepository
}

// NewPurchaseService creates a new PurchaseService
func NewPurchaseService(purchaseRepo *repository.PurchaseRepository, supplierRepo *repository.SupplierRepository, inventoryRepo *repository.InventoryRepository) *PurchaseService {
	return &PurchaseService{
		purchaseRepo:  purchaseRepo,
		supplierRepo:  supplierRepo,
		inventoryRepo: inventoryRepo,
	}
}

// CreatePurchaseOrder creates a new purchase order
func (s *PurchaseService) CreatePurchaseOrder(userID string, req model.CreatePurchaseOrderRequest) (model.PurchaseOrderResponse, error) {
	// Validate request
	if len(req.Items) == 0 {
		return model.PurchaseOrderResponse{}, errors.New("at least one item is required")
	}

	// Check if supplier exists
	supplier, err := s.supplierRepo.GetByID(req.SupplierID)
	if err != nil {
		return model.PurchaseOrderResponse{}, err
	}

	// Create order
	order := model.PurchaseOrder{
		SupplierID: req.SupplierID,
		Status:     "pending",
		TotalPrice: 0,
		Notes:      req.Notes,
		OrderDate:  time.Now(),
		CreatedBy:  userID,
	}

	createdOrder, err := s.purchaseRepo.CreateOrder(order)
	if err != nil {
		return model.PurchaseOrderResponse{}, err
	}

	// Create order items
	var totalPrice float64
	var orderItems []model.OrderItemResponse

	for _, itemReq := range req.Items {
		// Get inventory item
		inventoryItem, err := s.inventoryRepo.GetByID(itemReq.InventoryItemID)
		if err != nil {
			return model.PurchaseOrderResponse{}, err
		}

		// Calculate price
		itemPrice := inventoryItem.Price * float64(itemReq.Quantity)
		totalPrice += itemPrice

		// Create order item
		orderItem := model.OrderItem{
			PurchaseOrderID: createdOrder.ID,
			InventoryItemID: itemReq.InventoryItemID,
			Quantity:        itemReq.Quantity,
			Price:           itemPrice,
		}

		createdItem, err := s.purchaseRepo.CreateOrderItem(orderItem)
		if err != nil {
			return model.PurchaseOrderResponse{}, err
		}

		// Add to response
		orderItems = append(orderItems, model.OrderItemResponse{
			ID:            createdItem.ID,
			InventoryItem: inventoryItem.ToResponse(),
			Quantity:      createdItem.Quantity,
			Price:         createdItem.Price,
		})
	}

	// Update order total price
	createdOrder.TotalPrice = totalPrice
	_, err = s.purchaseRepo.CreateOrder(createdOrder)
	if err != nil {
		return model.PurchaseOrderResponse{}, err
	}

	// Create response
	response := model.PurchaseOrderResponse{
		ID:           createdOrder.ID,
		Supplier:     supplier.ToResponse(),
		Status:       createdOrder.Status,
		TotalPrice:   totalPrice,
		Notes:        createdOrder.Notes,
		OrderDate:    createdOrder.OrderDate,
		DeliveryDate: createdOrder.DeliveryDate,
		CreatedBy:    createdOrder.CreatedBy,
		Items:        orderItems,
		CreatedAt:    createdOrder.CreatedAt,
	}

	return response, nil
}

// GetPurchaseOrderByID gets a purchase order by ID
func (s *PurchaseService) GetPurchaseOrderByID(id string) (model.PurchaseOrderResponse, error) {
	// Get order
	order, err := s.purchaseRepo.GetOrderByID(id)
	if err != nil {
		return model.PurchaseOrderResponse{}, err
	}

	// Get supplier
	supplier, err := s.supplierRepo.GetByID(order.SupplierID)
	if err != nil {
		return model.PurchaseOrderResponse{}, err
	}

	// Get order items
	items, err := s.purchaseRepo.GetOrderItemsByOrderID(order.ID)
	if err != nil {
		return model.PurchaseOrderResponse{}, err
	}

	// Create response
	var orderItems []model.OrderItemResponse
	for _, item := range items {
		// Get inventory item
		inventoryItem, err := s.inventoryRepo.GetByID(item.InventoryItemID)
		if err != nil {
			return model.PurchaseOrderResponse{}, err
		}

		orderItems = append(orderItems, model.OrderItemResponse{
			ID:            item.ID,
			InventoryItem: inventoryItem.ToResponse(),
			Quantity:      item.Quantity,
			Price:         item.Price,
		})
	}

	response := model.PurchaseOrderResponse{
		ID:           order.ID,
		Supplier:     supplier.ToResponse(),
		Status:       order.Status,
		TotalPrice:   order.TotalPrice,
		Notes:        order.Notes,
		OrderDate:    order.OrderDate,
		DeliveryDate: order.DeliveryDate,
		CreatedBy:    order.CreatedBy,
		Items:        orderItems,
		CreatedAt:    order.CreatedAt,
	}

	return response, nil
}

// UpdatePurchaseOrderStatus updates a purchase order's status
func (s *PurchaseService) UpdatePurchaseOrderStatus(id, status, userID string) error {
	// Validate status
	validStatuses := map[string]bool{
		"pending":   true,
		"approved":  true,
		"received":  true,
		"cancelled": true,
	}

	if !validStatuses[status] {
		return errors.New("invalid status")
	}

	// Get order
	order, err := s.purchaseRepo.GetOrderByID(id)
	if err != nil {
		return err
	}

	// Check if status transition is valid
	if order.Status == "cancelled" {
		return errors.New("cannot change status of a cancelled order")
	}

	if order.Status == "received" && status != "received" {
		return errors.New("cannot change status of a received order")
	}

	// If status is changing to received, update inventory
	if status == "received" && order.Status != "received" {
		// Get order items
		items, err := s.purchaseRepo.GetOrderItemsByOrderID(id)
		if err != nil {
			return err
		}

		// Update inventory for each item
		for _, item := range items {
			// Get inventory item
			inventoryItem, err := s.inventoryRepo.GetByID(item.InventoryItemID)
			if err != nil {
				return err
			}

			// Update quantity
			newQuantity := inventoryItem.Quantity + item.Quantity
			if err := s.inventoryRepo.UpdateQuantity(item.InventoryItemID, newQuantity); err != nil {
				return err
			}

			// Create transaction
			transaction := model.InventoryTransaction{
				InventoryItemID: item.InventoryItemID,
				Quantity:        item.Quantity,
				Type:            "in",
				Source:          "purchase_order",
				SourceID:        id,
				CreatedBy:       userID,
			}

			if _, err := s.inventoryRepo.CreateTransaction(transaction); err != nil {
				return err
			}
		}

		// Set delivery date
		if err := s.purchaseRepo.UpdateDeliveryDate(id, time.Now()); err != nil {
			return err
		}
	}

	// Update status
	return s.purchaseRepo.UpdateOrderStatus(id, status)
}

// ListPurchaseOrders lists all purchase orders
func (s *PurchaseService) ListPurchaseOrders(limit, offset int) ([]model.PurchaseOrderResponse, error) {
	// Get orders
	orders, err := s.purchaseRepo.ListOrders(limit, offset)
	if err != nil {
		return nil, err
	}

	// Create responses
	var responses []model.PurchaseOrderResponse
	for _, order := range orders {
		// Get supplier
		supplier, err := s.supplierRepo.GetByID(order.SupplierID)
		if err != nil {
			return nil, err
		}

		// Get order items
		items, err := s.purchaseRepo.GetOrderItemsByOrderID(order.ID)
		if err != nil {
			return nil, err
		}

		// Create response
		var orderItems []model.OrderItemResponse
		for _, item := range items {
			// Get inventory item
			inventoryItem, err := s.inventoryRepo.GetByID(item.InventoryItemID)
			if err != nil {
				return nil, err
			}

			orderItems = append(orderItems, model.OrderItemResponse{
				ID:            item.ID,
				InventoryItem: inventoryItem.ToResponse(),
				Quantity:      item.Quantity,
				Price:         item.Price,
			})
		}

		response := model.PurchaseOrderResponse{
			ID:           order.ID,
			Supplier:     supplier.ToResponse(),
			Status:       order.Status,
			TotalPrice:   order.TotalPrice,
			Notes:        order.Notes,
			OrderDate:    order.OrderDate,
			DeliveryDate: order.DeliveryDate,
			CreatedBy:    order.CreatedBy,
			Items:        orderItems,
			CreatedAt:    order.CreatedAt,
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// ListPurchaseOrdersByStatus lists purchase orders by status
func (s *PurchaseService) ListPurchaseOrdersByStatus(status string, limit, offset int) ([]model.PurchaseOrderResponse, error) {
	// Validate status
	validStatuses := map[string]bool{
		"pending":   true,
		"approved":  true,
		"received":  true,
		"cancelled": true,
	}

	if !validStatuses[status] {
		return nil, errors.New("invalid status")
	}

	// Get orders
	orders, err := s.purchaseRepo.ListOrdersByStatus(status, limit, offset)
	if err != nil {
		return nil, err
	}

	// Create responses
	var responses []model.PurchaseOrderResponse
	for _, order := range orders {
		// Get supplier
		supplier, err := s.supplierRepo.GetByID(order.SupplierID)
		if err != nil {
			return nil, err
		}

		// Get order items
		items, err := s.purchaseRepo.GetOrderItemsByOrderID(order.ID)
		if err != nil {
			return nil, err
		}

		// Create response
		var orderItems []model.OrderItemResponse
		for _, item := range items {
			// Get inventory item
			inventoryItem, err := s.inventoryRepo.GetByID(item.InventoryItemID)
			if err != nil {
				return nil, err
			}

			orderItems = append(orderItems, model.OrderItemResponse{
				ID:            item.ID,
				InventoryItem: inventoryItem.ToResponse(),
				Quantity:      item.Quantity,
				Price:         item.Price,
			})
		}

		response := model.PurchaseOrderResponse{
			ID:           order.ID,
			Supplier:     supplier.ToResponse(),
			Status:       order.Status,
			TotalPrice:   order.TotalPrice,
			Notes:        order.Notes,
			OrderDate:    order.OrderDate,
			DeliveryDate: order.DeliveryDate,
			CreatedBy:    order.CreatedBy,
			Items:        orderItems,
			CreatedAt:    order.CreatedAt,
		}

		responses = append(responses, response)
	}

	return responses, nil
}
