package domain

type UserEquipmentItem struct {
	Model
	UserID          uint           `json:"user_id"`
	User            *User          `json:"user,omitempty"`
	EquipmentItemID uint           `json:"equipment_item_id"`
	EquipmentItem   *EquipmentItem `json:"equipment_item,omitempty"`
}
