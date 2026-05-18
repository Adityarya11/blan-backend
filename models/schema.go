package models

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	Username     string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Email        string `gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash string `gorm:"type:varchar(100);not null"`
	CreatedAt    time.Time

	// foreign keys
	Snippets []Snippet `gorm:"foreignKey:UserID"`
	Jobs     []Job     `gorm:"foreignKey:UserID"`
}

type Snippet struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	UserID    uint   `gorm:"not null;index"`
	Source    string `gorm:"type:text"` // later for strata
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Job struct {
	ID          string `gorm:"primaryKey;type:varchar(36)"`
	UserID      *uint  `gorm:"index"`
	Status      string `gorm:"type:varchar(20)"`
	CreatedAt   time.Time
	CompletedAt *time.Time
}
