package domain

import "github.com/jinzhu/gorm"

type Event struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	User   *User  `json:"user"`
	Body   string `json:"body"`
}
