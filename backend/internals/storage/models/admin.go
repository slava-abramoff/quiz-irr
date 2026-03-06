package models

type Admin struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	FullName string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"type:varchar(255);unique;not null;index"`
	Password string `gorm:"not null"`
	IsRoot   bool   `gorm:"not null;default:false"`
}

// id (UUID/Int, PK)

//    email (String, Unique)

//    password_hash (String)

//    full_name (String)
