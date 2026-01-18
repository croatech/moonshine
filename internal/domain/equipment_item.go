package domain

import "github.com/google/uuid"

type EquipmentItem struct {
	Model
	Name                string    `db:"name"`
	Slug                string    `db:"slug"`
	Attack              uint      `db:"attack"`
	Defense             uint      `db:"defense"`
	Hp                  uint      `db:"hp"`
	RequiredLevel       uint      `db:"required_level"`
	Price               uint      `db:"price"`
	Artifact            bool      `db:"artifact"`
	EquipmentCategoryID uuid.UUID `db:"equipment_category_id"`
	Image               string    `db:"image"`
}
