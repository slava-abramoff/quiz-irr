package models

type TestResult struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	TestID string `gorm:"type:uuid;not null;index"`
	Test   Test   `gorm:"foreignKey:TestID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Данные участника
	FullName string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"type:varchar(255);not null"`
	Org      string `gorm:"type:varchar(255)"`

	// Метрики результата
	Duration   int `gorm:"not null"`
	TotalScore int `gorm:"not null;index"`
}
