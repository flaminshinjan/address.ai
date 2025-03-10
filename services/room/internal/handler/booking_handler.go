package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/flaminshinjan/address.ai/pkg/common/auth"
	"github.com/flaminshinjan/address.ai/pkg/common/response"
	"github.com/flaminshinjan/address.ai/services/room/internal/model"
	"github.com/flaminshinjan/address.ai/services/room/internal/service"
)

// BookingHandler handles HTTP requests for bookings
type BookingHandler struct {
	bookingService *service.BookingService
	jwtSecret      string
}

// NewBookingHandler creates a new BookingHandler
func NewBookingHandler(bookingService *service.BookingService, jwtSecret string) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
		jwtSecret:      jwtSecret,
	}
}

// CreateBooking handles creating a new booking
// @Summary Create a new booking
// @Description Create a new booking with the provided details
// @Tags bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param booking body model.BookingRequest true "Booking Request"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /bookings [post]
func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID := r.Context().Value("user_id").(string)

	var req model.BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	booking, err := h.bookingService.CreateBooking(userID, req)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, "Booking created successfully", booking)
}

// GetBooking handles getting a booking by ID
// @Summary Get a booking by ID
// @Description Get a booking by its ID
// @Tags bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Booking ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /bookings/{id} [get]
func (h *BookingHandler) GetBooking(w http.ResponseWriter, r *http.Request) {
	// Get user ID and role from context
	userID := r.Context().Value("user_id").(string)
	role := r.Context().Value("role").(string)

	vars := mux.Vars(r)
	id := vars["id"]

	booking, err := h.bookingService.GetBookingByID(id)
	if err != nil {
		response.NotFound(w, "Booking not found")
		return
	}

	// Check if user is authorized to view this booking
	if booking.UserID != userID && role != "admin" {
		response.Forbidden(w, "You are not authorized to view this booking")
		return
	}

	response.Success(w, "Booking retrieved successfully", booking)
}

// GetUserBookings handles getting bookings for the authenticated user
// @Summary Get user bookings
// @Description Get bookings for the authenticated user
// @Tags bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /bookings/user [get]
func (h *BookingHandler) GetUserBookings(w http.ResponseWriter, r *http.Request) {
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

	bookings, err := h.bookingService.GetBookingsByUserID(userID, limit, offset)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.Success(w, "Bookings retrieved successfully", bookings)
}

// CancelBooking handles cancelling a booking
// @Summary Cancel a booking
// @Description Cancel a booking by its ID
// @Tags bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Booking ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /bookings/{id}/cancel [put]
func (h *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	// Get user ID and role from context
	userID := r.Context().Value("user_id").(string)
	role := r.Context().Value("role").(string)

	vars := mux.Vars(r)
	id := vars["id"]

	// Check if booking exists and belongs to the user
	booking, err := h.bookingService.GetBookingByID(id)
	if err != nil {
		response.NotFound(w, "Booking not found")
		return
	}

	// Check if user is authorized to cancel this booking
	if booking.UserID != userID && role != "admin" {
		response.Forbidden(w, "You are not authorized to cancel this booking")
		return
	}

	if err := h.bookingService.CancelBooking(id); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, "Booking cancelled successfully", nil)
}

// ListBookings handles listing all bookings (admin only)
// @Summary List all bookings
// @Description List all bookings with pagination (admin only)
// @Tags bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /bookings [get]
func (h *BookingHandler) ListBookings(w http.ResponseWriter, r *http.Request) {
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

	bookings, err := h.bookingService.ListBookings(limit, offset)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.Success(w, "Bookings retrieved successfully", bookings)
}

// RegisterRoutes registers the routes for the booking handler
func (h *BookingHandler) RegisterRoutes(router *mux.Router) {
	// Protected routes
	protected := router.PathPrefix("/bookings").Subrouter()
	protected.Use(func(next http.Handler) http.Handler {
		return auth.Middleware(h.jwtSecret, next)
	})

	protected.HandleFunc("", h.CreateBooking).Methods("POST")
	protected.HandleFunc("/user", h.GetUserBookings).Methods("GET")
	protected.HandleFunc("/{id}", h.GetBooking).Methods("GET")
	protected.HandleFunc("/{id}/cancel", h.CancelBooking).Methods("PUT")

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

	adminRouter.HandleFunc("", h.ListBookings).Methods("GET")
}
