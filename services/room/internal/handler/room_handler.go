package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/yourusername/hotel-management/pkg/common/auth"
	"github.com/yourusername/hotel-management/pkg/common/response"
	"github.com/yourusername/hotel-management/services/room/internal/model"
	"github.com/yourusername/hotel-management/services/room/internal/service"
)

// RoomHandler handles HTTP requests for rooms
type RoomHandler struct {
	roomService    *service.RoomService
	bookingService *service.BookingService
	jwtSecret      string
}

// NewRoomHandler creates a new RoomHandler
func NewRoomHandler(roomService *service.RoomService, bookingService *service.BookingService, jwtSecret string) *RoomHandler {
	return &RoomHandler{
		roomService:    roomService,
		bookingService: bookingService,
		jwtSecret:      jwtSecret,
	}
}

// CreateRoom handles creating a new room
// @Summary Create a new room
// @Description Create a new room with the provided details
// @Tags rooms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param room body model.Room true "Room"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /rooms [post]
func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	role := r.Context().Value("role").(string)
	if role != "admin" {
		response.Forbidden(w, "Admin access required")
		return
	}

	var room model.Room
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	createdRoom, err := h.roomService.CreateRoom(room)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, "Room created successfully", createdRoom)
}

// GetRoom handles getting a room by ID
// @Summary Get a room by ID
// @Description Get a room by its ID
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path string true "Room ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /rooms/{id} [get]
func (h *RoomHandler) GetRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	room, err := h.roomService.GetRoomByID(id)
	if err != nil {
		response.NotFound(w, "Room not found")
		return
	}

	response.Success(w, "Room retrieved successfully", room)
}

// UpdateRoom handles updating a room
// @Summary Update a room
// @Description Update a room with the provided details
// @Tags rooms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Room ID"
// @Param room body model.Room true "Room"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /rooms/{id} [put]
func (h *RoomHandler) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	role := r.Context().Value("role").(string)
	if role != "admin" {
		response.Forbidden(w, "Admin access required")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var room model.Room
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	updatedRoom, err := h.roomService.UpdateRoom(id, room)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, "Room updated successfully", updatedRoom)
}

// DeleteRoom handles deleting a room
// @Summary Delete a room
// @Description Delete a room by its ID
// @Tags rooms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Room ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /rooms/{id} [delete]
func (h *RoomHandler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	role := r.Context().Value("role").(string)
	if role != "admin" {
		response.Forbidden(w, "Admin access required")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.roomService.DeleteRoom(id); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, "Room deleted successfully", nil)
}

// ListRooms handles listing all rooms
// @Summary List all rooms
// @Description List all rooms with pagination
// @Tags rooms
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /rooms [get]
func (h *RoomHandler) ListRooms(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10 // Default limit
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0 // Default offset
	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	rooms, err := h.roomService.ListRooms(limit, offset)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.Success(w, "Rooms retrieved successfully", rooms)
}

// ListAvailableRooms handles listing all available rooms for the given dates
// @Summary List available rooms
// @Description List all available rooms for the given dates with pagination
// @Tags rooms
// @Accept json
// @Produce json
// @Param start_date query string true "Start Date (YYYY-MM-DD)"
// @Param end_date query string true "End Date (YYYY-MM-DD)"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /rooms/available [get]
func (h *RoomHandler) ListAvailableRooms(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Parse dates
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		response.BadRequest(w, "Invalid start date format (YYYY-MM-DD)")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		response.BadRequest(w, "Invalid end date format (YYYY-MM-DD)")
		return
	}

	limit := 10 // Default limit
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0 // Default offset
	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	rooms, err := h.roomService.ListAvailableRooms(startDate, endDate, limit, offset)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, "Available rooms retrieved successfully", rooms)
}

// UpdateRoomStatus handles updating a room's status
// @Summary Update room status
// @Description Update a room's status
// @Tags rooms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Room ID"
// @Param status body map[string]string true "Status"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /rooms/{id}/status [put]
func (h *RoomHandler) UpdateRoomStatus(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	role := r.Context().Value("role").(string)
	if role != "admin" {
		response.Forbidden(w, "Admin access required")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var statusData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&statusData); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	status, ok := statusData["status"]
	if !ok {
		response.BadRequest(w, "Status is required")
		return
	}

	if err := h.roomService.UpdateRoomStatus(id, status); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, "Room status updated successfully", nil)
}

// RegisterRoutes registers the routes for the room handler
func (h *RoomHandler) RegisterRoutes(router *mux.Router) {
	// Public routes
	router.HandleFunc("/rooms", h.ListRooms).Methods("GET")
	router.HandleFunc("/rooms/available", h.ListAvailableRooms).Methods("GET")
	router.HandleFunc("/rooms/{id}", h.GetRoom).Methods("GET")

	// Protected routes
	protected := router.PathPrefix("").Subrouter()
	protected.Use(func(next http.Handler) http.Handler {
		return auth.Middleware(h.jwtSecret, next)
	})

	// Admin routes
	adminRouter := protected.PathPrefix("/rooms").Subrouter()
	adminRouter.HandleFunc("", h.CreateRoom).Methods("POST")
	adminRouter.HandleFunc("/{id}", h.UpdateRoom).Methods("PUT")
	adminRouter.HandleFunc("/{id}", h.DeleteRoom).Methods("DELETE")
	adminRouter.HandleFunc("/{id}/status", h.UpdateRoomStatus).Methods("PUT")
}
