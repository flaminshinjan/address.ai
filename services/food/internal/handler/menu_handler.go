package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yourusername/hotel-management/pkg/common/auth"
	"github.com/yourusername/hotel-management/pkg/common/response"
	"github.com/yourusername/hotel-management/services/food/internal/model"
	"github.com/yourusername/hotel-management/services/food/internal/service"
)

// MenuHandler handles HTTP requests for menu items
type MenuHandler struct {
	menuService *service.MenuService
	jwtSecret   string
}

// NewMenuHandler creates a new MenuHandler
func NewMenuHandler(menuService *service.MenuService, jwtSecret string) *MenuHandler {
	return &MenuHandler{
		menuService: menuService,
		jwtSecret:   jwtSecret,
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
func (h *MenuHandler) CreateMenuItem(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	role := r.Context().Value("role").(string)
	if role != "admin" {
		response.Forbidden(w, "Admin access required")
		return
	}

	var item model.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	createdItem, err := h.menuService.CreateMenuItem(item)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, "Menu item created successfully", createdItem)
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
func (h *MenuHandler) GetMenuItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	item, err := h.menuService.GetMenuItemByID(id)
	if err != nil {
		response.NotFound(w, "Menu item not found")
		return
	}

	response.Success(w, "Menu item retrieved successfully", item)
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
func (h *MenuHandler) UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	role := r.Context().Value("role").(string)
	if role != "admin" {
		response.Forbidden(w, "Admin access required")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var item model.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	updatedItem, err := h.menuService.UpdateMenuItem(id, item)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, "Menu item updated successfully", updatedItem)
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
func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	role := r.Context().Value("role").(string)
	if role != "admin" {
		response.Forbidden(w, "Admin access required")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.menuService.DeleteMenuItem(id); err != nil {
		response.NotFound(w, "Menu item not found")
		return
	}

	response.Success(w, "Menu item deleted successfully", nil)
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
func (h *MenuHandler) ListMenuItems(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

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

	items, err := h.menuService.ListMenuItems(limit, offset)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.Success(w, "Menu items retrieved successfully", items)
}

// ListMenuItemsByCategory handles listing menu items by category
// @Summary List menu items by category
// @Description List menu items by category with pagination
// @Tags menu
// @Accept json
// @Produce json
// @Param category path string true "Category"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /menu/category/{category} [get]
func (h *MenuHandler) ListMenuItemsByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]

	// Get query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

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

	items, err := h.menuService.ListMenuItemsByCategory(category, limit, offset)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.Success(w, "Menu items retrieved successfully", items)
}

// ListAvailableMenuItems handles listing all available menu items
// @Summary List all available menu items
// @Description List all available menu items with pagination
// @Tags menu
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /menu/available [get]
func (h *MenuHandler) ListAvailableMenuItems(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

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

	items, err := h.menuService.ListAvailableMenuItems(limit, offset)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.Success(w, "Available menu items retrieved successfully", items)
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
func (h *MenuHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.menuService.ListCategories()
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.Success(w, "Categories retrieved successfully", categories)
}

// RegisterRoutes registers the routes for the menu handler
func (h *MenuHandler) RegisterRoutes(router *mux.Router) {
	// Public routes
	router.HandleFunc("/menu", h.ListMenuItems).Methods("GET")
	router.HandleFunc("/menu/available", h.ListAvailableMenuItems).Methods("GET")
	router.HandleFunc("/menu/categories", h.ListCategories).Methods("GET")
	router.HandleFunc("/menu/category/{category}", h.ListMenuItemsByCategory).Methods("GET")
	router.HandleFunc("/menu/{id}", h.GetMenuItem).Methods("GET")

	// Protected routes
	protected := router.PathPrefix("").Subrouter()
	protected.Use(func(next http.Handler) http.Handler {
		return auth.Middleware(h.jwtSecret, next)
	})

	// Admin routes
	adminRouter := protected.PathPrefix("/menu").Subrouter()
	adminRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := r.Context().Value("role").(string)
			if role != "admin" {
				response.Forbidden(w, "Admin access required")
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	adminRouter.HandleFunc("", h.CreateMenuItem).Methods("POST")
	adminRouter.HandleFunc("/{id}", h.UpdateMenuItem).Methods("PUT")
	adminRouter.HandleFunc("/{id}", h.DeleteMenuItem).Methods("DELETE")
}
