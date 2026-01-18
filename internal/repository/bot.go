package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"moonshine/internal/domain"
)

type BotRepository struct {
	db *sqlx.DB
}

func NewBotRepository(db *sqlx.DB) *BotRepository {
	return &BotRepository{db: db}
}

func (r *BotRepository) Create(bot *domain.Bot) error {
	query := `
		INSERT INTO bots (id, name, attack, defense, hp, level, avatar)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	if bot.ID == uuid.Nil {
		bot.ID = uuid.New()
	}

	_, err := r.db.Exec(query,
		bot.ID, bot.Name, bot.Attack, bot.Defense, bot.Hp, bot.Level, bot.Avatar,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *BotRepository) FindBotsByLocationID(locationID uuid.UUID) ([]*domain.Bot, error) {
	query := `
		SELECT id, name, level
		FROM bots b
		INNER JOIN location_bots lb ON lb.bot_id = b.id
		WHERE lb.location_id = $1 AND lb.deleted_at IS NULL
	`

	var bots []*domain.Bot
	err := r.db.Select(&bots, query, locationID)
	if err != nil {
		return nil, err
	}

	return bots, nil
}
