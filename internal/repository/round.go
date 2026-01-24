package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"moonshine/internal/domain"
)

type RoundRepository struct {
	db *sqlx.DB
}

func NewRoundRepository(db *sqlx.DB) *RoundRepository {

	return &RoundRepository{db: db}
}

func (r *RoundRepository) Create(fightID uuid.UUID, userHp, botHp uint) error {
	query := `
		INSERT INTO rounds (fight_id, player_hp, bot_hp, player_damage, bot_damage, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(query, fightID, userHp, botHp, 0, 0, domain.RoundStatusInProgress)
	return err
}

func (r *RoundRepository) FindByFightID(fightID uuid.UUID) ([]*domain.Round, error) {
	query := `
		SELECT id, created_at, deleted_at, fight_id, player_damage, bot_damage, 
			status, player_hp, bot_hp, player_attack_point, player_defense_point, 
			bot_attack_point, bot_defense_point
		FROM rounds 
		WHERE fight_id = $1 AND deleted_at IS NULL 
		ORDER BY created_at DESC
	`

	var rounds []*domain.Round
	err := r.db.Select(&rounds, query, fightID)
	if err != nil {
		return nil, err
	}

	return rounds, nil
}

func (r *RoundRepository) FinishRound(id uuid.UUID, botAttackPoint, botDefensePoint, playerAttackPoint, playerDefensePoint string,
	playerDmg, botDmg, finalPlayerHp, finalBotHp uint) error {
	query := `
		UPDATE rounds 
		SET bot_attack_point = $1,
		    bot_defense_point = $2,
		    player_attack_point = $3,
		    player_defense_point = $4,
		    player_damage = $5,
		    bot_damage = $6,
		    player_hp = $7,
		    bot_hp = $8,
		    status = $9
		WHERE id = $10
	`

	_, err := r.db.Exec(query, botAttackPoint, botDefensePoint, playerAttackPoint, playerDefensePoint, playerDmg, botDmg, finalPlayerHp, finalBotHp, domain.RoundStatusFinished, id)
	if err != nil {
		return err
	}

	return nil
}
