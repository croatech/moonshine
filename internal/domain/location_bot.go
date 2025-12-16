package domain

import "github.com/jinzhu/gorm"

type LocationBot struct {
	gorm.Model
	LocationID uint      `json:"location_id"`
	Location   *Location `json:"location"`
	BotID      uint      `json:"bot_id"`
	Bot        *Bot      `json:"bot"`
}
