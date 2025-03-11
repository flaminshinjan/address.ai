package handler

import (
	"net/http"
	"strconv"

	"github.com/flaminshinjan/address.ai/pkg/common/auth"
	"github.com/flaminshinjan/address.ai/services/supply/internal/model"
	"github.com/flaminshinjan/address.ai/services/supply/internal/service"
	"github.com/labstack/echo/v4"
)

// InventoryHandler handles HTTP requests for inventory items
type InventoryHandler struct {
	service   *service.InventoryService
	jwtSecret string
}

// NewInventoryHandler creates a new InventoryHandler
func NewInventoryHandler(service *service.InventoryService, jwtSecret string) *InventoryHandler {
	return &InventoryHandler{
		service:   service,
		jwtSecret: jwtSecret,
	}
}

// RegisterRoutes registers the routes for the inventory handler
func (h *InventoryHandler) RegisterRoutes(g *echo.Group) {
	// Public routes
	g.GET("/inventory", h.ListInventoryItems)
	g.GET("/inventory/:id", h.GetInventoryItem)
	g.GET("/inventory/categories", h.ListCategories)
	g.GET("/inventory/category/:category", h.ListInventoryItemsByCategory)
	g.GET("/inventory/low-stock", h.ListLowStockItems)

	// Protected routes (admin only)
	admin := g.Group("/admin/inventory")
	admin.Use(h.authMiddleware)

	admin.POST("", h.CreateInventoryItem)
	admin.PUT("/:id", h.UpdateInventoryItem)
	admin.PUT("/:id/quantity", h.UpdateInventoryQuantity)
	admin.DELETE("/:id", h.DeleteInventoryItem)
}

// CreateInventoryItem handles creating a new inventory item
func (h *InventoryHandler) CreateInventoryItem(c echo.Context) error {
	var item model.InventoryItem
	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request payload",
		})
	}

	createdItem, err := h.service.CreateInventoryItem(item)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Inventory item created successfully",
		"data":    createdItem,
	})
}

// GetInventoryItem handles getting an inventory item by ID
func (h *InventoryHandler) GetInventoryItem(c echo.Context) error {
	id := c.Param("id")

	item, err := h.service.GetInventoryItemByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Inventory item not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Inventory item retrieved successfully",
		"data":    item,
	})
}

// UpdateInventoryItem handles updating an inventory item
func (h *InventoryHandler) UpdateInventoryItem(c echo.Context) error {
	id := c.Param("id")

	var item model.InventoryItem
	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request payload",
		})
	}

	item.ID = id
	updatedItem, err := h.service.UpdateInventoryItem(id, item)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Inventory item updated successfully",
		"data":    updatedItem,
	})
}

// UpdateInventoryQuantity handles updating an inventory item's quantity
func (h *InventoryHandler) UpdateInventoryQuantity(c echo.Context) error {
	id := c.Param("id")

	var quantityData struct {
		Quantity float64 `json:"quantity"`
		Notes    string  `json:"notes"`
	}

	if err := c.Bind(&quantityData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request payload",
		})
	}

	// Get user ID from context
	userID := c.Get("user_id").(string)

	// Convert float64 to int for the service call
	quantity := int(quantityData.Quantity)

	if err := h.service.UpdateInventoryQuantity(id, quantity, userID, quantityData.Notes); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Inventory quantity updated successfully",
	})
}

// DeleteInventoryItem handles deleting an inventory item
func (h *InventoryHandler) DeleteInventoryItem(c echo.Context) error {
	id := c.Param("id")

	if err := h.service.DeleteInventoryItem(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Inventory item deleted successfully",
	})
}

// ListInventoryItems handles listing all inventory items
func (h *InventoryHandler) ListInventoryItems(c echo.Context) error {
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

	items, err := h.service.ListInventoryItems(limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to retrieve inventory items",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Inventory items retrieved successfully",
		"data":    items,
	})
}

// ListInventoryItemsByCategory handles listing inventory items by category
func (h *InventoryHandler) ListInventoryItemsByCategory(c echo.Context) error {
	category := c.Param("category")

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

	items, err := h.service.ListInventoryItemsByCategory(category, limit, offset)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Inventory items retrieved successfully",
		"data":    items,
	})
}

// ListLowStockItems handles listing inventory items with low stock
func (h *InventoryHandler) ListLowStockItems(c echo.Context) error {
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

	items, err := h.service.ListLowStockItems(limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to retrieve low stock items",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Low stock items retrieved successfully",
		"data":    items,
	})
}

// ListCategories handles listing all inventory categories
func (h *InventoryHandler) ListCategories(c echo.Context) error {
	categories, err := h.service.ListCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to retrieve categories",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Categories retrieved successfully",
		"data":    categories,
	})
}

// authMiddleware is a middleware to check if the user is authenticated and is an admin
func (h *InventoryHandler) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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
