package index

import (
	"context"
	"net/http"

	"github.com/dariubs/scaffold/app/config"
	"github.com/dariubs/scaffold/app/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleoauth2 "google.golang.org/api/oauth2/v2"
	"gorm.io/gorm"
)

func getGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.C.GoogleOAuth.ClientID,
		ClientSecret: config.C.GoogleOAuth.ClientSecret,
		RedirectURL:  config.C.GoogleOAuth.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func GoogleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		googleOauthConfig := getGoogleOAuthConfig()
		if googleOauthConfig.ClientID == "" {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"Title": "Login",
				"Error": "Google OAuth is not configured",
			})
			return
		}

		// Generate state token for CSRF protection
		state := uuid.New().String()
		session := sessions.Default(c)
		session.Set("oauth_state", state)
		if err := session.Save(); err != nil {
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{
				"Title": "Login",
				"Error": "Failed to save session",
			})
			return
		}

		url := googleOauthConfig.AuthCodeURL(state)
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func GoogleCallback(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		googleOauthConfig := getGoogleOAuthConfig()
		if googleOauthConfig.ClientID == "" {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"Title": "Login",
				"Error": "Google OAuth is not configured",
			})
			return
		}

		code := c.Query("code")
		state := c.Query("state")

		if code == "" {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"Title": "Login",
				"Error": "Authorization code not received",
			})
			return
		}

		// Validate state token (CSRF protection)
		session := sessions.Default(c)
		savedState := session.Get("oauth_state")
		if savedState == nil || savedState.(string) != state {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"Title": "Login",
				"Error": "Invalid state token",
			})
			return
		}

		// Clear state from session
		session.Delete("oauth_state")
		session.Save()

		// Exchange code for token
		token, err := googleOauthConfig.Exchange(context.Background(), code)
		if err != nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"Title": "Login",
				"Error": "Failed to exchange token",
			})
			return
		}

		// Get user info from Google
		client := googleOauthConfig.Client(context.Background(), token)
		oauth2Service, err := googleoauth2.New(client)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{
				"Title": "Login",
				"Error": "Failed to create OAuth service",
			})
			return
		}

		userInfo, err := oauth2Service.Userinfo.Get().Do()
		if err != nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"Title": "Login",
				"Error": "Failed to get user info",
			})
			return
		}

		// Check if user exists
		var user model.User
		err = db.Where("google_id = ? OR email = ?", userInfo.Id, userInfo.Email).First(&user).Error

		if err != nil {
			// User doesn't exist, create new one
			user = model.User{
				Username:    userInfo.Email, // Use email as username for OAuth users
				Email:       userInfo.Email,
				Password:    "", // No password for OAuth users
				Name:        userInfo.Name,
				AvatarURL:   userInfo.Picture,
				GoogleID:    userInfo.Id,
				LoginMethod: "google",
			}

			if err := db.Create(&user).Error; err != nil {
				c.HTML(http.StatusInternalServerError, "login.html", gin.H{
					"Title": "Login",
					"Error": "Failed to create user account",
				})
				return
			}
		} else {
			// User exists, update Google ID if not set
			if user.GoogleID == "" {
				user.GoogleID = userInfo.Id
				user.LoginMethod = "google"
				user.AvatarURL = userInfo.Picture
				user.Name = userInfo.Name
				db.Save(&user)
			}
		}

		// Set session
		session := sessions.Default(c)
		session.Set("user_id", user.Model.ID)
		session.Save()

		c.Redirect(http.StatusFound, "/")
	}
}
