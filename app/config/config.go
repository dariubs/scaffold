package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func isTruthy(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "true" || s == "1" || s == "yes"
}

type Config struct {
	Database struct {
		DSN string
	}
	Session struct {
		Secret string
	}
	Server struct {
		Port      string
		AdminPath string
	}
	Login struct {
		PasswordEnabled bool
		GoogleEnabled   bool
		GitHubEnabled   bool
		LinkedInEnabled bool
		XEnabled        bool
	}
	GoogleOAuth struct {
		ClientID     string
		ClientSecret string
		RedirectURL  string
	}
	GitHubOAuth struct {
		ClientID     string
		ClientSecret string
		RedirectURL  string
	}
	LinkedInOAuth struct {
		ClientID     string
		ClientSecret string
		RedirectURL  string
	}
	XOAuth struct {
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
	Resend struct {
		APIKey string
		From   string
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

	C.Server.AdminPath = os.Getenv("ADMIN_BASE_PATH")
	if C.Server.AdminPath == "" {
		C.Server.AdminPath = "admin"
	}

	// Login method enable flags (optional)
	if v := os.Getenv("LOGIN_PASSWORD_ENABLED"); v != "" {
		C.Login.PasswordEnabled = isTruthy(v)
	} else {
		C.Login.PasswordEnabled = true
	}
	if v := os.Getenv("LOGIN_GOOGLE_ENABLED"); v != "" {
		C.Login.GoogleEnabled = isTruthy(v)
	} else {
		C.Login.GoogleEnabled = true
	}
	C.Login.GitHubEnabled = isTruthy(os.Getenv("LOGIN_GITHUB_ENABLED"))
	C.Login.LinkedInEnabled = isTruthy(os.Getenv("LOGIN_LINKEDIN_ENABLED"))
	C.Login.XEnabled = isTruthy(os.Getenv("LOGIN_X_ENABLED"))

	// Google OAuth configuration (optional)
	C.GoogleOAuth.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
	C.GoogleOAuth.ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	C.GoogleOAuth.RedirectURL = os.Getenv("GOOGLE_REDIRECT_URL")

	// GitHub OAuth configuration (optional)
	C.GitHubOAuth.ClientID = os.Getenv("GITHUB_CLIENT_ID")
	C.GitHubOAuth.ClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
	C.GitHubOAuth.RedirectURL = os.Getenv("GITHUB_REDIRECT_URL")

	// LinkedIn OAuth configuration (optional)
	C.LinkedInOAuth.ClientID = os.Getenv("LINKEDIN_CLIENT_ID")
	C.LinkedInOAuth.ClientSecret = os.Getenv("LINKEDIN_CLIENT_SECRET")
	C.LinkedInOAuth.RedirectURL = os.Getenv("LINKEDIN_REDIRECT_URL")

	// X (Twitter) OAuth configuration (optional)
	C.XOAuth.ClientID = os.Getenv("X_CLIENT_ID")
	C.XOAuth.ClientSecret = os.Getenv("X_CLIENT_SECRET")
	C.XOAuth.RedirectURL = os.Getenv("X_REDIRECT_URL")

	// Cloudflare R2 configuration (optional)
	C.CloudflareR2.AccountID = os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	C.CloudflareR2.AccessKeyID = os.Getenv("CLOUDFLARE_ACCESS_KEY_ID")
	C.CloudflareR2.SecretAccessKey = os.Getenv("CLOUDFLARE_SECRET_ACCESS_KEY")
	C.CloudflareR2.Bucket = os.Getenv("CLOUDFLARE_R2_BUCKET")
	C.CloudflareR2.Region = os.Getenv("CLOUDFLARE_R2_REGION")
	if C.CloudflareR2.Region == "" {
		C.CloudflareR2.Region = "auto"
	}

	// Resend email configuration (optional)
	C.Resend.APIKey = os.Getenv("RESEND_API_KEY")
	C.Resend.From = os.Getenv("RESEND_FROM")

	return nil
}

// OAuthGoogleEnabled returns true if Google login is enabled and configured.
func (c *Config) OAuthGoogleEnabled() bool {
	return c.Login.GoogleEnabled && c.GoogleOAuth.ClientID != "" && c.GoogleOAuth.RedirectURL != ""
}

// OAuthGitHubEnabled returns true if GitHub login is enabled and configured.
func (c *Config) OAuthGitHubEnabled() bool {
	return c.Login.GitHubEnabled && c.GitHubOAuth.ClientID != "" && c.GitHubOAuth.RedirectURL != ""
}

// OAuthLinkedInEnabled returns true if LinkedIn login is enabled and configured.
func (c *Config) OAuthLinkedInEnabled() bool {
	return c.Login.LinkedInEnabled && c.LinkedInOAuth.ClientID != "" && c.LinkedInOAuth.RedirectURL != ""
}

// OAuthXEnabled returns true if X (Twitter) login is enabled and configured.
func (c *Config) OAuthXEnabled() bool {
	return c.Login.XEnabled && c.XOAuth.ClientID != "" && c.XOAuth.RedirectURL != ""
}
