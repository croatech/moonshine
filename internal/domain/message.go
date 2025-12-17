package domain

import "github.com/google/uuid"

type Message struct {
	Model
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid"`
	User        *User     `json:"user,omitempty"`
	Text        string    `json:"text"`
	RecipientID uuid.UUID `json:"recipient_id" gorm:"type:uuid"`
}
