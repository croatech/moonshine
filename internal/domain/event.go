package domain

type Event struct {
	Model
	UserID uint   `json:"user_id"`
	User   *User  `json:"user,omitempty"`
	Body   string `json:"body"`
}
