package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig holds CORS configuration.
type CORSConfig struct {
	AllowOrigins string
}

// CORS returns a middleware that handles CORS.
func CORS(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Determine allowed origin
		allowedOrigin := ""
		if config.AllowOrigins == "*" {
			allowedOrigin = "*"
		} else {
			origins := strings.Split(config.AllowOrigins, ",")
			for _, o := range origins {
				if strings.TrimSpace(o) == origin {
					allowedOrigin = origin
					break
				}
			}
		}

		if allowedOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
			c.Header("Access-Control-Expose-Headers", "X-Request-ID")
			c.Header("Access-Control-Max-Age", "86400")
		}

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
