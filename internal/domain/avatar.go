package domain

type Avatar struct {
	Model
	Private bool    `json:"private" db:"private"`
	Image   string  `json:"image" db:"image"`
	Users   []*User `json:"users,omitempty"`
}
