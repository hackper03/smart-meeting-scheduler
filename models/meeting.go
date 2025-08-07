package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Meeting struct {
	ID             string    `json:"meeting_id" gorm:"primaryKey"`
	Title          string    `json:"title" gorm:"not null"`
	ParticipantIDs string    `json:"-" gorm:"type:text"` // JSON string of participant IDs
	StartTime      time.Time `json:"start_time" gorm:"not null"`
	EndTime        time.Time `json:"end_time" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type MeetingResponse struct {
	MeetingID      string   `json:"meetingId"`
	Title          string   `json:"title"`
	ParticipantIDs []string `json:"participantIds"`
	StartTime      string   `json:"startTime"`
	EndTime        string   `json:"endTime"`
}

func (m *Meeting) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = "meeting_" + generateID()
	}
	return nil
}

func (m *Meeting) SetParticipantIDs(ids []string) error {
	data, err := json.Marshal(ids)
	if err != nil {
		return err
	}
	m.ParticipantIDs = string(data)
	return nil
}

func (m *Meeting) GetParticipantIDs() ([]string, error) {
	var ids []string
	err := json.Unmarshal([]byte(m.ParticipantIDs), &ids)
	return ids, err
}

func generateID() string {
	return time.Now().Format("20060102150405") + randomString(6)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
