package domain

type LocationBot struct {
	Model
	LocationID uint      `json:"location_id"`
	Location   *Location `json:"location,omitempty"`
	BotID      uint      `json:"bot_id"`
	Bot        *Bot      `json:"bot,omitempty"`
}
