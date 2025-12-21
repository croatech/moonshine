package domain

type EquipmentCategory struct {
	Model
	Name           string           `json:"name" db:"name"`
	Type           string           `json:"type" db:"type"`
	EquipmentItems []*EquipmentItem `json:"equipment_items,omitempty"`
}
