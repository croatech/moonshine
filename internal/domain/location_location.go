package domain

import "github.com/google/uuid"

type LocationLocation struct {
	Model
	LocationID     uuid.UUID `json:"location_id" db:"location_id"`
	Location       *Location `json:"location,omitempty" db:"-"`
	NearLocationID uuid.UUID `json:"near_location_id" db:"near_location_id"`
}
