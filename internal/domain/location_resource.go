package domain

import "github.com/google/uuid"

type LocationResource struct {
	Model
	LocationID uuid.UUID `json:"location_id"`
	Location   *Location `json:"location,omitempty"`
	ResourceID uuid.UUID `json:"resource_id"`
	Resource   *Resource `json:"resource,omitempty"`
}
