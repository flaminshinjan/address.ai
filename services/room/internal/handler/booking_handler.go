package handler

import (
	"net/http"
	"strconv"

	"github.com/flaminshinjan/address.ai/pkg/common/auth"
	"github.com/flaminshinjan/address.ai/services/room/internal/model"
	"github.com/flaminshinjan/address.ai/services/room/internal/service"
	"github.com/labstack/echo/v4"
)

// BookingHandler handles HTTP requests for bookings
type BookingHandler struct {
	service   *service.BookingService
	jwtSecret string
}

// NewBookingHandler creates a new BookingHandler
func NewBookingHandler(service *service.BookingService, jwtSecret string) *BookingHandler {
	return &BookingHandler{
		service:   service,
		jwtSecret: jwtSecret,
	}
}

// CreateBooking handles creating a new booking
func (h *BookingHandler) CreateBooking(c echo.Context) error {
	// Get user ID from context
	userID := c.Get("user_id").(string)

	var req model.BookingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request payload",
		})
	}

	booking, err := h.service.CreateBooking(userID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Booking created successfully",
		"data":    booking,
	})
}

// GetBooking handles getting a booking by ID
func (h *BookingHandler) GetBooking(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("user_id").(string)
	role := c.Get("role").(string)

	booking, err := h.service.GetBookingByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Booking not found",
		})
	}

	// Check if user is authorized to view this booking
	if booking.UserID != userID && role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"success": false,
			"error":   "You are not authorized to view this booking",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Booking retrieved successfully",
		"data":    booking,
	})
}

// CancelBooking handles canceling a booking
func (h *BookingHandler) CancelBooking(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("user_id").(string)
	role := c.Get("role").(string)

	// Check if booking exists and user is authorized
	existingBooking, err := h.service.GetBookingByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Booking not found",
		})
	}

	// Check if user is authorized to cancel this booking
	if existingBooking.UserID != userID && role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"success": false,
			"error":   "You are not authorized to cancel this booking",
		})
	}

	if err := h.service.CancelBooking(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Booking canceled successfully",
	})
}

// ListUserBookings handles listing all bookings for a user
func (h *BookingHandler) ListUserBookings(c echo.Context) error {
	userID := c.Get("user_id").(string)

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

	bookings, err := h.service.GetBookingsByUserID(userID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to retrieve bookings",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Bookings retrieved successfully",
		"data":    bookings,
	})
}

// ListAllBookings handles listing all bookings (admin only)
func (h *BookingHandler) ListAllBookings(c echo.Context) error {
	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"success": false,
			"error":   "Admin access required",
		})
	}

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

	bookings, err := h.service.ListBookings(limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to retrieve bookings",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Bookings retrieved successfully",
		"data":    bookings,
	})
}

// RegisterRoutes registers the routes for the booking handler
func (h *BookingHandler) RegisterRoutes(g *echo.Group) {
	// Protected routes
	bookings := g.Group("/bookings")
	bookings.Use(h.authMiddleware)

	bookings.POST("", h.CreateBooking)
	bookings.GET("/my", h.ListUserBookings)
	bookings.GET("/:id", h.GetBooking)
	bookings.DELETE("/:id", h.CancelBooking)

	// Admin routes
	admin := g.Group("/admin/bookings")
	admin.Use(h.authMiddleware, h.adminMiddleware)
	admin.GET("", h.ListAllBookings)
}

// authMiddleware is a middleware to check if the user is authenticated
func (h *BookingHandler) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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

// adminMiddleware is a middleware to check if the user is an admin
func (h *BookingHandler) adminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := c.Get("role").(string)
		if role != "admin" {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"success": false,
				"error":   "Admin access required",
			})
		}
		return next(c)
	}
}
