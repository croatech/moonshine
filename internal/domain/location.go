package domain

type Location struct {
	Model
	Name              string              `json:"name" db:"name"`
	Slug              string              `json:"slug" db:"slug"`
	Cell              bool                `json:"cell" db:"cell"`
	Inactive          bool                `json:"inactive" db:"inactive"`
	Image             string              `json:"image" db:"image"`
	ImageBg           string              `json:"image_bg" db:"image_bg"`
	LocationBots      []*LocationBot      `json:"location_bots,omitempty"`
	Users             []*User             `json:"users,omitempty"`
	LocationLocations []*LocationLocation `json:"location_locations,omitempty"`
	LocationResources []*LocationResource `json:"location_resources,omitempty"`
}
