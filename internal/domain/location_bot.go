package domain

import "github.com/google/uuid"

type LocationBot struct {
	Model
	LocationID uuid.UUID `json:"location_id" gorm:"type:uuid"`
	Location   *Location `json:"location,omitempty"`
	BotID      uuid.UUID `json:"bot_id" gorm:"type:uuid"`
	Bot        *Bot      `json:"bot,omitempty"`
}
