package domain

import "github.com/jinzhu/gorm"

type Bot struct {
	gorm.Model
	Name         string              `json:"name"`
	Attack       uint                `json:"attack"`
	Defense      uint                `json:"defense"`
	Hp           uint                `json:"hp"`
	Level        uint                `json:"level"`
	Avatar       string              `json:"avatar"`
	LocationBots []*LocationBot      `json:"location_bots"`
	Fights       []*Fight            `json:"fights"`
	Equipment    []*BotEquipmentItem `json:"bot_equipment_items"`
}
