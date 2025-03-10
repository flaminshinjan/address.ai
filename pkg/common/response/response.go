package response

import (
	"encoding/json"
	"net/http"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// JSON sends a JSON response
func JSON(w http.ResponseWriter, statusCode int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Success sends a success response
func Success(w http.ResponseWriter, message string, data interface{}) {
	JSON(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a created response
func Created(w http.ResponseWriter, message string, data interface{}) {
	JSON(w, http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// BadRequest sends a bad request response
func BadRequest(w http.ResponseWriter, message string) {
	JSON(w, http.StatusBadRequest, Response{
		Success: false,
		Error:   message,
	})
}

// Unauthorized sends an unauthorized response
func Unauthorized(w http.ResponseWriter, message string) {
	JSON(w, http.StatusUnauthorized, Response{
		Success: false,
		Error:   message,
	})
}

// Forbidden sends a forbidden response
func Forbidden(w http.ResponseWriter, message string) {
	JSON(w, http.StatusForbidden, Response{
		Success: false,
		Error:   message,
	})
}

// NotFound sends a not found response
func NotFound(w http.ResponseWriter, message string) {
	JSON(w, http.StatusNotFound, Response{
		Success: false,
		Error:   message,
	})
}

// InternalServerError sends an internal server error response
func InternalServerError(w http.ResponseWriter, err error) {
	JSON(w, http.StatusInternalServerError, Response{
		Success: false,
		Error:   "Internal server error",
	})
}
