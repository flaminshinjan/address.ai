package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/flaminshinjan/address.ai/pkg/common/db"
	"github.com/flaminshinjan/address.ai/pkg/common/logger"
	"github.com/flaminshinjan/address.ai/services/food/internal/handler"
	"github.com/flaminshinjan/address.ai/services/food/internal/repository"
	"github.com/flaminshinjan/address.ai/services/food/internal/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// @title Food Management Service API
// @version 1.0
// @description This is the food management service API for the hotel management system
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8083
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082" // Default port for food service
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
	menuRepo := repository.NewMenuRepository(database)
	orderRepo := repository.NewOrderRepository(database)

	// Initialize services
	menuService := service.NewMenuService(menuRepo)
	orderService := service.NewOrderService(orderRepo, menuRepo)

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
	h := handler.NewHandler(menuService, orderService, jwtSecret)

	// Register routes
	api := e.Group("/api/v1")
	h.RegisterRoutes(api)

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "Food service is healthy")
	})

	// Start server
	log.Printf("Food service starting on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
