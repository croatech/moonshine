package models

import (
	"github.com/jinzhu/gorm"
)

type Movement struct {
	gorm.Model
	UserID uint  `json:"user_id"`
	User   *User `json:"user"`
	Status uint  `json:"status"`
}
