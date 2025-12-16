package domain

import "github.com/jinzhu/gorm"

type Stuff struct {
	gorm.Model
	StuffableType string `json:"stuffable_type"`
	StuffableID   uint   `json:"stuffable_id"`
	UserID        uint   `json:"user_id"`
	User          *User  `json:"user"`
}
