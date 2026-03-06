package models

type Option struct {
	ID         uint     `gorm:"primaryKey;autoIncrement"`
	QuestionID uint     `gorm:"not null"`
	Question   Question `gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Контент и логика
	Text      string `gorm:"type:text;not null"`
	IsCorrect bool   `gorm:"not null;default:false"`
}

// id (UUID, PK)

// question_id (FK -> Questions.id)

// text (String) — Текст варианта.

// is_correct (Boolean) — Является ли этот вариант верным.
