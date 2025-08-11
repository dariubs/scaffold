package index

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dariubs/scaffold/app/model"
	"github.com/dariubs/scaffold/app/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UploadProfileImage handles profile image upload
func UploadProfileImage(db *gorm.DB, r2Service *utils.R2Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Get the uploaded file
		file, err := c.FormFile("profile_image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}

		// Validate file type
		allowedTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
		ext := strings.ToLower(filepath.Ext(file.Filename))
		isValidType := false
		for _, allowedType := range allowedTypes {
			if ext == allowedType {
				isValidType = true
				break
			}
		}

		if !isValidType {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only JPG, PNG, GIF, and WebP are allowed"})
			return
		}

		// Validate file size (max 5MB)
		if file.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File too large. Maximum size is 5MB"})
			return
		}

		// Upload to R2
		fileURL, err := r2Service.UploadProfileImage(file, userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
			return
		}

		// Get current user
		var user model.User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		// Delete old profile image if exists
		if user.AvatarURL != "" {
			if err := r2Service.DeleteFile(user.AvatarURL); err != nil {
				// Log error but don't fail the upload
				// You might want to add proper logging here
			}
		}

		// Update user profile with new image URL
		user.AvatarURL = fileURL
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":   "Profile image uploaded successfully",
			"image_url": fileURL,
		})
	}
}

// UploadImage handles general image upload
func UploadImage(db *gorm.DB, r2Service *utils.R2Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Get the uploaded file
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}

		// Get folder from query parameter
		folder := c.Query("folder")
		if folder == "" {
			folder = "general"
		}

		// Validate file type
		allowedTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
		ext := strings.ToLower(filepath.Ext(file.Filename))
		isValidType := false
		for _, allowedType := range allowedTypes {
			if ext == allowedType {
				isValidType = true
				break
			}
		}

		if !isValidType {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only JPG, PNG, GIF, and WebP are allowed"})
			return
		}

		// Validate file size (max 10MB for general images)
		if file.Size > 10*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File too large. Maximum size is 10MB"})
			return
		}

		// Upload to R2
		fileURL, err := r2Service.UploadImage(file, folder)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":   "Image uploaded successfully",
			"image_url": fileURL,
			"folder":    folder,
		})
	}
}

// DeleteImage handles image deletion
func DeleteImage(db *gorm.DB, r2Service *utils.R2Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		imageURL := c.PostForm("image_url")
		if imageURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image URL is required"})
			return
		}

		// Delete from R2
		if err := r2Service.DeleteFile(imageURL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Image deleted successfully",
		})
	}
}
