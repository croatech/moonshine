package domain

type Movement struct {
	Model
	UserID uint  `json:"user_id"`
	User   *User `json:"user,omitempty"`
	Status uint  `json:"status"`
}
