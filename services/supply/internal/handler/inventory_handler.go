package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/flaminshinjan/address.ai/pkg/common/auth"
	"github.com/flaminshinjan/address.ai/pkg/common/response"
	"github.com/flaminshinjan/address.ai/services/supply/internal/model"
	"github.com/flaminshinjan/address.ai/services/supply/internal/service"
)

// InventoryHandler handles HTTP requests for inventory items
type InventoryHandler struct {
	inventoryService *service.InventoryService
	jwtSecret        string
}

// NewInventoryHandler creates a new InventoryHandler
func NewInventoryHandler(inventoryService *service.InventoryService, jwtSecret string) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
		jwtSecret:        jwtSecret,
	}
}

// CreateInventoryItem handles creating a new inventory item
// @Summary Create a new inventory item
// @Description Create a new inventory item with the provided details
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item body model.InventoryItem true "Inventory Item"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /inventory [post]
func (h *InventoryHandler) CreateInventoryItem(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	role := r.Context().Value("role").(string)
	if role != "admin" {
		response.Forbidden(w, "Admin access required")
		return
	}

	var item model.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	createdItem, err := h.inventoryService.CreateInventoryItem(item)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, "Inventory item created successfully", createdItem)
}

// GetInventoryItem handles getting an inventory item by ID
// @Summary Get an inventory item by ID
// @Description Get an inventory item by its ID
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Inventory Item ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /inventory/{id} [get]
func (h *InventoryHandler) GetInventoryItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	item, err := h.inventoryService.GetInventoryItemByID(id)
	if err != nil {
		response.NotFound(w, "Inventory item not found")
		return
	}

	response.Success(w, "Inventory item retrieved successfully", item)
}

// UpdateInventoryItem handles updating an inventory item
// @Summary Update an inventory item
// @Description Update an inventory item with the provided details
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Inventory Item ID"
// @Param item body model.InventoryItem true "Inventory Item"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /inventory/{id} [put]
func (h *InventoryHandler) UpdateInventoryItem(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	role := r.Context().Value("role").(string)
	if role != "admin" {
		response.Forbidden(w, "Admin access required")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var item model.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	updatedItem, err := h.inventoryService.UpdateInventoryItem(id, item)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, "Inventory item updated successfully", updatedItem)
}

// UpdateInventoryQuantity handles updating an inventory item's quantity
// @Summary Update inventory quantity
// @Description Update an inventory item's quantity
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Inventory Item ID"
// @Param quantity body map[string]interface{} true "Quantity Update"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /inventory/{id}/quantity [put]
func (h *InventoryHandler) UpdateInventoryQuantity(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	role := r.Context().Value("role").(string)
	if role != "admin" {
		response.Forbidden(w, "Admin access required")
		return
	}

	// Get user ID from context
	userID := r.Context().Value("user_id").(string)

	vars := mux.Vars(r)
	id := vars["id"]

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	quantityFloat, ok := data["quantity"].(float64)
	if !ok {
		response.BadRequest(w, "Quantity is required and must be a number")
		return
	}

	quantity := int(quantityFloat)
	notes, _ := data["notes"].(string)

	if err := h.inventoryService.UpdateInventoryQuantity(id, quantity, userID, notes); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, "Inventory quantity updated successfully", nil)
}

// DeleteInventoryItem handles deleting an inventory item
// @Summary Delete an inventory item
// @Description Delete an inventory item by its ID
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Inventory Item ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /inventory/{id} [delete]
func (h *InventoryHandler) DeleteInventoryItem(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	role := r.Context().Value("role").(string)
	if role != "admin" {
		response.Forbidden(w, "Admin access required")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.inventoryService.DeleteInventoryItem(id); err != nil {
		response.NotFound(w, "Inventory item not found")
		return
	}

	response.Success(w, "Inventory item deleted successfully", nil)
}

// ListInventoryItems handles listing all inventory items
// @Summary List all inventory items
// @Description List all inventory items with pagination
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /inventory [get]
func (h *InventoryHandler) ListInventoryItems(w http.ResponseWriter, r *http.Request) {
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

	items, err := h.inventoryService.ListInventoryItems(limit, offset)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.Success(w, "Inventory items retrieved successfully", items)
}

// ListInventoryItemsByCategory handles listing inventory items by category
// @Summary List inventory items by category
// @Description List inventory items by category with pagination
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category path string true "Category"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /inventory/category/{category} [get]
func (h *InventoryHandler) ListInventoryItemsByCategory(w http.ResponseWriter, r *http.Request) {
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

	items, err := h.inventoryService.ListInventoryItemsByCategory(category, limit, offset)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.Success(w, "Inventory items retrieved successfully", items)
}

// ListLowStockItems handles listing inventory items with quantity below min_quantity
// @Summary List low stock items
// @Description List inventory items with quantity below min_quantity with pagination
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /inventory/low-stock [get]
func (h *InventoryHandler) ListLowStockItems(w http.ResponseWriter, r *http.Request) {
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

	items, err := h.inventoryService.ListLowStockItems(limit, offset)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.Success(w, "Low stock items retrieved successfully", items)
}

// ListCategories handles listing all inventory categories
// @Summary List all inventory categories
// @Description List all inventory categories
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /inventory/categories [get]
func (h *InventoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.inventoryService.ListCategories()
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.Success(w, "Categories retrieved successfully", categories)
}

// RegisterRoutes registers the routes for the inventory handler
func (h *InventoryHandler) RegisterRoutes(router *mux.Router) {
	// Protected routes
	protected := router.PathPrefix("/inventory").Subrouter()
	protected.Use(func(next http.Handler) http.Handler {
		return auth.Middleware(h.jwtSecret, next)
	})

	protected.HandleFunc("", h.ListInventoryItems).Methods("GET")
	protected.HandleFunc("/categories", h.ListCategories).Methods("GET")
	protected.HandleFunc("/category/{category}", h.ListInventoryItemsByCategory).Methods("GET")
	protected.HandleFunc("/low-stock", h.ListLowStockItems).Methods("GET")
	protected.HandleFunc("/{id}", h.GetInventoryItem).Methods("GET")

	// Admin routes
	adminRouter := protected.PathPrefix("").Subrouter()
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

	adminRouter.HandleFunc("", h.CreateInventoryItem).Methods("POST")
	adminRouter.HandleFunc("/{id}", h.UpdateInventoryItem).Methods("PUT")
	adminRouter.HandleFunc("/{id}", h.DeleteInventoryItem).Methods("DELETE")
	adminRouter.HandleFunc("/{id}/quantity", h.UpdateInventoryQuantity).Methods("PUT")
}
