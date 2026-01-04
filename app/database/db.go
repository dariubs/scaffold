package database

import (
	"database/sql"
	"log"

	"github.com/dariubs/scaffold/app/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Load configuration
	err := config.Load()
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	db, err := gorm.Open(postgres.Open(config.C.Database.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get database instance:", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	// sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	log.Println("Database connected successfully")
}

// GetDB returns the database instance (for dependency injection pattern)
func GetDB() *gorm.DB {
	return DB
}

// GetSQLDB returns the underlying sql.DB instance for connection pool management
func GetSQLDB() (*sql.DB, error) {
	return DB.DB()
}
