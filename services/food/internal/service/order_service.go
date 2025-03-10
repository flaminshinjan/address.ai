package service

import (
	"errors"

	"github.com/yourusername/hotel-management/services/food/internal/model"
	"github.com/yourusername/hotel-management/services/food/internal/repository"
)

// OrderService handles business logic for food orders
type OrderService struct {
	orderRepo *repository.OrderRepository
	menuRepo  *repository.MenuRepository
}

// NewOrderService creates a new OrderService
func NewOrderService(orderRepo *repository.OrderRepository, menuRepo *repository.MenuRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		menuRepo:  menuRepo,
	}
}

// CreateOrder creates a new food order
func (s *OrderService) CreateOrder(userID string, req model.CreateOrderRequest) (model.FoodOrderResponse, error) {
	// Validate request
	if len(req.Items) == 0 {
		return model.FoodOrderResponse{}, errors.New("at least one item is required")
	}

	// Create order
	order := model.FoodOrder{
		UserID:     userID,
		RoomID:     req.RoomID,
		Status:     "pending",
		TotalPrice: 0,
		Notes:      req.Notes,
	}

	createdOrder, err := s.orderRepo.CreateOrder(order)
	if err != nil {
		return model.FoodOrderResponse{}, err
	}

	// Create order items
	var totalPrice float64
	var orderItems []model.OrderItemResponse

	for _, itemReq := range req.Items {
		// Get menu item
		menuItem, err := s.menuRepo.GetMenuItemByID(itemReq.MenuItemID)
		if err != nil {
			return model.FoodOrderResponse{}, err
		}

		// Check if menu item is available
		if !menuItem.IsAvailable {
			return model.FoodOrderResponse{}, errors.New("menu item is not available: " + menuItem.Name)
		}

		// Calculate price
		itemPrice := menuItem.Price * float64(itemReq.Quantity)
		totalPrice += itemPrice

		// Create order item
		orderItem := model.OrderItem{
			OrderID:    createdOrder.ID,
			MenuItemID: itemReq.MenuItemID,
			Quantity:   itemReq.Quantity,
			Price:      itemPrice,
			Notes:      itemReq.Notes,
		}

		createdItem, err := s.orderRepo.CreateOrderItem(orderItem)
		if err != nil {
			return model.FoodOrderResponse{}, err
		}

		// Add to response
		orderItems = append(orderItems, model.OrderItemResponse{
			ID:       createdItem.ID,
			MenuItem: menuItem.ToResponse(),
			Quantity: createdItem.Quantity,
			Price:    createdItem.Price,
			Notes:    createdItem.Notes,
		})
	}

	// Update order total price
	createdOrder.TotalPrice = totalPrice
	_, err = s.orderRepo.CreateOrder(createdOrder)
	if err != nil {
		return model.FoodOrderResponse{}, err
	}

	// Create response
	response := model.FoodOrderResponse{
		ID:         createdOrder.ID,
		UserID:     createdOrder.UserID,
		RoomID:     createdOrder.RoomID,
		Status:     createdOrder.Status,
		TotalPrice: totalPrice,
		Notes:      createdOrder.Notes,
		Items:      orderItems,
		CreatedAt:  createdOrder.CreatedAt,
	}

	return response, nil
}

// GetOrderByID gets a food order by ID
func (s *OrderService) GetOrderByID(id string) (model.FoodOrderResponse, error) {
	// Get order
	order, err := s.orderRepo.GetOrderByID(id)
	if err != nil {
		return model.FoodOrderResponse{}, err
	}

	// Get order items
	items, err := s.orderRepo.GetOrderItemsByOrderID(order.ID)
	if err != nil {
		return model.FoodOrderResponse{}, err
	}

	// Create response
	var orderItems []model.OrderItemResponse
	for _, item := range items {
		// Get menu item
		menuItem, err := s.menuRepo.GetMenuItemByID(item.MenuItemID)
		if err != nil {
			return model.FoodOrderResponse{}, err
		}

		orderItems = append(orderItems, model.OrderItemResponse{
			ID:       item.ID,
			MenuItem: menuItem.ToResponse(),
			Quantity: item.Quantity,
			Price:    item.Price,
			Notes:    item.Notes,
		})
	}

	response := model.FoodOrderResponse{
		ID:         order.ID,
		UserID:     order.UserID,
		RoomID:     order.RoomID,
		Status:     order.Status,
		TotalPrice: order.TotalPrice,
		Notes:      order.Notes,
		Items:      orderItems,
		CreatedAt:  order.CreatedAt,
	}

	return response, nil
}

