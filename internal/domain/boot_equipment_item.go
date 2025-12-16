package domain

type BotEquipmentItem struct {
	Model
	BotID           uint           `json:"bot_id"`
	Bot             *Bot           `json:"bot,omitempty"`
	EquipmentItemID uint           `json:"equipment_item_id"`
	EquipmentItem   *EquipmentItem `json:"equipment_item,omitempty"`
}
