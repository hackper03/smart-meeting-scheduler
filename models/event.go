package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"not null;index"`
	Title     string    `json:"title" gorm:"not null"`
	StartTime time.Time `json:"start_time" gorm:"not null;index"`
	EndTime   time.Time `json:"end_time" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = "event_" + generateID()
	}
	return nil
}
