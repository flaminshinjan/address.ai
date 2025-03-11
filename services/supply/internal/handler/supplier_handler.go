package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/flaminshinjan/address.ai/pkg/common/auth"
	"github.com/flaminshinjan/address.ai/pkg/common/response"
	"github.com/flaminshinjan/address.ai/services/supply/internal/model"
	"github.com/flaminshinjan/address.ai/services/supply/internal/service"
	"github.com/gorilla/mux"
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

// ListSuppliers handles listing all suppliers
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

	suppliers, err := h.service.ListSuppliers(limit, offset)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.Success(w, "Suppliers retrieved successfully", suppliers)
}

// ListActiveSuppliers handles listing all active suppliers
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

	suppliers, err := h.service.ListActiveSuppliers(limit, offset)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.Success(w, "Active suppliers retrieved successfully", suppliers)
}

// GetSupplier handles getting a supplier by ID
func (h *SupplierHandler) GetSupplier(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	supplier, err := h.service.GetSupplierByID(id)
	if err != nil {
		response.NotFound(w, "Supplier not found")
		return
	}

	response.Success(w, "Supplier retrieved successfully", supplier)
}

// CreateSupplier handles creating a new supplier
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

	createdSupplier, err := h.service.CreateSupplier(supplier)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, "Supplier created successfully", createdSupplier)
}

// UpdateSupplier handles updating a supplier
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

	supplier.ID = id
	updatedSupplier, err := h.service.UpdateSupplier(id, supplier)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, "Supplier updated successfully", updatedSupplier)
}

// DeleteSupplier handles deleting a supplier
func (h *SupplierHandler) DeleteSupplier(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	role := r.Context().Value("role").(string)
	if role != "admin" {
		response.Forbidden(w, "Admin access required")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.DeleteSupplier(id); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, "Supplier deleted successfully", nil)
}
