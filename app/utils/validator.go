package utils

import (
	"regexp"
	"strings"
	"unicode"
)

// ValidateEmail validates an email address
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePasswordStrength validates password strength
// Returns true if password meets minimum requirements:
// - At least 8 characters
// - Contains at least one uppercase letter
// - Contains at least one lowercase letter
// - Contains at least one number
func ValidatePasswordStrength(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	_ = hasSpecial // reserved for optional future requirement
	return hasUpper && hasLower && hasNumber
}

// ValidateUsername validates username
// Returns true if username:
// - Is between 3 and 20 characters
// - Contains only alphanumeric characters and underscores
func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}

// SanitizeString removes leading/trailing whitespace and limits length
func SanitizeString(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	if len(s) > maxLen {
		s = s[:maxLen]
	}
	return s
}

