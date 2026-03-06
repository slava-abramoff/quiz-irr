package database

import (
	"fmt"
	"quiz-irr/internals/storage/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// cfg *config.Database

func ConnectDB() (*gorm.DB, error) {

	// dsn := fmt.Sprintf(
	// 	"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
	// 	cfg.Host,
	// 	cfg.User,
	// 	cfg.Password,
	// 	cfg.DB,
	// 	cfg.Port,
	// )
	//
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		"localhost",
		"user",
		"password",
		"my_database",
		5432,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
