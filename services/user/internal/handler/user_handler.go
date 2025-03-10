package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/flaminshinjan/address.ai/pkg/common/auth"
	"github.com/flaminshinjan/address.ai/pkg/common/response"
	"github.com/flaminshinjan/address.ai/services/user/internal/model"
	"github.com/flaminshinjan/address.ai/services/user/internal/service"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	service   *service.UserService
	jwtSecret string
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(service *service.UserService, jwtSecret string) *UserHandler {
	return &UserHandler{
		service:   service,
		jwtSecret: jwtSecret,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.UserRegistration true "User Registration"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var reg model.UserRegistration
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	user, token, err := h.service.Register(reg)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	responseData := map[string]interface{}{
		"user":  user,
		"token": token,
	}

	response.Created(w, "User registered successfully", responseData)
}

// Login handles user login
// @Summary Login a user
// @Description Login a user with username and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.UserLogin true "User Login"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var login model.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	user, token, err := h.service.Login(login)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	responseData := map[string]interface{}{
		"user":  user,
		"token": token,
	}

	response.Success(w, "User logged in successfully", responseData)
}

// GetProfile handles getting the user's profile
// @Summary Get user profile
// @Description Get the profile of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /profile [get]
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := r.Context().Value("user_id").(string)

	user, err := h.service.GetByID(userID)
	if err != nil {
		response.NotFound(w, "User not found")
		return
	}

	response.Success(w, "User profile retrieved successfully", user)
}

// UpdateProfile handles updating the user's profile
// @Summary Update user profile
// @Description Update the profile of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body model.User true "User Update"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /profile [put]
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := r.Context().Value("user_id").(string)

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	// Ensure user ID matches
	user.ID = userID

	updatedUser, err := h.service.Update(userID, user)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, "User profile updated successfully", updatedUser)
}

// UpdatePassword handles updating the user's password
// @Summary Update user password
// @Description Update the password of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param passwords body map[string]string true "Password Update"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /profile/password [put]
func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := r.Context().Value("user_id").(string)

	var passwords map[string]string
	if err := json.NewDecoder(r.Body).Decode(&passwords); err != nil {
		response.BadRequest(w, "Invalid request payload")
		return
	}

	currentPassword, ok := passwords["current_password"]
	if !ok {
		response.BadRequest(w, "Current password is required")
		return
	}

	newPassword, ok := passwords["new_password"]
	if !ok {
		response.BadRequest(w, "New password is required")
		return
	}

	if err := h.service.UpdatePassword(userID, currentPassword, newPassword); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, "Password updated successfully", nil)
}

// GetUser handles getting a user by ID
// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.service.GetByID(id)
	if err != nil {
		response.NotFound(w, "User not found")
		return
	}

	response.Success(w, "User retrieved successfully", user)
}

// ListUsers handles listing all users
// @Summary List users
// @Description List all users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users [get]
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
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

	users, err := h.service.List(limit, offset)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.Success(w, "Users retrieved successfully", users)
}

// DeleteUser handles deleting a user
// @Summary Delete user
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.Delete(id); err != nil {
		response.NotFound(w, "User not found")
		return
	}

	response.Success(w, "User deleted successfully", nil)
}

// RegisterRoutes registers the routes for the user handler
func (h *UserHandler) RegisterRoutes(router *mux.Router) {
	// Public routes
	router.HandleFunc("/register", h.Register).Methods("POST")
	router.HandleFunc("/login", h.Login).Methods("POST")

	// Protected routes
	protected := router.PathPrefix("").Subrouter()
	protected.Use(func(next http.Handler) http.Handler {
		return auth.Middleware(h.jwtSecret, next)
	})

	protected.HandleFunc("/profile", h.GetProfile).Methods("GET")
	protected.HandleFunc("/profile", h.UpdateProfile).Methods("PUT")
	protected.HandleFunc("/profile/password", h.UpdatePassword).Methods("PUT")

	// Admin routes
	adminRouter := protected.PathPrefix("/users").Subrouter()
	adminRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := r.Context().Value("role").(string)
			if role != "admin" {
				response.Forbidden(w, "Admin access required")
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	adminRouter.HandleFunc("", h.ListUsers).Methods("GET")
	adminRouter.HandleFunc("/{id}", h.GetUser).Methods("GET")
	adminRouter.HandleFunc("/{id}", h.DeleteUser).Methods("DELETE")
}
