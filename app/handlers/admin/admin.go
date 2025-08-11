package admin

import "github.com/gin-gonic/gin"

func AdminHome() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(200, "admin.home.html", gin.H{
			"Title": "Admin Dashboard",
		})
	}
}
