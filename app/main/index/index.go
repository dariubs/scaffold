package main

import (
	"log"
	"os"

	"github.com/dariubs/scaffold/app/database"
	"github.com/dariubs/scaffold/app/handlers/index"
	"github.com/dariubs/scaffold/app/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database.InitDB()

	// Initialize R2 service
	r2Service, err := utils.NewR2Service()
	if err != nil {
		log.Printf("Warning: R2 service not available: %v", err)
		r2Service = nil
	}

	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("views/index/*")

	// Session middleware
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "devsecret123" // fallback for dev
	}
	r.Use(sessions.Sessions("scaffoldsession", cookie.NewStore([]byte(sessionSecret))))

	// Routes
	r.GET("/", index.Home(database.DB))
	r.GET("/login", index.LoginForm())
	r.POST("/login", index.Login(database.DB))
	r.GET("/register", index.RegisterForm())
	r.POST("/register", index.Register(database.DB))
	r.GET("/logout", index.Logout())
	r.GET("/profile", index.Profile(database.DB))

	// Google OAuth routes
	r.GET("/auth/google", index.GoogleLogin())
	r.GET("/auth/google/callback", index.GoogleCallback(database.DB))

	// File upload routes (only if R2 service is available)
	if r2Service != nil {
		r.POST("/upload/profile-image", index.UploadProfileImage(database.DB, r2Service))
		r.POST("/upload/image", index.UploadImage(database.DB, r2Service))
		r.POST("/delete/image", index.DeleteImage(database.DB, r2Service))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3782"
	}

	log.Printf("Starting server on port %s", port)
	r.Run(":" + port)
}
