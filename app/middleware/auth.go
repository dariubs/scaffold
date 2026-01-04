package middleware

import (
	"net/http"

	"github.com/dariubs/scaffold/app/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RequireAuth checks if the user is authenticated
func RequireAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		if userID == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Check if user exists
		var user model.User
		if err := db.First(&user, userID).Error; err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Store user in context for use in handlers
		c.Set("user", user)
		c.Set("user_id", userID)

		c.Next()
	}
}
