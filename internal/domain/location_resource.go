package domain

type LocationResource struct {
	Model
	LocationID uint      `json:"location_id"`
	Location   *Location `json:"location,omitempty"`
	ResourceID uint      `json:"resource_id"`
	Resource   *Resource `json:"resource,omitempty"`
}
