package domain

import "github.com/google/uuid"

type BotEquipmentItem struct {
	Model
	BotID           uuid.UUID      `json:"bot_id" gorm:"type:uuid"`
	Bot             *Bot           `json:"bot,omitempty"`
	EquipmentItemID uuid.UUID      `json:"equipment_item_id" gorm:"type:uuid"`
	EquipmentItem   *EquipmentItem `json:"equipment_item,omitempty"`
}
