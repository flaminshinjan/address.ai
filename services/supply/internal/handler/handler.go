package handler

import (
	"github.com/flaminshinjan/address.ai/services/supply/internal/service"
	"github.com/labstack/echo/v4"
)

// Handler is the main handler for the supply microservice
type Handler struct {
	SupplierHandler  *SupplierHandler
	InventoryHandler *InventoryHandler
	PurchaseHandler  *PurchaseHandler
}

// NewHandler creates a new Handler
func NewHandler(
	supplierService *service.SupplierService,
	inventoryService *service.InventoryService,
	purchaseService *service.PurchaseService,
	jwtSecret string,
) *Handler {
	return &Handler{
		SupplierHandler:  NewSupplierHandler(supplierService, jwtSecret),
		InventoryHandler: NewInventoryHandler(inventoryService, jwtSecret),
		PurchaseHandler:  NewPurchaseHandler(purchaseService, jwtSecret),
	}
}

// RegisterRoutes registers all routes for the supply microservice
func (h *Handler) RegisterRoutes(g *echo.Group) {
	// Register supplier routes
	h.SupplierHandler.RegisterRoutes(g)

	// Register inventory routes
	h.InventoryHandler.RegisterRoutes(g)

	// Register purchase routes
	h.PurchaseHandler.RegisterRoutes(g)
}
