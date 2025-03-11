package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func main() {
	// Replace with your actual JWT secret
	jwtSecret := "your_jwt_secret_here"

	// Create the Claims
	claims := jwt.MapClaims{
		"user_id":  "admin123",
		"username": "admin",
		"role":     "admin",
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		return
	}

	fmt.Printf("Generated admin token: %s\n", tokenString)
}