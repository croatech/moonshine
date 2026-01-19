package repository

import (
	"database/sql"
	"errors"
	"strings"

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
			username, email, password, name, avatar_id, location_id,
			attack, defense, current_hp, exp, free_stats, gold, hp, level
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(query,
		user.Username, user.Email, user.Password, user.Name, user.AvatarID, user.LocationID,
		user.Attack, user.Defense, user.CurrentHp, user.Exp, user.FreeStats, user.Gold, user.Hp, user.Level,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if isUniqueConstraintError(err) {
			return ErrUserExists
		}
		return err
	}
	return nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, created_at, deleted_at, username, email, password, name, 
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
		SELECT id, created_at, deleted_at, username, email, password, name, 
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

func (r *UserRepository) UpdateGold(userID uuid.UUID, newGold uint) error {
	query := `UPDATE users SET gold = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, newGold, userID)
	return err
}

func (r *UserRepository) UpdateAvatarID(userID uuid.UUID, avatarID *uuid.UUID) error {
	query := `UPDATE users SET avatar_id = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, avatarID, userID)
	return err
}

func (r *UserRepository) RegenerateAllUsersHealth(percent float64) (int64, error) {
	query := `
		UPDATE users 
		SET current_hp = LEAST(
			current_hp + GREATEST(5, ROUND(hp * $1 / 100.0)), 
			hp
		)
		WHERE deleted_at IS NULL AND current_hp < hp
	`
	result, err := r.db.Exec(query, percent)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "23505") ||
		strings.Contains(errStr, "unique constraint") ||
		strings.Contains(errStr, "duplicate key")
}

func (r *UserRepository) UpdateLocationID(userID uuid.UUID, locationID uuid.UUID) error {
	query := `UPDATE users SET location_id = $1 WHERE id = $2`
	_, err := r.db.Exec(query, locationID, userID)
	return err
}

func (r *UserRepository) InFight(userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM fights WHERE user_id = $1 AND status = $2)`

	exists := false
	err := r.db.Get(&exists, query, userID, domain.FightStatusInProgress)

	return exists, err
}
