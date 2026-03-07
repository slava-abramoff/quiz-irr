package models

import "time"

type Test struct {
	ID       string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AuthorID uint   `gorm:"not null"`            // Внешний ключ (FK)
	Author   Admin  `gorm:"foreignKey:AuthorID"` // Сама структура для связки

	// Ключевые данные теста
	Title    string `gorm:"type:varchar(255);not null"`
	Desc     string `gorm:"type:text"`
	IsActive bool   `gorm:"not null;default:false"`

	// Работа со временем
	StartAt   *time.Time `gorm:"type:timestamptz"`
	EndAt     *time.Time `gorm:"type:timestamptz"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
}

// id (UUID, PK)

// author_id (FK -> Admins.id)

// title (String) — Название викторины.

// description (Text) — Описание или инструкция.

// is_active (Boolean) — Можно ли сейчас проходить этот тест.

// created_at (Timestamp)
