package domain

import "github.com/google/uuid"

type Stuff struct {
	Model
	StuffableType string    `json:"stuffable_type"`
	StuffableID   uuid.UUID `json:"stuffable_id" gorm:"type:uuid"`
	UserID        uuid.UUID `json:"user_id" gorm:"type:uuid"`
	User          *User  `json:"user,omitempty"`
}
