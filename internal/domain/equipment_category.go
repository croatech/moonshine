package domain

type EquipmentCategory struct {
	Model
	Name           string           `json:"name"`
	Type           string           `json:"type"`
	EquipmentItems []*EquipmentItem `json:"equipment_items,omitempty"`
}
