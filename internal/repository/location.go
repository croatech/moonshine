package repository

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"moonshine/internal/domain"
)

var (
	ErrLocationNotFound = errors.New("location not found")
	ErrLocationExists   = errors.New("location already exists")
)

type LocationRepository struct {
	db *sqlx.DB
}

func NewLocationRepository(db *sqlx.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) Create(location *domain.Location) error {
	query := `
		INSERT INTO locations (id, name, slug, cell, inactive, image, image_bg)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	if location.ID == uuid.Nil {
		location.ID = uuid.New()
	}

	_, err := r.db.Exec(query,
		location.ID, location.Name, location.Slug, location.Cell, location.Inactive,
		location.Image, location.ImageBg,
	)
	if err != nil {
		// Check for unique constraint violation
		if isUniqueConstraintError(err) {
			return ErrLocationExists
		}
		return err
	}
	return nil
}

func (r *LocationRepository) FindByID(id uuid.UUID) (*domain.Location, error) {
	query := `
		SELECT id, created_at, deleted_at, name, slug, cell, inactive, image, image_bg
		FROM locations
		WHERE id = $1 AND deleted_at IS NULL
	`

	location := &domain.Location{}
	err := r.db.Get(location, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLocationNotFound
		}
		return nil, err
	}

	return location, nil
}

func (r *LocationRepository) FindStartLocation() (*domain.Location, error) {
	query := `
		SELECT id, created_at, deleted_at, name, slug, cell, inactive, image, image_bg
		FROM locations
		WHERE slug = $1 AND deleted_at IS NULL
	`

	location := &domain.Location{}
	err := r.db.Get(location, query, "moonshine")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLocationNotFound
		}
		return nil, err
	}

	return location, nil
}

func (r *LocationRepository) FindBySlug(slug string) (*domain.Location, error) {
	query := `
		SELECT id, created_at, deleted_at, name, slug, cell, inactive, image, image_bg
		FROM locations
		WHERE slug = $1 AND deleted_at IS NULL
	`

	location := &domain.Location{}
	err := r.db.Get(location, query, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLocationNotFound
		}
		return nil, err
	}

	return location, nil
}

