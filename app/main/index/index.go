package main

import (
	"log"

	"github.com/dariubs/scaffold/app/config"
	"github.com/dariubs/scaffold/app/database"
	"github.com/dariubs/scaffold/app/handlers/health"
	"github.com/dariubs/scaffold/app/handlers/index"
	"github.com/dariubs/scaffold/app/middleware"
	"github.com/dariubs/scaffold/app/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	err := config.Load()
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

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
	r.Use(sessions.Sessions("scaffoldsession", cookie.NewStore([]byte(config.C.Session.Secret))))

	// Health check routes (before other middleware)
	healthGroup := r.Group("")
	{
		healthGroup.GET("/health", health.Health())
		healthGroup.GET("/readiness", health.Readiness())
	}

	// Routes
	r.GET("/", index.Home(database.DB))
	r.GET("/login", index.LoginForm())
	r.POST("/login", index.Login(database.DB))
	r.GET("/register", index.RegisterForm())
	r.POST("/register", index.Register(database.DB))
	r.GET("/logout", index.Logout())

	// Protected routes
	protected := r.Group("")
	protected.Use(middleware.RequireAuth(database.DB))
	{
		protected.GET("/profile", index.Profile(database.DB))
	}

	// Google OAuth routes
	r.GET("/auth/google", index.GoogleLogin())
	r.GET("/auth/google/callback", index.GoogleCallback(database.DB))

	// File upload routes (only if R2 service is available, protected)
	if r2Service != nil {
		uploadGroup := r.Group("")
		uploadGroup.Use(middleware.RequireAuth(database.DB))
		{
			uploadGroup.POST("/upload/profile-image", index.UploadProfileImage(database.DB, r2Service))
			uploadGroup.POST("/upload/image", index.UploadImage(database.DB, r2Service))
			uploadGroup.POST("/delete/image", index.DeleteImage(database.DB, r2Service))
		}
	}

	log.Printf("Starting server on port %s", config.C.Server.Port)
	r.Run(":" + config.C.Server.Port)
}
