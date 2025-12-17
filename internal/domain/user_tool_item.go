package domain

import "github.com/google/uuid"

type UserToolItem struct {
	Model
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid"`
	User       *User     `json:"user,omitempty"`
	ToolItemID uuid.UUID `json:"tool_item_id" gorm:"type:uuid"`
	ToolItem   *ToolItem `json:"tool_item,omitempty"`
}
