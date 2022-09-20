package models

import "github.com/jinzhu/gorm"

type UserToolItem struct {
	gorm.Model
	UserID     uint      `json:"user_id"`
	User       *User     `json:"user"`
	ToolItemID uint      `json:"tool_item_id"`
	ToolItem   *ToolItem `json:"tool_item"`
}
