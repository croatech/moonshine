package domain

import (
	"github.com/jinzhu/gorm"
)

type Fight struct {
	gorm.Model
	UserID          uint     `json:"user_id"`
	User            *User    `json:"user"`
	BotID           uint     `json:"bot_id"`
	Bot             *Bot     `json:"bot"`
	Status          uint     `json:"status" sql:"DEFAULT:0"`
	WinnerType      string   `json:"winner_type"`
	DroppedGold     uint     `json:"dropped_gold"`
	WinnerID        uint     `json:"winner_id"`
	DroppedItemID   uint     `json:"dropped_item_id"`
	DroppedItemType string   `json:"dropped_item_type"`
	Rounds          []*Round `json:"rounds"`
}
