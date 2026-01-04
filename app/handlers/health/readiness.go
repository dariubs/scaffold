package health

import (
	"net/http"

	"github.com/dariubs/scaffold/app/database"
	"github.com/gin-gonic/gin"
)

// Readiness checks if the application is ready to serve traffic
func Readiness() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check database connection
		sqlDB, err := database.GetSQLDB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error":  "database connection unavailable",
			})
			return
		}

		// Ping database
		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error":  "database ping failed",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	}
}

