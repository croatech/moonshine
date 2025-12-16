package domain

type Stuff struct {
	Model
	StuffableType string `json:"stuffable_type"`
	StuffableID   uint   `json:"stuffable_id"`
	UserID        uint   `json:"user_id"`
	User          *User  `json:"user,omitempty"`
}
