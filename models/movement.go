package models

import (
	"github.com/jinzhu/gorm"
)

type Movement struct {
	gorm.Model
	Path   []int `json:"path" gorm:"type:int[]; default: []"`
	UserID uint  `json:"user_id"`
	User   *User `json:"user"`
	Status uint  `json:"status"`
}
