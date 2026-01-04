package admin

import (
	"net/http"

	"github.com/dariubs/scaffold/app/model"
	"github.com/gin-gonic/gin"
)

func AdminHome() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (set by admin middleware)
		user, exists := c.Get("user")
		if !exists {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Title": "Error",
				"Error": "User not found in context",
			})
			return
		}

		adminUser, ok := user.(model.User)
		if !ok {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Title": "Error",
				"Error": "Invalid user type",
			})
			return
		}

		c.HTML(http.StatusOK, "admin.home.html", gin.H{
			"Title": "Admin Dashboard",
			"User":  adminUser,
		})
	}
}
