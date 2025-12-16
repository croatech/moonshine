package domain

import "github.com/jinzhu/gorm"

type BotEquipmentItem struct {
	gorm.Model
	BotID           uint           `json:"bot_id"`
	Bot             *Bot           `json:"bot"`
	EquipmentItemID uint           `json:"equipment_item_id"`
	EquipmentItem   *EquipmentItem `json:"equipment_item"`
}
