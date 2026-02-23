package index

import (
	"net/http"

	"github.com/dariubs/scaffold/app/config"
	"github.com/dariubs/scaffold/app/model"
	"github.com/dariubs/scaffold/app/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func loginFormData() gin.H {
	return gin.H{
		"LoginPassword": config.C.Login.PasswordEnabled,
		"LoginGoogle":   config.C.OAuthGoogleEnabled(),
		"LoginGitHub":  config.C.OAuthGitHubEnabled(),
		"LoginLinkedIn": config.C.OAuthLinkedInEnabled(),
		"LoginX":       config.C.OAuthXEnabled(),
	}
}

func Home(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		var user *model.User
		if userID != nil {
			db.First(&user, userID)
		}

		c.HTML(http.StatusOK, "home.html", gin.H{
			"User":  user,
			"Title": "Welcome to Scaffold",
		})
	}
}

func LoginForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := gin.H{"Title": "Login"}
		for k, v := range loginFormData() {
			data[k] = v
		}
		if errMsg := c.Query("error"); errMsg != "" {
			data["Error"] = "Login failed. Please try again."
		}
		c.HTML(http.StatusOK, "login.html", data)
	}
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.C.Login.PasswordEnabled {
			c.Redirect(http.StatusFound, "/login?error=password_disabled")
			return
		}
		username := c.PostForm("username")
		password := c.PostForm("password")

		renderLoginError := func(msg string) {
			data := gin.H{"Title": "Login", "Error": msg}
			for k, v := range loginFormData() {
				data[k] = v
			}
			c.HTML(http.StatusOK, "login.html", data)
		}

		var user model.User
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			renderLoginError("Invalid username or password")
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			renderLoginError("Invalid username or password")
			return
		}

		session := sessions.Default(c)
		session.Set("user_id", user.Model.ID)
		session.Save()

		c.Redirect(http.StatusFound, "/")
	}
}

func RegisterForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := gin.H{"Title": "Register"}
		for k, v := range loginFormData() {
			data[k] = v
		}
		c.HTML(http.StatusOK, "register.html", data)
	}
}

func Register(db *gorm.DB, emailService *utils.EmailService) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		email := c.PostForm("email")
		password := c.PostForm("password")
		name := c.PostForm("name")

		registerData := func(errMsg string) gin.H {
			data := gin.H{"Title": "Register"}
			if errMsg != "" {
				data["Error"] = errMsg
			}
			for k, v := range loginFormData() {
				data[k] = v
			}
			return data
		}

		// Check if user already exists
		var existingUser model.User
		if err := db.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
			c.HTML(http.StatusOK, "register.html", registerData("Username or email already exists"))
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.HTML(http.StatusOK, "register.html", registerData("Error creating account"))
			return
		}

		user := model.User{
			Username: username,
			Email:    email,
			Password: string(hashedPassword),
			Name:     name,
		}

		if err := db.Create(&user).Error; err != nil {
			c.HTML(http.StatusOK, "register.html", registerData("Error creating account"))
			return
		}

		// Send welcome email if Resend is configured (non-blocking)
		if emailService != nil {
			go func() {
				_ = emailService.SendWelcome(user.Email, user.Name)
			}()
		}

		// Auto-login after registration
		session := sessions.Default(c)
		session.Set("user_id", user.Model.ID)
		session.Save()

		c.Redirect(http.StatusFound, "/")
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(http.StatusFound, "/")
	}
}

func Profile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (set by auth middleware)
		user, exists := c.Get("user")
		if !exists {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		userModel, ok := user.(model.User)
		if !ok {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		c.HTML(http.StatusOK, "profile.html", gin.H{
			"User":  userModel,
			"Title": "Profile",
		})
	}
}
