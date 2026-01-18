package domain

type EquipmentCategory struct {
	Model
	Name string `db:"name"`
	Type string `db:"type"`
}
