package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Todos        []Todo
}

type Todo struct {
	gorm.Model
	UserID  uint
	Title   string `gorm:"not null"`
	IsDone  bool   `gorm:"default:false"`
	Duedate *time.Time
}
