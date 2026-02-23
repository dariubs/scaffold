package middleware

import (
	"net/http"

	"github.com/dariubs/scaffold/app/utils"
	"github.com/gin-gonic/gin"
)

// Recovery middleware catches panics and logs them
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		utils.Logger.Error("Panic recovered",
			"error", recovered,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
			"headers", c.Request.Header,
		)

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	})
}
