package domain

type UserToolItem struct {
	Model
	UserID     uint      `json:"user_id"`
	User       *User     `json:"user,omitempty"`
	ToolItemID uint      `json:"tool_item_id"`
	ToolItem   *ToolItem `json:"tool_item,omitempty"`
}
