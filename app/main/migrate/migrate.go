package main

import (
	"log"

	"github.com/dariubs/scaffold/app/database"
	"github.com/dariubs/scaffold/app/model"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database connection
	database.InitDB()

	log.Println("Starting database migrations...")

	// Run migrations
	err = runMigrations()
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("All migrations completed successfully!")
}

func runMigrations() error {
	db := database.DB

	// Migration 1: Create users table
	log.Println("Running migration: Create users table")
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		return err
	}

	// Migration 2: Add any additional indexes or constraints
	log.Println("Running migration: Add additional indexes and constraints")

	// Example: Add a composite index if needed
	// err = db.Exec("CREATE INDEX IF NOT EXISTS idx_users_email_login_method ON users(email, login_method)").Error
	// if err != nil {
	// 	return err
	// }

	// Migration 3: Seed initial data if needed
	log.Println("Running migration: Seed initial data")

	// Example: Create admin user if it doesn't exist
	var adminUser model.User
	result := db.Where("email = ?", "admin@example.com").First(&adminUser)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			log.Println("Creating admin user...")
			adminUser = model.User{
				Username:    "admin",
				Email:       "admin@example.com",
				Password:    "admin123", // In production, use hashed password
				Name:        "Administrator",
				LoginMethod: "password",
			}
			err = db.Create(&adminUser).Error
			if err != nil {
				return err
			}
			log.Println("Admin user created successfully")
		} else {
			return result.Error
		}
	} else {
		log.Println("Admin user already exists")
	}

	return nil
}
