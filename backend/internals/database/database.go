package database

import (
	"fmt"
	"os"
	"quiz-irr/internals/storage/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// cfg *config.Database

func ConnectDB(envExist bool) (*gorm.DB, error) {
	var dsn string

	if !envExist {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
			"localhost",
			"user",
			"password",
			"my_database",
			5432,
		)
	} else {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_DB"),
			os.Getenv("POSTGRES_PORT"),
		)
	}

	// Default to silent to avoid noisy SQL logs in console.
	// Override with GORM_LOG_LEVEL: silent|error|warn|info
	logLevel := logger.Silent
	switch os.Getenv("GORM_LOG_LEVEL") {
	case "info":
		logLevel = logger.Info
	case "warn":
		logLevel = logger.Warn
	case "error":
		logLevel = logger.Error
	case "silent", "":
		logLevel = logger.Silent
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`).Error; err != nil {
		return nil, fmt.Errorf("failed to enable pgcrypto: %v", err)
	}

	// TODO: Ручные миграции, отдельным скриптом
	if true {
		err = db.AutoMigrate(
			&models.Admin{},
			&models.Test{},
			&models.Question{},
			&models.Option{},
			&models.TestResult{},
			&models.RawSubmission{},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to migrate database: %w", err)
		}
	}
	return db, nil
}
