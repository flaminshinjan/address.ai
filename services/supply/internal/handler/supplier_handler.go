package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yourusername/hotel-management/pkg/common/auth"
	"github.com/yourusername/hotel-management/pkg/common/response"
	"github.com/yourusername/hotel-management/services/supply/internal/model"
	"github.com/yourusername/hotel-management/services/supply/internal/service"
)

// SupplierHandler handles HTTP requests for suppliers
type SupplierHandler struct {
	supplierService *service.SupplierService
	jwtSecret       string
}

// NewSupplierHandler creates a new SupplierHandler
func NewSupplierHandler(supplierService *service.SupplierService, jwtSecret string) *SupplierHandler {
	return &SupplierHandler{
		supplierService: supplierService,
		jwtSecret:       jwtSecret,
	}
}

// CreateSupplier handles creating a new supplier
// @Summary Create a new supplier
// @Description Create a new supplier with the provided details
// @Tags suppliers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param supplier body model.Supplier true "Supplier"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /suppliers [post]
func (h *SupplierHandler) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	role := r.Context().Value("role").(string)
	if role != "admin" {
		response.Forbidden(w, "Admin access required")
		return
	}

	var supplier model.Supplier
	if err := json.NewDecoder(r.Body).Decode(&supplier); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	createdSupplier, err := h.supplierService.CreateSupplier(supplier)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, "Supplier created successfully", createdSupplier)
}

// GetSupplier handles getting a supplier by ID
// @Summary Get a supplier by ID
// @Description Get a supplier by its ID
// @Tags suppliers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Supplier ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /suppliers/{id} [get]
func (h *SupplierHandler) GetSupplier(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	supplier, err := h.supplierService.GetSupplierByID(id)
	if err != nil {
		response.NotFound(w, "Supplier not found")
		return
	}

	response.Success(w, "Supplier retrieved successfully", supplier)
}

// UpdateSupplier handles updating a supplier
// @Summary Update a supplier
// @Description Update a supplier with the provided details
// @Tags suppliers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Supplier ID"
// @Param supplier body model.Supplier true "Supplier"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /suppliers/{id} [put]
func (h *SupplierHandler) UpdateSupplier(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	role := r.Context().Value("role").(string)
	if role != "admin" {
		response.Forbidden(w, "Admin access required")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var supplier model.Supplier
	if err := json.NewDecoder(r.Body).Decode(&supplier); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	updatedSupplier, err := h.supplierService.UpdateSupplier(id, supplier)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, "Supplier updated successfully", updatedSupplier)
}

// DeleteSupplier handles deleting a supplier
// @Summary Delete a supplier
// @Description Delete a supplier by its ID
// @Tags suppliers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Supplier ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /suppliers/{id} [delete]
func (h *SupplierHandler) DeleteSupplier(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	role := r.Context().Value("role").(string)
	if role != "admin" {
		response.Forbidden(w, "Admin access required")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.supplierService.DeleteSupplier(id); err != nil {
		response.NotFound(w, "Supplier not found")
		return
	}

	response.Success(w, "Supplier deleted successfully", nil)
}

// ListSuppliers handles listing all suppliers
// @Summary List all suppliers
// @Description List all suppliers with pagination
// @Tags suppliers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /suppliers [get]
func (h *SupplierHandler) ListSuppliers(w http.ResponseWriter, r *http.Request) {
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

	suppliers, err := h.supplierService.ListSuppliers(limit, offset)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.Success(w, "Suppliers retrieved successfully", suppliers)
}

// ListActiveSuppliers handles listing all active suppliers
// @Summary List all active suppliers
// @Description List all active suppliers with pagination
// @Tags suppliers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /suppliers/active [get]
func (h *SupplierHandler) ListActiveSuppliers(w http.ResponseWriter, r *http.Request) {
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

	suppliers, err := h.supplierService.ListActiveSuppliers(limit, offset)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.Success(w, "Active suppliers retrieved successfully", suppliers)
}

// RegisterRoutes registers the routes for the supplier handler
func (h *SupplierHandler) RegisterRoutes(router *mux.Router) {
	// Protected routes
	protected := router.PathPrefix("/suppliers").Subrouter()
	protected.Use(func(next http.Handler) http.Handler {
		return auth.Middleware(h.jwtSecret, next)
	})

	protected.HandleFunc("", h.ListSuppliers).Methods("GET")
	protected.HandleFunc("/active", h.ListActiveSuppliers).Methods("GET")
	protected.HandleFunc("/{id}", h.GetSupplier).Methods("GET")

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

	adminRouter.HandleFunc("", h.CreateSupplier).Methods("POST")
	adminRouter.HandleFunc("/{id}", h.UpdateSupplier).Methods("PUT")
	adminRouter.HandleFunc("/{id}", h.DeleteSupplier).Methods("DELETE")
}
