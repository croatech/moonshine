package domain

import "github.com/google/uuid"

type BotEquipmentItem struct {
	Model
	BotID           uuid.UUID      `json:"bot_id"`
	Bot             *Bot           `json:"bot,omitempty"`
	EquipmentItemID uuid.UUID      `json:"equipment_item_id"`
	EquipmentItem   *EquipmentItem `json:"equipment_item,omitempty"`
}


