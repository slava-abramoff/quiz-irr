package models

import (
	"time"

	"gorm.io/datatypes"
)

type RawSubmission struct {
	ID     string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TestID string `gorm:"type:uuid;not null;index"`
	Test   Test   `gorm:"foreignKey:TestID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	ParticipantData datatypes.JSON `gorm:"type:jsonb;not null"`
	AnswersPayload  datatypes.JSON `gorm:"type:jsonb;not null"`

	//Статус и время
	Status    string    `gorm:"type:varchar(20);not null;default:'pending';index"`
	StartAt   time.Time `gorm:"not null"`
	EndAt     time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
