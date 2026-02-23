package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dariubs/scaffold/app/config"
	"github.com/dariubs/scaffold/app/database"
	"github.com/dariubs/scaffold/app/handlers/admin"
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

	// Initialize email service (Resend; optional)
	emailService, _ := utils.NewEmailService()

	r := gin.Default()

	// Load HTML templates from both index and admin directories
	t, err := template.New("").ParseGlob("views/index/*.html")
	if err != nil {
		log.Fatal("Failed to parse index templates:", err)
	}
	t, err = t.ParseGlob("views/admin/*.html")
	if err != nil {
		log.Fatal("Failed to parse admin templates:", err)
	}
	r.SetHTMLTemplate(t)

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
	r.POST("/register", index.Register(database.DB, emailService))
	r.GET("/logout", index.Logout())

	// Protected routes
	protected := r.Group("")
	protected.Use(middleware.RequireAuth(database.DB))
	{
		protected.GET("/profile", index.Profile(database.DB))
	}

	// OAuth routes (only for enabled providers)
	if config.C.OAuthGoogleEnabled() {
		r.GET("/auth/google", index.GoogleLogin())
		r.GET("/auth/google/callback", index.GoogleCallback(database.DB))
	}
	if config.C.OAuthGitHubEnabled() {
		r.GET("/auth/github", index.GitHubLogin())
		r.GET("/auth/github/callback", index.GitHubCallback(database.DB))
	}
	if config.C.OAuthLinkedInEnabled() {
		r.GET("/auth/linkedin", index.LinkedInLogin())
		r.GET("/auth/linkedin/callback", index.LinkedInCallback(database.DB))
	}
	if config.C.OAuthXEnabled() {
		r.GET("/auth/x", index.XLogin())
		r.GET("/auth/x/callback", index.XCallback(database.DB))
	}

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

	// Admin routes (mount at configurable base path)
	adminGroup := r.Group("/" + config.C.Server.AdminPath)
	adminGroup.Use(middleware.RequireAdmin(database.DB))
	{
		adminGroup.GET("/", admin.AdminHome())
	}

	srv := &http.Server{
		Addr:    ":" + config.C.Server.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Starting server on port %s", config.C.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
