package models

import "github.com/jinzhu/gorm"

type EquipmentCategory struct {
	gorm.Model
	Name           string           `json:"name"`
	Type           string           `json:"type"`
	EquipmentItems []*EquipmentItem `json:"equipment_items"`
}
