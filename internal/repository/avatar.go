package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"moonshine/internal/domain"
)

var (
	ErrAvatarNotFound = errors.New("avatar not found")
	ErrAvatarExists   = errors.New("avatar already exists")
)

type AvatarRepository struct {
	db *sqlx.DB
}

func NewAvatarRepository(db *sqlx.DB) *AvatarRepository {
	return &AvatarRepository{db: db}
}

func (r *AvatarRepository) Create(avatar *domain.Avatar) error {
	query := `
		INSERT INTO avatars (id, image, private, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	now := time.Now()
	if avatar.ID == uuid.Nil {
		avatar.ID = uuid.New()
	}
	if avatar.CreatedAt.IsZero() {
		avatar.CreatedAt = now
	}
	if avatar.UpdatedAt.IsZero() {
		avatar.UpdatedAt = now
	}

	_, err := r.db.Exec(query,
		avatar.ID, avatar.Image, avatar.Private,
		avatar.CreatedAt, avatar.UpdatedAt,
	)
	if err != nil {
		// Check for unique constraint violation
		if isUniqueConstraintError(err) {
			return ErrAvatarExists
		}
		return err
	}
	return nil
}

func (r *AvatarRepository) FindByID(id uuid.UUID) (*domain.Avatar, error) {
	query := `
		SELECT id, created_at, updated_at, deleted_at, image, private
		FROM avatars
		WHERE id = $1 AND deleted_at IS NULL
	`

	avatar := &domain.Avatar{}
	err := r.db.Get(avatar, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAvatarNotFound
		}
		return nil, err
	}

	return avatar, nil
}

func (r *AvatarRepository) FindByImage(image string) (*domain.Avatar, error) {
	query := `
		SELECT id, created_at, updated_at, deleted_at, image, private
		FROM avatars
		WHERE image = $1 AND deleted_at IS NULL
	`

	avatar := &domain.Avatar{}
	err := r.db.Get(avatar, query, image)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAvatarNotFound
		}
		return nil, err
	}

	return avatar, nil
}

func (r *AvatarRepository) FindFirst() (*domain.Avatar, error) {
	query := `
		SELECT id, created_at, updated_at, deleted_at, image, private
		FROM avatars
		WHERE deleted_at IS NULL
		LIMIT 1
	`

	avatar := &domain.Avatar{}
	err := r.db.Get(avatar, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAvatarNotFound
		}
		return nil, err
	}

	return avatar, nil
}
