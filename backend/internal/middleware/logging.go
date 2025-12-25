package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// RequestIDKey is the context key for the request ID.
const RequestIDKey = "requestId"

// Logging returns a middleware that logs HTTP requests.
func Logging(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate request ID
		requestID := uuid.New().String()
		c.Set(RequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status and errors
		status := c.Writer.Status()
		errors := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Build log event
		event := logger.Info()
		if status >= 400 && status < 500 {
			event = logger.Warn()
		} else if status >= 500 {
			event = logger.Error()
		}

		event.
			Str("requestId", requestID).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("query", c.Request.URL.RawQuery).
			Int("status", status).
			Dur("latency", latency).
			Str("clientIP", c.ClientIP()).
			Str("userAgent", c.Request.UserAgent())

		if errors != "" {
			event.Str("errors", errors)
		}

		event.Msg("request completed")
	}
}

// GetRequestID returns the request ID from the context.
func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get(RequestIDKey); exists {
		return id.(string)
	}
	return ""
}
