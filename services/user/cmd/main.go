package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/flaminshinjan/address.ai/pkg/common/db"
	"github.com/flaminshinjan/address.ai/pkg/common/logger"
	"github.com/flaminshinjan/address.ai/services/user/internal/handler"
	"github.com/flaminshinjan/address.ai/services/user/internal/repository"
	"github.com/flaminshinjan/address.ai/services/user/internal/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// @title User Service API
// @version 1.0
// @description This is the user service API for the hotel management system
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
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
		port = "8081" // Default port for user service
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
	userRepo := repository.NewUserRepository(database)

	// Initialize services
	userService := service.NewUserService(userRepo, jwtSecret)

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
	userHandler := handler.NewUserHandler(userService, jwtSecret)

	// Register routes
	api := e.Group("/api/v1")
	userHandler.RegisterRoutes(api)

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "User service is healthy")
	})

	// Start server
	log.Printf("User service starting on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
