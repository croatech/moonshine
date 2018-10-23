package models

import "github.com/jinzhu/gorm"

type Location struct {
	gorm.Model
	Name              string              `json:"name"`
	Slug              string              `json:"slug"`
	Cell              bool                `json:"cell" sql:"DEFAULT:true"`
	Inactive          bool                `json:"inactive" sql:"DEFAULT:true"`
	ParentID          uint                `json:"parent_id"`
	Parent            *Location           `json:"parent"`
	LocationBots      []*LocationBot      `json:"location_bots"`
	Users             []*User             `json:"users"`
	LocationLocations []*LocationLocation `json:"location_locations"`
	LocationResources []*LocationResource `json:"location_resources"`
}
