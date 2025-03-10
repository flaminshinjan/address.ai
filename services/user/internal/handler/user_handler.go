package handler

import (
	"net/http"
	"strconv"

	"github.com/flaminshinjan/address.ai/pkg/common/auth"
	"github.com/flaminshinjan/address.ai/services/user/internal/model"
	"github.com/flaminshinjan/address.ai/services/user/internal/service"
	"github.com/labstack/echo/v4"
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
func (h *UserHandler) Register(c echo.Context) error {
	var req model.UserRegistration
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request payload",
		})
	}

	user, token, err := h.service.Register(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "User registered successfully",
		"data":    user,
		"token":   token,
	})
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
func (h *UserHandler) Login(c echo.Context) error {
	var req model.UserLogin
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request payload",
		})
	}

	user, token, err := h.service.Login(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"data":    user,
		"token":   token,
	})
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
func (h *UserHandler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id").(string)

	user, err := h.service.GetByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "User not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Profile retrieved successfully",
		"data":    user,
	})
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
func (h *UserHandler) UpdateProfile(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var user model.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request payload",
		})
	}

	user.ID = userID
	updatedUser, err := h.service.Update(userID, user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Profile updated successfully",
		"data":    updatedUser,
	})
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
func (h *UserHandler) UpdatePassword(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request payload",
		})
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Current password and new password are required",
		})
	}

	err := h.service.UpdatePassword(userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Password updated successfully",
	})
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
func (h *UserHandler) GetUser(c echo.Context) error {
	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"success": false,
			"error":   "Admin access required",
		})
	}

	id := c.Param("id")
	user, err := h.service.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "User not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "User retrieved successfully",
		"data":    user,
	})
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
func (h *UserHandler) ListUsers(c echo.Context) error {
	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"success": false,
			"error":   "Admin access required",
		})
	}

	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

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
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to retrieve users",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Users retrieved successfully",
		"data":    users,
	})
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
func (h *UserHandler) DeleteUser(c echo.Context) error {
	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"success": false,
			"error":   "Admin access required",
		})
	}

	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "User not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "User deleted successfully",
	})
}

// RegisterRoutes registers the routes for the user handler
func (h *UserHandler) RegisterRoutes(g *echo.Group) {
	// Public routes
	g.POST("/auth/register", h.Register)
	g.POST("/auth/login", h.Login)

	// Protected routes
	users := g.Group("/users")
	users.Use(h.authMiddleware)

	users.GET("/profile", h.GetProfile)
	users.PUT("/profile", h.UpdateProfile)
	users.PUT("/password", h.UpdatePassword)

	// Admin routes
	users.GET("", h.ListUsers)
	users.GET("/:id", h.GetUser)
	users.DELETE("/:id", h.DeleteUser)
}

// authMiddleware is a middleware to check if the user is authenticated
func (h *UserHandler) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"error":   "Authorization token is required",
			})
		}

		// Remove "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		claims, err := auth.ValidateToken(token, h.jwtSecret)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"error":   "Invalid or expired token",
			})
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		return next(c)
	}
}
