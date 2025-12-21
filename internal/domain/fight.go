package domain

import "github.com/google/uuid"

type Fight struct {
	Model
	UserID          uuid.UUID `json:"user_id"`
	User            *User     `json:"user,omitempty"`
	BotID           uuid.UUID `json:"bot_id"`
	Bot             *Bot      `json:"bot,omitempty"`
	Status          uint      `json:"status"`
	WinnerType      string    `json:"winner_type"`
	DroppedGold     uint      `json:"dropped_gold"`
	WinnerID        uuid.UUID `json:"winner_id"`
	DroppedItemID   uuid.UUID `json:"dropped_item_id"`
	DroppedItemType string   `json:"dropped_item_type"`
	Rounds          []*Round `json:"rounds,omitempty"`
}
