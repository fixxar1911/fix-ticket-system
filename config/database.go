package config

import (
	"fix-ticket-system/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() {
	dbType := getEnv("DB_TYPE", "sqlite")
	var dsn string
	var driver gorm.Dialector
	if dbType == "postgres" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_PORT", "5432"),
			getEnv("DB_USER", "postgres"),
			getEnv("DB_PASSWORD", "postgres"),
			getEnv("DB_NAME", "tickets"),
		)
		driver = postgres.Open(dsn)
	} else {
		dsn = "file::memory:?cache=shared"
		driver = sqlite.Open(dsn)
	}
	db, err := gorm.Open(driver, &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Auto-migrate the schema
	db.AutoMigrate(&models.Ticket{})
	// Set the global DB variable
	DB = db
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
