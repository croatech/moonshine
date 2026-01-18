package repository

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"moonshine/internal/domain"
)

var (
	ErrBotNotFound = errors.New("bot not found")
)

type BotRepository struct {
	db *sqlx.DB
}

func NewBotRepository(db *sqlx.DB) *BotRepository {
	return &BotRepository{db: db}
}

func (r *BotRepository) Create(bot *domain.Bot) error {
	query := `
		INSERT INTO bots (id, name, slug, attack, defense, hp, level, avatar)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	if bot.ID == uuid.Nil {
		bot.ID = uuid.New()
	}

	_, err := r.db.Exec(query,
		bot.ID, bot.Name, bot.Slug, bot.Attack, bot.Defense, bot.Hp, bot.Level, bot.Avatar,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *BotRepository) FindBotsByLocationID(locationID uuid.UUID) ([]*domain.Bot, error) {
	query := `
		SELECT b.id, b.created_at, b.deleted_at, b.name, b.slug, b.attack, b.defense, b.hp, b.level, b.avatar
		FROM bots b
		INNER JOIN location_bots lb ON lb.bot_id = b.id
		WHERE lb.location_id = $1 AND b.deleted_at IS NULL AND lb.deleted_at IS NULL
	`

	var bots []*domain.Bot
	err := r.db.Select(&bots, query, locationID)
	if err != nil {
		return nil, err
	}

	return bots, nil
}

func (r *BotRepository) FindBySlug(slug string) (*domain.Bot, error) {
	query := `
		SELECT id, created_at, deleted_at, name, slug, attack, defense, hp, level, avatar
		FROM bots
		WHERE slug = $1 AND deleted_at IS NULL
	`

	bot := &domain.Bot{}
	err := r.db.Get(bot, query, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBotNotFound
		}
		return nil, err
	}

	return bot, nil
}
