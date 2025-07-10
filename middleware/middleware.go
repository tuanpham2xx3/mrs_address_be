package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger returns a Gin middleware for logging HTTP requests.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Printf("[%s] %s %s %d %v %s",
			time.Now().Format("2006/01/02 15:04:05"),
			method,
			path,
			statusCode,
			latency,
			clientIP,
		)
	}
}

// CORS returns a gin middleware for CORS
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "43200") // 12 hours

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RateLimit returns a simple rate limiting middleware
func RateLimit() gin.HandlerFunc {
	// Simple in-memory rate limiter (for production, use Redis-based)
	return func(c *gin.Context) {
		// This is a placeholder - implement proper rate limiting as needed
		c.Next()
	}
}

// AdminAuth returns a middleware for admin authentication
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simple token-based auth (implement proper auth as needed)
		token := c.GetHeader("Authorization")

		// Example: Bearer admin-secret-token
		if token != "Bearer admin-secret-token" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized access",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Recovery returns a custom recovery middleware
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Internal server error",
				})
			}
		}()
		c.Next()
	}
}
