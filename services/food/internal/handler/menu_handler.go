package handler

import (
	"net/http"
	"strconv"

	"github.com/flaminshinjan/address.ai/pkg/common/auth"
	"github.com/flaminshinjan/address.ai/services/food/internal/model"
	"github.com/flaminshinjan/address.ai/services/food/internal/service"
	"github.com/labstack/echo/v4"
)

// MenuHandler handles HTTP requests for menu items
type MenuHandler struct {
	service   *service.MenuService
	jwtSecret string
}

// NewMenuHandler creates a new MenuHandler
func NewMenuHandler(service *service.MenuService, jwtSecret string) *MenuHandler {
	return &MenuHandler{
		service:   service,
		jwtSecret: jwtSecret,
	}
}

// CreateMenuItem handles creating a new menu item
// @Summary Create a new menu item
// @Description Create a new menu item with the provided details
// @Tags menu
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item body model.MenuItem true "Menu Item"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /menu [post]
func (h *MenuHandler) CreateMenuItem(c echo.Context) error {
	// Check if user is admin
	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"success": false,
			"error":   "Admin access required",
		})
	}

	var item model.MenuItem
	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request payload",
		})
	}

	createdItem, err := h.service.CreateMenuItem(item)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Menu item created successfully",
		"data":    createdItem,
	})
}

// GetMenuItem handles getting a menu item by ID
// @Summary Get a menu item by ID
// @Description Get a menu item by its ID
// @Tags menu
// @Accept json
// @Produce json
// @Param id path string true "Menu Item ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /menu/{id} [get]
func (h *MenuHandler) GetMenuItem(c echo.Context) error {
	id := c.Param("id")

	item, err := h.service.GetMenuItemByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Menu item not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Menu item retrieved successfully",
		"data":    item,
	})
}

// UpdateMenuItem handles updating a menu item
// @Summary Update a menu item
// @Description Update a menu item with the provided details
// @Tags menu
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Menu Item ID"
// @Param item body model.MenuItem true "Menu Item"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /menu/{id} [put]
func (h *MenuHandler) UpdateMenuItem(c echo.Context) error {
	// Check if user is admin
	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"success": false,
			"error":   "Admin access required",
		})
	}

	id := c.Param("id")

	var item model.MenuItem
	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request payload",
		})
	}

	item.ID = id
	updatedItem, err := h.service.UpdateMenuItem(id, item)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Menu item updated successfully",
		"data":    updatedItem,
	})
}

// DeleteMenuItem handles deleting a menu item
// @Summary Delete a menu item
// @Description Delete a menu item by its ID
// @Tags menu
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Menu Item ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /menu/{id} [delete]
func (h *MenuHandler) DeleteMenuItem(c echo.Context) error {
	// Check if user is admin
	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"success": false,
			"error":   "Admin access required",
		})
	}

	id := c.Param("id")

	if err := h.service.DeleteMenuItem(id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Menu item not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Menu item deleted successfully",
	})
}

// ListMenuItems handles listing all menu items
// @Summary List all menu items
// @Description List all menu items with pagination
// @Tags menu
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /menu [get]
func (h *MenuHandler) ListMenuItems(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")
	category := c.QueryParam("category")

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

	var items interface{}
	var err error

	if category != "" {
		items, err = h.service.ListMenuItemsByCategory(category, limit, offset)
	} else {
		items, err = h.service.ListMenuItems(limit, offset)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to retrieve menu items",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Menu items retrieved successfully",
		"data":    items,
	})
}

// ListCategories handles listing all menu categories
// @Summary List all menu categories
// @Description List all menu categories
// @Tags menu
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /menu/categories [get]
func (h *MenuHandler) ListCategories(c echo.Context) error {
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

// RegisterRoutes registers the routes for the menu handler
func (h *MenuHandler) RegisterRoutes(g *echo.Group) {
	// Menu routes
	menu := g.Group("/menu")

	// Public routes
	menu.GET("", h.ListMenuItems)
	menu.GET("/categories", h.ListCategories)
	menu.GET("/:id", h.GetMenuItem)

	// Protected routes
	admin := menu.Group("")
	admin.Use(h.authMiddleware)

	admin.POST("", h.CreateMenuItem)
	admin.PUT("/:id", h.UpdateMenuItem)
	admin.DELETE("/:id", h.DeleteMenuItem)
}

// authMiddleware is a middleware to check if the user is authenticated
func (h *MenuHandler) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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
