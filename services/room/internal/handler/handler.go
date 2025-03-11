package handler

import (
	"github.com/flaminshinjan/address.ai/services/room/internal/service"
	"github.com/labstack/echo/v4"
)

// Handler is the main handler for the room service
type Handler struct {
	RoomHandler    *RoomHandler
	BookingHandler *BookingHandler
}

// NewHandler creates a new Handler
func NewHandler(
	roomService *service.RoomService,
	bookingService *service.BookingService,
	jwtSecret string,
) *Handler {
	return &Handler{
		RoomHandler:    NewRoomHandler(roomService, bookingService, jwtSecret),
		BookingHandler: NewBookingHandler(bookingService, jwtSecret),
	}
}

// RegisterRoutes registers all routes for the room service
func (h *Handler) RegisterRoutes(g *echo.Group) {
	// Register room routes
	h.RoomHandler.RegisterRoutes(g)

	// Register booking routes
	h.BookingHandler.RegisterRoutes(g)
}
