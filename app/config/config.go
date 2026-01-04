package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database struct {
		DSN string
	}
	Session struct {
		Secret string
	}
	Server struct {
		Port      string
		AdminPort string
	}
	GoogleOAuth struct {
		ClientID     string
		ClientSecret string
		RedirectURL  string
	}
	CloudflareR2 struct {
		AccountID       string
		AccessKeyID     string
		SecretAccessKey string
		Bucket          string
		Region          string
	}
}

var C *Config

func Load() error {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	C = &Config{}

	// Database configuration
	C.Database.DSN = os.Getenv("DB_DSN")
	if C.Database.DSN == "" {
		return fmt.Errorf("DB_DSN is required")
	}

	// Session configuration
	C.Session.Secret = os.Getenv("SESSION_SECRET")
	if C.Session.Secret == "" {
		return fmt.Errorf("SESSION_SECRET is required")
	}

	// Server configuration
	C.Server.Port = os.Getenv("PORT")
	if C.Server.Port == "" {
		C.Server.Port = "3782"
	}

	C.Server.AdminPort = os.Getenv("ADMIN_PORT")
	if C.Server.AdminPort == "" {
		C.Server.AdminPort = "3781"
	}

	// Google OAuth configuration (optional)
	C.GoogleOAuth.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
	C.GoogleOAuth.ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	C.GoogleOAuth.RedirectURL = os.Getenv("GOOGLE_REDIRECT_URL")

	// Cloudflare R2 configuration (optional)
	C.CloudflareR2.AccountID = os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	C.CloudflareR2.AccessKeyID = os.Getenv("CLOUDFLARE_ACCESS_KEY_ID")
	C.CloudflareR2.SecretAccessKey = os.Getenv("CLOUDFLARE_SECRET_ACCESS_KEY")
	C.CloudflareR2.Bucket = os.Getenv("CLOUDFLARE_R2_BUCKET")
	C.CloudflareR2.Region = os.Getenv("CLOUDFLARE_R2_REGION")
	if C.CloudflareR2.Region == "" {
		C.CloudflareR2.Region = "auto"
	}

	return nil
}
