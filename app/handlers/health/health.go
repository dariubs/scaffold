package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health returns a simple health check endpoint
func Health() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}

