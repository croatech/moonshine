package models

import (
	"github.com/jinzhu/gorm"
)

type Message struct {
	gorm.Model
	UserID      int    `json:"user_id"`
	User        *User  `json:"user"`
	Text        string `json:"text"`
	RecipientID uint   `json:"recipient_id"`
}
