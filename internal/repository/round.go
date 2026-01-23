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

func (r *RoundRepository) Create(fightID uuid.UUID, userHp uint16, botHp uint16) error {
	query := `
		INSERT INTO rounds (fight_id, player_hp, bot_hp, player_damage, bot_damage)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query, fightID, userHp, botHp, 0, 0)
	return err
}
