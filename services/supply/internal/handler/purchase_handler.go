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

// PurchaseHandler handles HTTP requests for purchase orders
type PurchaseHandler struct {
	purchaseService *service.PurchaseService
	jwtSecret       string
}

// NewPurchaseHandler creates a new PurchaseHandler
func NewPurchaseHandler(purchaseService *service.PurchaseService, jwtSecret string) *PurchaseHandler {
	return &PurchaseHandler{
		purchaseService: purchaseService,
		jwtSecret:       jwtSecret,
	}
}

// CreatePurchaseOrder handles creating a new purchase order
// @Summary Create a new purchase order
// @Description Create a new purchase order with the provided details
// @Tags purchase-orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order body model.CreatePurchaseOrderRequest true "Purchase Order Request"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /purchase-orders [post]
func (h *PurchaseHandler) CreatePurchaseOrder(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	role := r.Context().Value("role").(string)
	if role != "admin" {
		response.Forbidden(w, "Admin access required")
		return
	}

	// Get user ID from context
	userID := r.Context().Value("user_id").(string)

	var req model.CreatePurchaseOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	order, err := h.purchaseService.CreatePurchaseOrder(userID, req)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, "Purchase order created successfully", order)
}

// GetPurchaseOrder handles getting a purchase order by ID
// @Summary Get a purchase order by ID
// @Description Get a purchase order by its ID
// @Tags purchase-orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Purchase Order ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /purchase-orders/{id} [get]
func (h *PurchaseHandler) GetPurchaseOrder(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	role := r.Context().Value("role").(string)
	if role != "admin" {
		response.Forbidden(w, "Admin access required")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	order, err := h.purchaseService.GetPurchaseOrderByID(id)
	if err != nil {
		response.NotFound(w, "Purchase order not found")
		return
	}

	response.Success(w, "Purchase order retrieved successfully", order)
}

// UpdatePurchaseOrderStatus handles updating a purchase order's status
// @Summary Update purchase order status
// @Description Update a purchase order's status
// @Tags purchase-orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Purchase Order ID"
// @Param status body map[string]string true "Status"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /purchase-orders/{id}/status [put]
func (h *PurchaseHandler) UpdatePurchaseOrderStatus(w http.ResponseWriter, r *http.Request) {
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

	var statusData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&statusData); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	status, ok := statusData["status"]
	if !ok {
		response.BadRequest(w, "Status is required")
		return
	}

	if err := h.purchaseService.UpdatePurchaseOrderStatus(id, status, userID); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, "Purchase order status updated successfully", nil)
}

// ListPurchaseOrders handles listing all purchase orders
// @Summary List all purchase orders
// @Description List all purchase orders with pagination
// @Tags purchase-orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /purchase-orders [get]
func (h *PurchaseHandler) ListPurchaseOrders(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	role := r.Context().Value("role").(string)
	if role != "admin" {
		response.Forbidden(w, "Admin access required")
		return
	}

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

	orders, err := h.purchaseService.ListPurchaseOrders(limit, offset)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.Success(w, "Purchase orders retrieved successfully", orders)
}

// ListPurchaseOrdersByStatus handles listing purchase orders by status
// @Summary List purchase orders by status
// @Description List purchase orders by status with pagination
// @Tags purchase-orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status path string true "Status"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /purchase-orders/status/{status} [get]
func (h *PurchaseHandler) ListPurchaseOrdersByStatus(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	role := r.Context().Value("role").(string)
	if role != "admin" {
		response.Forbidden(w, "Admin access required")
		return
	}

	vars := mux.Vars(r)
	status := vars["status"]

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

	orders, err := h.purchaseService.ListPurchaseOrdersByStatus(status, limit, offset)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, "Purchase orders retrieved successfully", orders)
}

// RegisterRoutes registers the routes for the purchase handler
func (h *PurchaseHandler) RegisterRoutes(router *mux.Router) {
	// Protected routes
	protected := router.PathPrefix("/purchase-orders").Subrouter()
	protected.Use(func(next http.Handler) http.Handler {
		return auth.Middleware(h.jwtSecret, next)
	})

	// Admin routes (all purchase order routes require admin access)
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

	adminRouter.HandleFunc("", h.CreatePurchaseOrder).Methods("POST")
	adminRouter.HandleFunc("", h.ListPurchaseOrders).Methods("GET")
	adminRouter.HandleFunc("/status/{status}", h.ListPurchaseOrdersByStatus).Methods("GET")
	adminRouter.HandleFunc("/{id}", h.GetPurchaseOrder).Methods("GET")
	adminRouter.HandleFunc("/{id}/status", h.UpdatePurchaseOrderStatus).Methods("PUT")
}
