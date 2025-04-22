package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/shinjan/address.ai/services/room/config"
	"github.com/shinjan/address.ai/services/room/handlers"
)

func main() {
	// Initialize database
	config.InitDB()

	// Create router
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/api/v1/rooms", handlers.GetRooms).Methods("GET")
	r.HandleFunc("/api/v1/bookings", handlers.CreateBooking).Methods("POST")
	r.HandleFunc("/api/v1/bookings", handlers.GetBookings).Methods("GET")

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	// Start server
	log.Println("Starting server on :8082")
	log.Fatal(http.ListenAndServe(":8082", c.Handler(r)))
}
