package models

import "github.com/jinzhu/gorm"

type Avatar struct {
	gorm.Model
	Private bool    `json:"private" sql:"DEFAULT:true"`
	Image   string  `json:"image"`
	Users   []*User `json:"users"`
}
