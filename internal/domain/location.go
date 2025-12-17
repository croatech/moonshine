package domain

type Location struct {
	Model
	Name              string              `json:"name"`
	Slug              string              `json:"slug"`
	Cell              bool                `json:"cell" gorm:"default:true"`
	Inactive          bool                `json:"inactive" gorm:"default:false"`
	Image             string              `json:"image"`
	ImageBg           string              `json:"image_bg"`
	LocationBots      []*LocationBot      `json:"location_bots,omitempty"`
	Users             []*User             `json:"users,omitempty"`
	LocationLocations []*LocationLocation `json:"location_locations,omitempty"`
	LocationResources []*LocationResource `json:"location_resources,omitempty"`
}
