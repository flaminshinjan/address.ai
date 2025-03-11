package handler

import (
	"net/http"
	"strconv"

	"github.com/flaminshinjan/address.ai/pkg/common/auth"
	"github.com/flaminshinjan/address.ai/services/room/internal/model"
	"github.com/flaminshinjan/address.ai/services/room/internal/service"
	"github.com/labstack/echo/v4"
)

// RoomHandler handles HTTP requests for rooms
type RoomHandler struct {
	roomService    *service.RoomService
	bookingService *service.BookingService
	jwtSecret      string
}

// NewRoomHandler creates a new RoomHandler
func NewRoomHandler(roomService *service.RoomService, bookingService *service.BookingService, jwtSecret string) *RoomHandler {
	return &RoomHandler{
		roomService:    roomService,
		bookingService: bookingService,
		jwtSecret:      jwtSecret,
	}
}

// ListRooms handles listing all rooms
func (h *RoomHandler) ListRooms(c echo.Context) error {
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

	rooms, err := h.roomService.ListRooms(limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to retrieve rooms",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Rooms retrieved successfully",
		"data":    rooms,
	})
}

// GetRoom handles getting a room by ID
func (h *RoomHandler) GetRoom(c echo.Context) error {
	id := c.Param("id")

	room, err := h.roomService.GetRoomByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Room not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Room retrieved successfully",
		"data":    room,
	})
}

// CreateRoom handles creating a new room
func (h *RoomHandler) CreateRoom(c echo.Context) error {
	// Check if user is admin
	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"success": false,
			"error":   "Admin access required",
		})
	}

	var room model.Room
	if err := c.Bind(&room); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request payload",
		})
	}

	createdRoom, err := h.roomService.CreateRoom(room)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Room created successfully",
		"data":    createdRoom,
	})
}

// UpdateRoom handles updating a room
func (h *RoomHandler) UpdateRoom(c echo.Context) error {
	// Check if user is admin
	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"success": false,
			"error":   "Admin access required",
		})
	}

	id := c.Param("id")

	var room model.Room
	if err := c.Bind(&room); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request payload",
		})
	}

	room.ID = id
	updatedRoom, err := h.roomService.UpdateRoom(id, room)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Room updated successfully",
		"data":    updatedRoom,
	})
}

// DeleteRoom handles deleting a room
func (h *RoomHandler) DeleteRoom(c echo.Context) error {
	// Check if user is admin
	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"success": false,
			"error":   "Admin access required",
		})
	}

	id := c.Param("id")

	if err := h.roomService.DeleteRoom(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Room deleted successfully",
	})
}

// UpdateRoomStatus handles updating a room's status
func (h *RoomHandler) UpdateRoomStatus(c echo.Context) error {
	// Check if user is admin
	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"success": false,
			"error":   "Admin access required",
		})
	}

	id := c.Param("id")

	var statusData map[string]string
	if err := c.Bind(&statusData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request payload",
		})
	}

	status, ok := statusData["status"]
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Status is required",
		})
	}

	if err := h.roomService.UpdateRoomStatus(id, status); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Room status updated successfully",
	})
}

// RegisterRoutes registers the routes for the room handler
func (h *RoomHandler) RegisterRoutes(g *echo.Group) {
	// Public routes
	g.GET("/rooms", h.ListRooms)
	g.GET("/rooms/:id", h.GetRoom)

	// Protected routes
	admin := g.Group("/admin/rooms")
	admin.Use(h.authMiddleware)

	admin.POST("", h.CreateRoom)
	admin.PUT("/:id", h.UpdateRoom)
	admin.DELETE("/:id", h.DeleteRoom)
	admin.PUT("/:id/status", h.UpdateRoomStatus)
}

// authMiddleware is a middleware to check if the user is authenticated
func (h *RoomHandler) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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
