package db

import (
	"fmt"
	"os"

	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func InitDB(logger *zap.Logger) {
	logger.Info("Initializing database connection")

	dbHost := os.Getenv("DATABASE_HOST")
	dbUser := os.Getenv("DATABASE_USER")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")
	dbPort := os.Getenv("DATABASE_PORT")

	logger.Debug("Database configuration",
		zap.String("DATABASE_HOST", dbHost),
		zap.String("DATABASE_USER", dbUser),
		zap.String("DATABASE_NAME", dbName),
		zap.String("DATABASE_PORT", dbPort))

	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" {
		logger.Fatal("Missing database environment variables. Make sure .env file exists and is properly configured")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	logger.Debug("Attempting to connect to database")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return
	}

	logger.Info("Successfully connected to database")

	logger.Info("Running auto-migration")
	err = db.AutoMigrate(&models.User{}, &models.Customer{}, &models.Vehicle{}, &models.Input{})
	if err != nil {
		logger.Error("Failed to run auto-migration", zap.Error(err))
		return
	}

	logger.Info("Auto-migration completed successfully")

	DB = db
}
