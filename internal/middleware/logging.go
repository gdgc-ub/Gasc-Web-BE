package middleware

import (
	"ProjectGolang/pkg/log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// LoggerConfig creates a Fiber middleware for structured request logging
func LoggerConfig() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		path := c.Path()
		method := c.Method()

		// Process request
		err := c.Next()

		// Calculate latency
		latency := time.Since(start)
		status := c.Response().StatusCode()

		// Determine log level based on status code
		logFields := log.Fields{
			"method":     method,
			"path":       path,
			"status":     status,
			"latency_ms": latency.Milliseconds(),
			"ip":         c.IP(),
			"host":       c.Hostname(),
			"user_agent": c.Get("User-Agent"),
			"referer":    c.Get("Referer"),
		}

		// Get request body if available
		if c.Request().Body() != nil && len(c.Request().Body()) > 0 {
			// Only log request body for non-file uploads or sensitive routes
			// You might want to filter out sensitive routes
			logFields["request_body"] = string(c.Request().Body())
		}

		// Log based on status code
		if status >= 500 {
			log.Error(logFields, "Server error")
		} else if status >= 400 {
			log.Warn(logFields, "Client error")
		} else {
			log.Info(logFields, "Success")
		}

		return err
	}
}

// For use with the existing middleware interface
func newLoggingMiddleware(logger *logrus.Logger) *loggingMiddleware {
	return &loggingMiddleware{
		logger: logger,
	}
}

type loggingMiddleware struct {
	logger *logrus.Logger
}

func (m *middleware) NewLoggingMiddleware(ctx *fiber.Ctx) error {
	start := time.Now()

	// Process request
	err := ctx.Next()

	// Calculate request processing time
	latency := time.Since(start)

	// Log request details
	m.log.WithFields(logrus.Fields{
		"client_ip":     ctx.IP(),
		"method":        ctx.Method(),
		"path":          ctx.Path(),
		"status":        ctx.Response().StatusCode(),
		"latency_ms":    latency.Milliseconds(),
		"request_body":  string(ctx.Request().Body()),
		"response_size": len(ctx.Response().Body()),
		"user_agent":    ctx.Get("User-Agent"),
	}).Info("HTTP Request")

	return err
}
