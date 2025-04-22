package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the service
type Config struct {
	DBHost      string
	DBPort      int
	DBUser      string
	DBPassword  string
	DBName      string
	JWTSecret   string
	ServicePort int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	servicePort, err := strconv.Atoi(getEnv("SERVICE_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVICE_PORT: %w", err)
	}

	return &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      dbPort,
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "postgres"),
		DBName:      getEnv("DB_NAME", "hotel_db"),
		JWTSecret:   getEnv("JWT_SECRET", "default_jwt_secret"),
		ServicePort: servicePort,
	}, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetDBConnString returns the database connection string
func (c *Config) GetDBConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}
