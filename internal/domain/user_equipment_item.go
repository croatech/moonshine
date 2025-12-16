package domain

import "github.com/jinzhu/gorm"

type UserEquipmentItem struct {
	gorm.Model
	UserID          uint           `json:"user_id"`
	User            *User          `json:"user"`
	EquipmentItemID uint           `json:"equipment_item_id"`
	EquipmentItem   *EquipmentItem `json:"equipment_item"`
}
