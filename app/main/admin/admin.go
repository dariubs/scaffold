package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dariubs/scaffold/app/config"
	"github.com/dariubs/scaffold/app/database"
	"github.com/dariubs/scaffold/app/handlers/admin"
	"github.com/dariubs/scaffold/app/middleware"
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

	r := gin.Default()
	r.LoadHTMLGlob("views/admin/*")

	// Session middleware (must use same secret as main app for session sharing)
	r.Use(sessions.Sessions("scaffoldsession", cookie.NewStore([]byte(config.C.Session.Secret))))

	// Admin routes - require admin authentication
	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.RequireAdmin(database.DB))
	{
		adminGroup.GET("/", admin.AdminHome())
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + config.C.Server.AdminPort,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting admin server on port %s", config.C.Server.AdminPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Admin server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down admin server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Admin server forced to shutdown:", err)
	}

	log.Println("Admin server exited")
}
