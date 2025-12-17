package domain

import "github.com/google/uuid"

type Movement struct {
	Model
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid"`
	User   *User `json:"user,omitempty"`
	Status uint  `json:"status"`
}
