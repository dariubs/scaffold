package middleware

import (
	"time"

	"github.com/dariubs/scaffold/app/utils"
	"github.com/gin-gonic/gin"
)

// RequestLogger logs HTTP requests with structured logging
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Build log entry
		entry := utils.Logger.WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"query":      raw,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"latency":    latency.String(),
			"latency_ms": latency.Milliseconds(),
		})

		// Log error if status is 500 or above
		if c.Writer.Status() >= 500 {
			entry.Error("HTTP Request Error")
		} else {
			entry.Info("HTTP Request")
		}
	}
}

