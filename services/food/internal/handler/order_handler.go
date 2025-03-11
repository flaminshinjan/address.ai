package handler

import (
	"net/http"
	"strconv"

	"github.com/flaminshinjan/address.ai/pkg/common/auth"
	"github.com/flaminshinjan/address.ai/services/food/internal/model"
	"github.com/flaminshinjan/address.ai/services/food/internal/service"
	"github.com/labstack/echo/v4"
)

// OrderHandler handles HTTP requests for food orders
type OrderHandler struct {
	service   *service.OrderService
	jwtSecret string
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(service *service.OrderService, jwtSecret string) *OrderHandler {
	return &OrderHandler{
		service:   service,
		jwtSecret: jwtSecret,
	}
}

// CreateOrder handles creating a new food order
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	// Get user ID from context
	userID := c.Get("user_id").(string)

	var order model.CreateOrderRequest
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request payload",
		})
	}

	createdOrder, err := h.service.CreateOrder(userID, order)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Order created successfully",
		"data":    createdOrder,
	})
}

// GetOrder handles getting an order by ID
func (h *OrderHandler) GetOrder(c echo.Context) error {
	// Get user ID and role from context
	userID := c.Get("user_id").(string)
	role := c.Get("role").(string)

	id := c.Param("id")

	order, err := h.service.GetOrderByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Order not found",
		})
	}

	// Check if user is authorized to view this order
	if role != "admin" && order.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"success": false,
			"error":   "You are not authorized to view this order",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Order retrieved successfully",
		"data":    order,
	})
}

// UpdateOrderStatus handles updating an order's status
func (h *OrderHandler) UpdateOrderStatus(c echo.Context) error {
	// Check if user is admin or staff
	role := c.Get("role").(string)
	if role != "admin" && role != "staff" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"success": false,
			"error":   "Admin or staff access required",
		})
	}

	id := c.Param("id")

	var statusUpdate struct {
		Status string `json:"status"`
	}
	if err := c.Bind(&statusUpdate); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request payload",
		})
	}

	if statusUpdate.Status == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Status is required",
		})
	}

	err := h.service.UpdateOrderStatus(id, statusUpdate.Status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Get updated order
	updatedOrder, err := h.service.GetOrderByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to retrieve updated order",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Order status updated successfully",
		"data":    updatedOrder,
	})
}

// ListOrders handles listing all orders
func (h *OrderHandler) ListOrders(c echo.Context) error {
	// Check if user is admin or staff
	role := c.Get("role").(string)
	if role != "admin" && role != "staff" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"success": false,
			"error":   "Admin or staff access required",
		})
	}

	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")
	status := c.QueryParam("status")

	limit := 10 // Default limit
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0 // Default offset
	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	var orders interface{}
	var err error

	if status != "" {
		orders, err = h.service.ListOrdersByStatus(status, limit, offset)
	} else {
		orders, err = h.service.ListOrders(limit, offset)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to retrieve orders",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Orders retrieved successfully",
		"data":    orders,
	})
}

// ListUserOrders handles listing orders for a specific user
func (h *OrderHandler) ListUserOrders(c echo.Context) error {
	// Get user ID from context
	userID := c.Get("user_id").(string)

	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")
	status := c.QueryParam("status")

	limit := 10 // Default limit
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0 // Default offset
	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	var orders interface{}
	var err error

	if status != "" {
		// For simplicity, we'll just use the user ID filter and ignore status for now
		orders, err = h.service.ListOrdersByUserID(userID, limit, offset)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error":   "Failed to retrieve orders",
			})
		}
	} else {
		orders, err = h.service.ListOrdersByUserID(userID, limit, offset)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error":   "Failed to retrieve orders",
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Orders retrieved successfully",
		"data":    orders,
	})
}

// RegisterRoutes registers the routes for the order handler
func (h *OrderHandler) RegisterRoutes(g *echo.Group) {
	// Order routes
	orders := g.Group("/orders")
	orders.Use(h.authMiddleware)

	// User routes
	orders.POST("", h.CreateOrder)
	orders.GET("/my-orders", h.ListUserOrders)
	orders.GET("/:id", h.GetOrder)

	// Admin/Staff routes
	adminStaff := orders.Group("")
	adminStaff.Use(h.adminStaffMiddleware)

	adminStaff.GET("", h.ListOrders)
	adminStaff.PUT("/:id/status", h.UpdateOrderStatus)
}

// authMiddleware is a middleware to check if the user is authenticated
func (h *OrderHandler) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"error":   "Authorization token is required",
			})
		}

		// Remove "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		claims, err := auth.ValidateToken(token, h.jwtSecret)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"error":   "Invalid or expired token",
			})
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		return next(c)
	}
}

// adminStaffMiddleware is a middleware to check if the user is an admin or staff
func (h *OrderHandler) adminStaffMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := c.Get("role").(string)
		if role != "admin" && role != "staff" {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"success": false,
				"error":   "Admin or staff access required",
			})
		}
		return next(c)
	}
}
