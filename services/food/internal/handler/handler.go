package handler

import (
	"github.com/flaminshinjan/address.ai/services/food/internal/service"
	"github.com/labstack/echo/v4"
)

// Handler is the main handler for the food service
type Handler struct {
	MenuHandler  *MenuHandler
	OrderHandler *OrderHandler
}

// NewHandler creates a new Handler
func NewHandler(
	menuService *service.MenuService,
	orderService *service.OrderService,
	jwtSecret string,
) *Handler {
	return &Handler{
		MenuHandler:  NewMenuHandler(menuService, jwtSecret),
		OrderHandler: NewOrderHandler(orderService, jwtSecret),
	}
}

// RegisterRoutes registers all routes for the food service
func (h *Handler) RegisterRoutes(g *echo.Group) {
	// Register menu routes
	h.MenuHandler.RegisterRoutes(g)

	// Register order routes
	h.OrderHandler.RegisterRoutes(g)
}
