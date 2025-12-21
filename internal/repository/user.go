package repository

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"moonshine/internal/domain"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

type UserRepository struct {
	db           *sqlx.DB
	locationRepo *LocationRepository
	avatarRepo   *AvatarRepository
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db:           db,
		locationRepo: NewLocationRepository(db),
		avatarRepo:   NewAvatarRepository(db),
	}
}

func (r *UserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (
			id, username, email, password, name, avatar_id, location_id,
			current_hp, exp, free_stats, gold, hp, level,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)
	`

	now := time.Now()
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = now
	}

	// Set default location if not provided (moonshine)
	if user.LocationID == uuid.Nil {
		location, err := r.locationRepo.FindStartLocation()
		if err != nil {
			return errors.New("default location 'moonshine' not found")
		}
		user.LocationID = location.ID
	}

	// Set default avatar if not provided (first available avatar)
	if user.AvatarID == nil {
		avatar, err := r.avatarRepo.FindFirst()
		if err == nil {
			user.AvatarID = &avatar.ID
		}
		// Avatar is optional, so we don't return error if not found
	}

	_, err := r.db.Exec(query,
		user.ID, user.Username, user.Email, user.Password, user.Name, user.AvatarID, user.LocationID,
		user.CurrentHp, user.Exp, user.FreeStats, user.Gold, user.Hp, user.Level,
		user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		// Check for unique constraint violation
		if isUniqueConstraintError(err) {
			return ErrUserExists
		}
		return err
	}
	return nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, created_at, updated_at, deleted_at, username, email, password, name, 
			avatar_id, location_id, attack, defense, current_hp, exp, fishing_skill, fishing_slot,
			free_stats, gold, hp, level, lumberjacking_skill, lumberjacking_slot,
			chest_equipment_item_id, belt_equipment_item_id, head_equipment_item_id,
			neck_equipment_item_id, weapon_equipment_item_id, shield_equipment_item_id,
			legs_equipment_item_id, feet_equipment_item_id, arms_equipment_item_id,
			hands_equipment_item_id, ring1_equipment_item_id, ring2_equipment_item_id,
			ring3_equipment_item_id, ring4_equipment_item_id
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	user := &domain.User{}
	err := r.db.Get(user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByUsername(username string) (*domain.User, error) {
	query := `
		SELECT id, created_at, updated_at, deleted_at, username, email, password, name, 
			avatar_id, location_id, attack, defense, current_hp, exp, fishing_skill, fishing_slot,
			free_stats, gold, hp, level, lumberjacking_skill, lumberjacking_slot,
			chest_equipment_item_id, belt_equipment_item_id, head_equipment_item_id,
			neck_equipment_item_id, weapon_equipment_item_id, shield_equipment_item_id,
			legs_equipment_item_id, feet_equipment_item_id, arms_equipment_item_id,
			hands_equipment_item_id, ring1_equipment_item_id, ring2_equipment_item_id,
			ring3_equipment_item_id, ring4_equipment_item_id
		FROM users
		WHERE username = $1 AND deleted_at IS NULL
	`

	user := &domain.User{}
	err := r.db.Get(user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	// PostgreSQL unique constraint violation error code is 23505
	errStr := err.Error()
	return strings.Contains(errStr, "23505") ||
		strings.Contains(errStr, "unique constraint") ||
		strings.Contains(errStr, "duplicate key")
}
