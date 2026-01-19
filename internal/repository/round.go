package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RoundRepository struct {
	db *sqlx.DB
}

func NewRoundRepository(db *sqlx.DB) *RoundRepository {

	return &RoundRepository{db: db}
}

func (r *RoundRepository) Create(fightID uuid.UUID) error {
	query := `
		INSERT INTO rounds (fight_id)
		VALUES ($1)
	`

	_, err := r.db.Exec(query, fightID)
	return err
}
