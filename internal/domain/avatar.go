package domain

type Avatar struct {
	Model
	Private bool    `json:"private" gorm:"default:true"`
	Image   string  `json:"image"`
	Users   []*User `json:"users,omitempty"`
}
