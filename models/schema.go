package models

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	Username     string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Email        string `gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time

	Snippets []Snippet `gorm:"foreignKey:UserID"`
}

type Snippet struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	UserID    uint   `gorm:"not null;index"`
	Source    string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
