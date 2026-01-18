package domain

import (
	"time"

	"github.com/google/uuid"
)

type MovementStatus string

const (
	MovementStatusActive   MovementStatus = "active"
	MovementStatusFinished MovementStatus = "finished"
)

type Movement struct {
	Model
	UserID uuid.UUID      `json:"user_id" db:"user_id"`
	User   *User          `json:"user,omitempty" db:"-"`
	Status MovementStatus `json:"status" db:"status"`
}

type MovementCell struct {
	ID         uuid.UUID `json:"id" db:"id"`
	MovementID uuid.UUID `json:"movement_id" db:"movement_id"`
	FromCellID uuid.UUID `json:"from_cell_id" db:"from_cell_id"`
	ToCellID   uuid.UUID `json:"to_cell_id" db:"to_cell_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
