package models

import "github.com/jinzhu/gorm"

type Bot struct {
	gorm.Model
	Name         string         `json:"name"`
	Attack       uint           `json:"attack"`
	Defense      uint           `json:"defense"`
	Hp           uint           `json:"hp"`
	Inventory    []int          `json:"inventory" gorm:"type:int[]; default: []"`
	Level        uint           `json:"level"`
	Avatar       string         `json:"avatar"`
	LocationBots []*LocationBot `json:"location_bots"`
	Fights       []*Fight       `json:"fights"`
}
