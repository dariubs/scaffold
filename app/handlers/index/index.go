package index

import (
	"net/http"

	"github.com/dariubs/scaffold/app/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

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
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Title": "Login",
		})
	}
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		var user model.User
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"Title": "Login",
				"Error": "Invalid username or password",
			})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"Title": "Login",
				"Error": "Invalid username or password",
			})
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
		c.HTML(http.StatusOK, "register.html", gin.H{
			"Title": "Register",
		})
	}
}

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		email := c.PostForm("email")
		password := c.PostForm("password")
		name := c.PostForm("name")

		// Check if user already exists
		var existingUser model.User
		if err := db.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
			c.HTML(http.StatusOK, "register.html", gin.H{
				"Title": "Register",
				"Error": "Username or email already exists",
			})
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.HTML(http.StatusOK, "register.html", gin.H{
				"Title": "Register",
				"Error": "Error creating account",
			})
			return
		}

		user := model.User{
			Username: username,
			Email:    email,
			Password: string(hashedPassword),
			Name:     name,
		}

		if err := db.Create(&user).Error; err != nil {
			c.HTML(http.StatusOK, "register.html", gin.H{
				"Title": "Register",
				"Error": "Error creating account",
			})
			return
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