// UpdateOrderStatus updates a food order's status
func (s *OrderService) UpdateOrderStatus(id, status string) error {
	// Validate status
	validStatuses := map[string]bool{
		"pending":   true,
		"preparing": true,
		"delivered": true,
		"cancelled": true,
	}

	if !validStatuses[status] {
		return errors.New("invalid status")
	}

	// Check if order exists
	_, err := s.orderRepo.GetOrderByID(id)
	if err != nil {
		return err
	}

	return s.orderRepo.UpdateOrderStatus(id, status)
}

// ListOrdersByUserID lists food orders by user ID
func (s *OrderService) ListOrdersByUserID(userID string, limit, offset int) ([]model.FoodOrderResponse, error) {
	// Get orders
	orders, err := s.orderRepo.ListOrdersByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Create responses
	var responses []model.FoodOrderResponse
	for _, order := range orders {
		// Get order items
		items, err := s.orderRepo.GetOrderItemsByOrderID(order.ID)
		if err != nil {
			return nil, err
		}

		// Create response
		var orderItems []model.OrderItemResponse
		for _, item := range items {
			// Get menu item
			menuItem, err := s.menuRepo.GetMenuItemByID(item.MenuItemID)
			if err != nil {
				return nil, err
			}

			orderItems = append(orderItems, model.OrderItemResponse{
				ID:       item.ID,
				MenuItem: menuItem.ToResponse(),
				Quantity: item.Quantity,
				Price:    item.Price,
				Notes:    item.Notes,
			})
		}

		response := model.FoodOrderResponse{
			ID:         order.ID,
			UserID:     order.UserID,
			RoomID:     order.RoomID,
			Status:     order.Status,
			TotalPrice: order.TotalPrice,
			Notes:      order.Notes,
			Items:      orderItems,
			CreatedAt:  order.CreatedAt,
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// ListOrders lists all food orders
func (s *OrderService) ListOrders(limit, offset int) ([]model.FoodOrderResponse, error) {
	// Get orders
	orders, err := s.orderRepo.ListOrders(limit, offset)
	if err != nil {
		return nil, err
	}

	// Create responses
	var responses []model.FoodOrderResponse
	for _, order := range orders {
		// Get order items
		items, err := s.orderRepo.GetOrderItemsByOrderID(order.ID)
		if err != nil {
			return nil, err
		}

		// Create response
		var orderItems []model.OrderItemResponse
		for _, item := range items {
			// Get menu item
			menuItem, err := s.menuRepo.GetMenuItemByID(item.MenuItemID)
			if err != nil {
				return nil, err
			}

			orderItems = append(orderItems, model.OrderItemResponse{
				ID:       item.ID,
				MenuItem: menuItem.ToResponse(),
				Quantity: item.Quantity,
				Price:    item.Price,
				Notes:    item.Notes,
			})
		}

		response := model.FoodOrderResponse{
			ID:         order.ID,
			UserID:     order.UserID,
			RoomID:     order.RoomID,
			Status:     order.Status,
			TotalPrice: order.TotalPrice,
			Notes:      order.Notes,
			Items:      orderItems,
			CreatedAt:  order.CreatedAt,
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// ListOrdersByStatus lists food orders by status
func (s *OrderService) ListOrdersByStatus(status string, limit, offset int) ([]model.FoodOrderResponse, error) {
	// Validate status
	validStatuses := map[string]bool{
		"pending":   true,
		"preparing": true,
		"delivered": true,
		"cancelled": true,
	}

	if !validStatuses[status] {
		return nil, errors.New("invalid status")
	}

	// Get orders
	orders, err := s.orderRepo.ListOrdersByStatus(status, limit, offset)
	if err != nil {
		return nil, err
	}

	// Create responses
	var responses []model.FoodOrderResponse
	for _, order := range orders {
		// Get order items
		items, err := s.orderRepo.GetOrderItemsByOrderID(order.ID)
		if err != nil {
			return nil, err
		}

		// Create response
		var orderItems []model.OrderItemResponse
		for _, item := range items {
			// Get menu item
			menuItem, err := s.menuRepo.GetMenuItemByID(item.MenuItemID)
			if err != nil {
				return nil, err
			}

			orderItems = append(orderItems, model.OrderItemResponse{
				ID:       item.ID,
				MenuItem: menuItem.ToResponse(),
				Quantity: item.Quantity,
				Price:    item.Price,
				Notes:    item.Notes,
			})
		}

		response := model.FoodOrderResponse{
			ID:         order.ID,
			UserID:     order.UserID,
			RoomID:     order.RoomID,
			Status:     order.Status,
			TotalPrice: order.TotalPrice,
			Notes:      order.Notes,
			Items:      orderItems,
			CreatedAt:  order.CreatedAt,
		}

		responses = append(responses, response)
	}

	return responses, nil
}
