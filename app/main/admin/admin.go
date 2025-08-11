package main

import (
	"log"
	"os"

	"github.com/dariubs/scaffold/app/database"
	"github.com/dariubs/scaffold/app/handlers/admin"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database.InitDB()

	r := gin.Default()
	r.LoadHTMLGlob("views/admin/*")

	// Basic auth middleware
	username := os.Getenv("ADMIN_USER")
	if username == "" {
		username = "admin"
	}
	password := os.Getenv("ADMIN_PASS")
	if password == "" {
		password = "admin123"
	}

	// Admin routes
	adminGroup := r.Group("/admin")
	adminGroup.Use(gin.BasicAuth(gin.Accounts{
		username: password,
	}))
	{
		adminGroup.GET("/", admin.AdminHome())
	}

	port := os.Getenv("ADMIN_PORT")
	if port == "" {
		port = "3781"
	}

	log.Printf("Starting admin server on port %s", port)
	r.Run(":" + port)
}
