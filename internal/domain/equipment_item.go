package domain

import "github.com/google/uuid"

type EquipmentItem struct {
	Model
	Name                string             `json:"name"`
	Attack              uint               `json:"attack" gorm:"default:0"`
	Defense             uint               `json:"defense" gorm:"default:0"`
	Hp                  uint               `json:"hp" gorm:"default:0"`
	RequiredLevel       uint               `json:"required_level" gorm:"default:1"`
	Price               uint               `json:"price"`
	Artifact            bool               `json:"artifact" gorm:"default:false"`
	EquipmentCategoryID uuid.UUID          `json:"equipment_category_id" gorm:"type:uuid"`
	EquipmentCategory   *EquipmentCategory `json:"equipment_category,omitempty"`
	Image               string             `json:"image"`
}
