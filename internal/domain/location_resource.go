package domain

import "github.com/jinzhu/gorm"

type LocationResource struct {
	gorm.Model
	LocationID uint      `json:"location_id"`
	Location   *Location `json:"location"`
	ResourceID uint      `json:"resource_id"`
	Resource   *Resource `json:"resource"`
}
