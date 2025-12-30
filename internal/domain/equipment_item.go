package domain

import "github.com/google/uuid"

type EquipmentItem struct {
	Model
	Name                string             `json:"name" db:"name"`
	Slug                string             `json:"slug" db:"slug"`
	Attack              uint               `json:"attack" db:"attack"`
	Defense             uint               `json:"defense" db:"defense"`
	Hp                  uint               `json:"hp" db:"hp"`
	RequiredLevel       uint               `json:"required_level" db:"required_level"`
	Price               uint               `json:"price" db:"price"`
	Artifact            bool               `json:"artifact" db:"artifact"`
	EquipmentCategoryID uuid.UUID          `json:"equipment_category_id" db:"equipment_category_id"`
	EquipmentCategory   *EquipmentCategory `json:"equipment_category,omitempty"`
	Image               string             `json:"image" db:"image"`
}
