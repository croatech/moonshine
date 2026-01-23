package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"moonshine/internal/domain"
)

type FightRepository struct {
	db *sqlx.DB
}

func NewFightRepository(db *sqlx.DB) *FightRepository {

	return &FightRepository{db: db}
}

func (r *FightRepository) Create(fight *domain.Fight) (uuid.UUID, error) {
	query := `
		INSERT INTO fights (user_id, bot_id)
		VALUES ($1, $2)
		RETURNING id
	`

	err := r.db.QueryRow(query,
		fight.UserID, fight.BotID,
	).Scan(&fight.ID)
	if err != nil {
		return fight.ID, err
	}
	return fight.ID, err
}

func (r *FightRepository) FindActiveByUserID(userID uuid.UUID) (*domain.Fight, error) {
	query := `
		SELECT id, created_at, deleted_at, user_id, bot_id, status, dropped_gold, dropped_item_id
		FROM fights
		WHERE user_id = $1 AND status = $2 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`

	fight := &domain.Fight{}
	err := r.db.Get(fight, query, userID, domain.FightStatusInProgress)
	if err != nil {
		return nil, err
	}

	return fight, nil
}
