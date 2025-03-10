package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/flaminshinjan/address.ai/pkg/common/auth"
	"github.com/flaminshinjan/address.ai/pkg/common/response"
	"github.com/flaminshinjan/address.ai/services/food/internal/model"
	"github.com/flaminshinjan/address.ai/services/food/internal/service"
)

// OrderHandler handles HTTP requests for food orders
type OrderHandler struct {
	orderService *service.OrderService
	jwtSecret    string
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(orderService *service.OrderService, jwtSecret string) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		jwtSecret:    jwtSecret,
	}
}

// CreateOrder handles creating a new food order
// @Summary Create a new food order
// @Description Create a new food order with the provided details
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order body model.CreateOrderRequest true "Order Request"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID := r.Context().Value("user_id").(string)

	var req model.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	order, err := h.orderService.CreateOrder(userID, req)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, "Order created successfully", order)
}

// GetOrder handles getting a food order by ID
// @Summary Get a food order by ID
// @Description Get a food order by its ID
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	// Get user ID and role from context
	userID := r.Context().Value("user_id").(string)
	role := r.Context().Value("role").(string)

	vars := mux.Vars(r)
	id := vars["id"]

	order, err := h.orderService.GetOrderByID(id)
	if err != nil {
		response.NotFound(w, "Order not found")
		return
	}

	// Check if user is authorized to view this order
	if order.UserID != userID && role != "admin" {
		response.Forbidden(w, "You are not authorized to view this order")
		return
	}

	response.Success(w, "Order retrieved successfully", order)
}

// UpdateOrderStatus handles updating a food order's status
// @Summary Update order status
// @Description Update a food order's status
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Param status body map[string]string true "Status"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders/{id}/status [put]
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	role := r.Context().Value("role").(string)
	if role != "admin" {
		response.Forbidden(w, "Admin access required")
		return
	}

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

	if err := h.orderService.UpdateOrderStatus(id, status); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, "Order status updated successfully", nil)
}

// GetUserOrders handles getting food orders for the authenticated user
// @Summary Get user orders
// @Description Get food orders for the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders/user [get]
func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID := r.Context().Value("user_id").(string)

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

	orders, err := h.orderService.ListOrdersByUserID(userID, limit, offset)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.Success(w, "Orders retrieved successfully", orders)
}

// ListOrders handles listing all food orders (admin only)
// @Summary List all food orders
// @Description List all food orders with pagination (admin only)
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders [get]
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
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

	orders, err := h.orderService.ListOrders(limit, offset)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.Success(w, "Orders retrieved successfully", orders)
}

// ListOrdersByStatus handles listing food orders by status (admin only)
// @Summary List orders by status
// @Description List food orders by status with pagination (admin only)
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status path string true "Status"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders/status/{status} [get]
func (h *OrderHandler) ListOrdersByStatus(w http.ResponseWriter, r *http.Request) {
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

	orders, err := h.orderService.ListOrdersByStatus(status, limit, offset)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, "Orders retrieved successfully", orders)
}

// RegisterRoutes registers the routes for the order handler
func (h *OrderHandler) RegisterRoutes(router *mux.Router) {
	// Protected routes
	protected := router.PathPrefix("/orders").Subrouter()
	protected.Use(func(next http.Handler) http.Handler {
		return auth.Middleware(h.jwtSecret, next)
	})

	protected.HandleFunc("", h.CreateOrder).Methods("POST")
	protected.HandleFunc("/user", h.GetUserOrders).Methods("GET")
	protected.HandleFunc("/{id}", h.GetOrder).Methods("GET")

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

	adminRouter.HandleFunc("", h.ListOrders).Methods("GET")
	adminRouter.HandleFunc("/status/{status}", h.ListOrdersByStatus).Methods("GET")
	adminRouter.HandleFunc("/{id}/status", h.UpdateOrderStatus).Methods("PUT")
}
