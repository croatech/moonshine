package domain

import "github.com/google/uuid"

type LocationBot struct {
	Model
	LocationID uuid.UUID `json:"location_id"`
	Location   *Location `json:"location,omitempty"`
	BotID      uuid.UUID `json:"bot_id"`
	Bot        *Bot      `json:"bot,omitempty"`
}
