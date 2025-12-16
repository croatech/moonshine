package domain

type LocationLocation struct {
	Model
	LocationID     uint      `json:"location_id"`
	Location       *Location `json:"location,omitempty"`
	NearLocationID uint      `json:"near_location_id"`
}
