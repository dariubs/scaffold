package database

import (
	"log"
	"os"

	"github.com/dariubs/scaffold/app/model"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN not set in .env")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	DB = db

	// Auto migrate the schema - only User model
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	log.Println("Database connected and migrated successfully")
}
