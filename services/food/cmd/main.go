package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/flaminshinjan/address.ai/pkg/common/config"
	"github.com/flaminshinjan/address.ai/pkg/common/db"
	"github.com/flaminshinjan/address.ai/services/food/internal/handler"
	"github.com/flaminshinjan/address.ai/services/food/internal/repository"
	"github.com/flaminshinjan/address.ai/services/food/internal/service"
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
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	database, err := db.Connect(cfg.GetDBConnString())
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

	// Initialize handlers
	menuHandler := handler.NewMenuHandler(menuService, cfg.JWTSecret)
	orderHandler := handler.NewOrderHandler(orderService, cfg.JWTSecret)

	// Initialize router
	router := mux.NewRouter()

	// Register routes
	menuHandler.RegisterRoutes(router)
	orderHandler.RegisterRoutes(router)

	// Swagger documentation
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)

	// Start server
	port := fmt.Sprintf(":%d", cfg.ServicePort)
	log.Printf("Food management service starting on port %s", port)
	if err := http.ListenAndServe(port, corsHandler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
