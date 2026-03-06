package models

type Question struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	TestID string `gorm:"type:uuid;not null"`
	Test   Test   `gorm:"foreignKey:TestID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Контент
	Text   string `gorm:"type:text;not null"`
	Type   string `gorm:"type:varchar(50);not null;default:'multiple'"` // Например
	Points int    `gorm:"not null;default:0"`
}
