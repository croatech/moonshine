package models

import "github.com/jinzhu/gorm"

type EquipmentItem struct {
	gorm.Model
	Name                string             `json:"name"`
	Attack              uint               `json:"attack" sql:"DEFAULT:0"`
	Defense             uint               `json:"defense" sql:"DEFAULT:0"`
	Hp                  uint               `json:"hp" sql:"DEFAULT:0"`
	RequiredLevel       uint               `json:"required_level" sql:"DEFAULT:1"`
	Price               uint               `json:"price"`
	Artifact            bool               `json:"artifact" sql:"DEFAULT:false"`
	EquipmentCategoryID uint               `json:"equipment_category_id"`
	EquipmentCategory   *EquipmentCategory `json:"equipment_category"`
	Image               string             `json:"image"`
}
