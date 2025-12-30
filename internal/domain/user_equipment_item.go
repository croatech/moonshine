package domain

import "github.com/google/uuid"

type UserEquipmentItem struct {
	Model
	UserID          uuid.UUID      `json:"user_id"`
	User            *User          `json:"user,omitempty"`
	EquipmentItemID uuid.UUID      `json:"equipment_item_id"`
	EquipmentItem   *EquipmentItem `json:"equipment_item,omitempty"`
}
