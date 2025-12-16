package domain

type Message struct {
	Model
	UserID      uint   `json:"user_id"`
	User        *User  `json:"user,omitempty"`
	Text        string `json:"text"`
	RecipientID uint   `json:"recipient_id"`
}
