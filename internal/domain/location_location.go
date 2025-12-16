package domain

import "github.com/jinzhu/gorm"

type LocationLocation struct {
	gorm.Model
	LocationID     uint      `json:"location_id"`
	Location       *Location `json:"location"`
	NearLocationID uint      `json:"near_location_id"`
}
