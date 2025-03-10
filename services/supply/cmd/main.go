package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/flaminshinjan/address.ai/pkg/common/db"
	"github.com/flaminshinjan/address.ai/services/supply/internal/handler"
	"github.com/flaminshinjan/address.ai/services/supply/internal/repository"
	"github.com/flaminshinjan/address.ai/services/supply/internal/service"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083" // Default port for supply service
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	// Initialize database connection
	db, err := db.Connect(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	supplierRepo := repository.NewSupplierRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)
	purchaseRepo := repository.NewPurchaseRepository(db)

	// Initialize services
	supplierService := service.NewSupplierService(supplierRepo)
	inventoryService := service.NewInventoryService(inventoryRepo)
	purchaseService := service.NewPurchaseService(purchaseRepo, supplierRepo, inventoryRepo)

	// Initialize handlers
	h := handler.NewHandler(supplierService, inventoryService, purchaseService, jwtSecret)

	// Initialize router
	router := mux.NewRouter()

	// Register API routes
	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	h.RegisterRoutes(apiRouter)

	// Add health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Supply service is healthy"))
	}).Methods("GET")

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Supply service starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
