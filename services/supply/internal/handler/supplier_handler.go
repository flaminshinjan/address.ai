package handler

import (
	"net/http"
	"strconv"

	"github.com/flaminshinjan/address.ai/pkg/common/auth"
	"github.com/flaminshinjan/address.ai/services/supply/internal/model"
	"github.com/flaminshinjan/address.ai/services/supply/internal/service"
	"github.com/labstack/echo/v4"
)

// SupplierHandler handles HTTP requests for suppliers
type SupplierHandler struct {
	service   *service.SupplierService
	jwtSecret string
}

// NewSupplierHandler creates a new SupplierHandler
func NewSupplierHandler(service *service.SupplierService, jwtSecret string) *SupplierHandler {
	return &SupplierHandler{
		service:   service,
		jwtSecret: jwtSecret,
	}
}

// RegisterRoutes registers the routes for the supplier handler
func (h *SupplierHandler) RegisterRoutes(g *echo.Group) {
	// Public routes
	g.GET("/suppliers", h.ListSuppliers)
	g.GET("/suppliers/:id", h.GetSupplier)

	// Protected routes (admin only)
	admin := g.Group("/admin/suppliers")
	admin.Use(h.authMiddleware)

	admin.POST("", h.CreateSupplier)
	admin.PUT("/:id", h.UpdateSupplier)
	admin.DELETE("/:id", h.DeleteSupplier)
}

// ListSuppliers handles listing all suppliers
func (h *SupplierHandler) ListSuppliers(c echo.Context) error {
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

	suppliers, err := h.service.ListSuppliers(limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to retrieve suppliers",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Suppliers retrieved successfully",
		"data":    suppliers,
	})
}

// GetSupplier handles getting a supplier by ID
func (h *SupplierHandler) GetSupplier(c echo.Context) error {
	id := c.Param("id")

	supplier, err := h.service.GetSupplierByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Supplier not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Supplier retrieved successfully",
		"data":    supplier,
	})
}

// CreateSupplier handles creating a new supplier
func (h *SupplierHandler) CreateSupplier(c echo.Context) error {
	var supplier model.Supplier
	if err := c.Bind(&supplier); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request payload",
		})
	}

	createdSupplier, err := h.service.CreateSupplier(supplier)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Supplier created successfully",
		"data":    createdSupplier,
	})
}

// UpdateSupplier handles updating a supplier
func (h *SupplierHandler) UpdateSupplier(c echo.Context) error {
	id := c.Param("id")

	var supplier model.Supplier
	if err := c.Bind(&supplier); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request payload",
		})
	}

	supplier.ID = id
	updatedSupplier, err := h.service.UpdateSupplier(id, supplier)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Supplier updated successfully",
		"data":    updatedSupplier,
	})
}

// DeleteSupplier handles deleting a supplier
func (h *SupplierHandler) DeleteSupplier(c echo.Context) error {
	id := c.Param("id")

	if err := h.service.DeleteSupplier(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Supplier deleted successfully",
	})
}

// authMiddleware is a middleware to check if the user is authenticated and is an admin
func (h *SupplierHandler) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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
