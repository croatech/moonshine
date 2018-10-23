package models

import "github.com/jinzhu/gorm"

type ToolItem struct {
	gorm.Model
	Name           string        `json:"name"`
	Price          uint          `json:"price"`
	RequiredSkill  uint          `json:"required_skill"`
	ToolCategoryID uint          `json:"tool_category_id"`
	ToolCategory   *ToolCategory `json:"tool_category"`
	Image          string        `json:"image"`
}
