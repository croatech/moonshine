package domain

type Fight struct {
	Model
	UserID          uint     `json:"user_id"`
	User            *User    `json:"user,omitempty"`
	BotID           uint     `json:"bot_id"`
	Bot             *Bot     `json:"bot,omitempty"`
	Status          uint     `json:"status" gorm:"default:0"`
	WinnerType      string   `json:"winner_type"`
	DroppedGold     uint     `json:"dropped_gold"`
	WinnerID        uint     `json:"winner_id"`
	DroppedItemID   uint     `json:"dropped_item_id"`
	DroppedItemType string   `json:"dropped_item_type"`
	Rounds          []*Round `json:"rounds,omitempty"`
}
