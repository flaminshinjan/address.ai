package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/flaminshinjan/address.ai/pkg/common/db"
	"github.com/flaminshinjan/address.ai/services/supply/internal/handler"
	"github.com/flaminshinjan/address.ai/services/supply/internal/repository"
	"github.com/flaminshinjan/address.ai/services/supply/internal/service"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084" // Default port for supply service
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info" // Default log level
	}

	// Initialize database connection
	database, err := db.Connect(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	migrationsPath := filepath.Join("migrations")
	if err := db.RunMigrations(database, migrationsPath); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	supplierRepo := repository.NewSupplierRepository(database)
	inventoryRepo := repository.NewInventoryRepository(database)
	purchaseRepo := repository.NewPurchaseRepository(database)

	// Initialize services
	supplierService := service.NewSupplierService(supplierRepo)
	inventoryService := service.NewInventoryService(inventoryRepo)
	purchaseService := service.NewPurchaseService(purchaseRepo, supplierRepo, inventoryRepo)

	// Initialize router
	router := mux.NewRouter()

	// Initialize handlers
	h := handler.NewHandler(supplierService, inventoryService, purchaseService, jwtSecret)

	// Register API routes
	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	h.RegisterRoutes(apiRouter)

	// Add health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Supply service is healthy"))
	}).Methods("GET")

	// CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)

	// Start server
	log.Printf("Supply service starting on port %s", port)
	if err := http.ListenAndServe(":"+port, corsHandler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
