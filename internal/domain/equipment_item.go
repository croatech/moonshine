package domain

import "github.com/google/uuid"

type EquipmentItem struct {
	Model
	Name                string             `json:"name"`
	Attack              uint               `json:"attack"`
	Defense             uint               `json:"defense"`
	Hp                  uint               `json:"hp"`
	RequiredLevel       uint               `json:"required_level"`
	Price               uint               `json:"price"`
	Artifact            bool               `json:"artifact"`
	EquipmentCategoryID uuid.UUID          `json:"equipment_category_id"`
	EquipmentCategory   *EquipmentCategory `json:"equipment_category,omitempty"`
	Image               string             `json:"image"`
}
