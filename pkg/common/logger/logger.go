package logger

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

// Configure sets up the Echo logger middleware with custom configuration
func Configure(e *echo.Echo) {
	// Set log level based on environment
	e.Logger.SetLevel(log.INFO)

	// Configure logger middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} | ${status} | ${method} ${uri} | ${latency_human} | ${remote_ip} | ${error}\n",
	}))
}

// SetLogLevel sets the log level for the Echo instance
func SetLogLevel(e *echo.Echo, level string) {
	switch level {
	case "debug":
		e.Logger.SetLevel(log.DEBUG)
	case "info":
		e.Logger.SetLevel(log.INFO)
	case "warn":
		e.Logger.SetLevel(log.WARN)
	case "error":
		e.Logger.SetLevel(log.ERROR)
	default:
		e.Logger.SetLevel(log.INFO)
	}
}
