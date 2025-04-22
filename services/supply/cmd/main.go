package main

import (
	"log"
	"os"

	"github.com/flaminshinjan/address.ai/pkg/common/db"
	"github.com/flaminshinjan/address.ai/pkg/common/logger"
	"github.com/flaminshinjan/address.ai/services/supply/internal/handler"
	"github.com/flaminshinjan/address.ai/services/supply/internal/repository"
	"github.com/flaminshinjan/address.ai/services/supply/internal/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	migrationsPath := "/app/migrations"
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

	// Initialize Echo
	e := echo.New()
	e.HideBanner = false

	// Configure logger
	logger.Configure(e)
	logger.SetLogLevel(e, logLevel)

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize handlers
	h := handler.NewHandler(supplierService, inventoryService, purchaseService, jwtSecret)

	// Register routes
	api := e.Group("/api/v1")
	h.RegisterRoutes(api)

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "Supply service is healthy")
	})

	// Start server
	log.Printf("Supply service starting on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
