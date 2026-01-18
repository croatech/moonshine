package domain

import "github.com/google/uuid"

type Resource struct {
	Model
	Name              string              `json:"name"`
	ItemID            uuid.UUID          `json:"item_id"`
	Price             uint                `json:"price"`
	Type              string              `json:"type"`
	Image             string              `json:"image"`
	LocationResources []*LocationResource `json:"location_resources,omitempty"`
}
