package handler

import (
	"github.com/flaminshinjan/address.ai/services/supply/internal/service"
	"github.com/gorilla/mux"
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
func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Register supplier routes
	h.SupplierHandler.RegisterRoutes(router)

	// Register inventory routes
	h.InventoryHandler.RegisterRoutes(router)

	// Register purchase routes
	h.PurchaseHandler.RegisterRoutes(router)
}
