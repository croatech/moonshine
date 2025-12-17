package domain

import "github.com/google/uuid"

type LocationResource struct {
	Model
	LocationID uuid.UUID `json:"location_id" gorm:"type:uuid"`
	Location   *Location `json:"location,omitempty"`
	ResourceID uuid.UUID `json:"resource_id" gorm:"type:uuid"`
	Resource   *Resource `json:"resource,omitempty"`
}
