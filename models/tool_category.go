package models

import "github.com/jinzhu/gorm"

type ToolCategory struct {
	gorm.Model
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	ToolItems []*ToolItem `json:"tool_items"`
}
