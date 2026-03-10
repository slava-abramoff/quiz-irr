package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type RawSubmission struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TestID uuid.UUID `gorm:"type:uuid;not null;index"`
	Test   Test      `gorm:"foreignKey:TestID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	FullName string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"type:varchar(255);not null"`
	Org      string `gorm:"type:varchar(255)"`

	AnswersPayload datatypes.JSON `gorm:"type:jsonb"`

	//Статус и время
	Status    string     `gorm:"type:varchar(20);not null;default:'started';index"`
	StartAt   *time.Time `gorm:"type:timestamptz"`
	EndAt     *time.Time `gorm:"type:timestamptz"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
}
