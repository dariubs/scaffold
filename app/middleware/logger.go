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

		latency := time.Since(start)
		attrs := []any{
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", path,
			"query", raw,
			"ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"latency", latency.String(),
			"latency_ms", latency.Milliseconds(),
		}

		if c.Writer.Status() >= 500 {
			utils.Logger.Error("HTTP Request Error", attrs...)
		} else {
			utils.Logger.Info("HTTP Request", attrs...)
		}
	}
}
