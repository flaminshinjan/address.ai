package handler

import (
	"net/http"
	"strconv"

	"github.com/flaminshinjan/address.ai/pkg/common/auth"
	"github.com/flaminshinjan/address.ai/services/supply/internal/model"
	"github.com/flaminshinjan/address.ai/services/supply/internal/service"
	"github.com/labstack/echo/v4"
)

// PurchaseHandler handles HTTP requests for purchase orders
type PurchaseHandler struct {
	service   *service.PurchaseService
	jwtSecret string
}

// NewPurchaseHandler creates a new PurchaseHandler
func NewPurchaseHandler(service *service.PurchaseService, jwtSecret string) *PurchaseHandler {
	return &PurchaseHandler{
		service:   service,
		jwtSecret: jwtSecret,
	}
}

// RegisterRoutes registers the routes for the purchase handler
func (h *PurchaseHandler) RegisterRoutes(g *echo.Group) {
	// All purchase order routes require admin access
	admin := g.Group("/purchase-orders")
	admin.Use(h.authMiddleware)

	admin.POST("", h.CreatePurchaseOrder)
	admin.GET("", h.ListPurchaseOrders)
	admin.GET("/status/:status", h.ListPurchaseOrdersByStatus)
	admin.GET("/:id", h.GetPurchaseOrder)
	admin.PUT("/:id/status", h.UpdatePurchaseOrderStatus)
}

// CreatePurchaseOrder handles creating a new purchase order
func (h *PurchaseHandler) CreatePurchaseOrder(c echo.Context) error {
	// Get user ID from context
	userID := c.Get("user_id").(string)

	var req model.CreatePurchaseOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request payload",
		})
	}

	order, err := h.service.CreatePurchaseOrder(userID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Purchase order created successfully",
		"data":    order,
	})
}

// GetPurchaseOrder handles getting a purchase order by ID
func (h *PurchaseHandler) GetPurchaseOrder(c echo.Context) error {
	id := c.Param("id")

	order, err := h.service.GetPurchaseOrderByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Purchase order not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Purchase order retrieved successfully",
		"data":    order,
	})
}

// UpdatePurchaseOrderStatus handles updating a purchase order's status
func (h *PurchaseHandler) UpdatePurchaseOrderStatus(c echo.Context) error {
	// Get user ID from context
	userID := c.Get("user_id").(string)
	id := c.Param("id")

	var statusData struct {
		Status string `json:"status"`
	}

	if err := c.Bind(&statusData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request payload",
		})
	}

	if statusData.Status == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Status is required",
		})
	}

	if err := h.service.UpdatePurchaseOrderStatus(id, statusData.Status, userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Purchase order status updated successfully",
	})
}

// ListPurchaseOrders handles listing all purchase orders
func (h *PurchaseHandler) ListPurchaseOrders(c echo.Context) error {
	// Get query parameters
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

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

	orders, err := h.service.ListPurchaseOrders(limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to retrieve purchase orders",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Purchase orders retrieved successfully",
		"data":    orders,
	})
}

// ListPurchaseOrdersByStatus handles listing purchase orders by status
func (h *PurchaseHandler) ListPurchaseOrdersByStatus(c echo.Context) error {
	status := c.Param("status")

	// Get query parameters
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

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

	orders, err := h.service.ListPurchaseOrdersByStatus(status, limit, offset)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Purchase orders retrieved successfully",
		"data":    orders,
	})
}

// authMiddleware is a middleware to check if the user is authenticated and is an admin
func (h *PurchaseHandler) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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

		// Check if user is admin
		if claims.Role != "admin" {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"success": false,
				"error":   "Admin access required",
			})
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		return next(c)
	}
}
