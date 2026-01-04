package middleware

import (
	"net/http"

	"github.com/dariubs/scaffold/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Recovery middleware catches panics and logs them
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		utils.Logger.WithFields(logrus.Fields{
			"error":   recovered,
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"ip":      c.ClientIP(),
			"headers": c.Request.Header,
		}).Error("Panic recovered")

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	})
}

