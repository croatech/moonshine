package domain

import "github.com/google/uuid"

type Event struct {
	Model
	UserID uuid.UUID `json:"user_id"`
	User   *User  `json:"user,omitempty"`
	Body   string `json:"body"`
}
