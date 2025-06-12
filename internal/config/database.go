package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"hospital-project/internal/models"
)

// Database configuration
type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	URL      string
}

// NewDatabaseConfig creates a new database configuration from environment variables
func NewDatabaseConfig() *Database {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, using environment variables")
	}

	return &Database{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		Name:     getEnv("DB_NAME", "hospital"),
		URL:      getEnv("DB_URL", ""),
	}
}

// Connect establishes a connection to the database
func (db *Database) Connect() (*gorm.DB, error) {
	var dsn string

	if db.URL != "" {
		dsn = db.URL // use full connection string from DB_URL
	} else {
		// Fallback to individual parts (local dev)
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			db.Host, db.Port, db.User, db.Password, db.Name,
		)
	}

	// Configure logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Connect to database
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return gormDB, nil
}

// MigrateDB performs database migrations
func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Patient{},
	)
}

// Helper function to get environment variable with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
