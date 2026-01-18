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
	ErrShortestPath     = errors.New("shortest path resolving error")
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

func (r *LocationRepository) FindCellsByLocationID(locationID uuid.UUID) ([]*domain.Location, error) {
	query := `
		SELECT id, created_at, deleted_at, name, slug, cell, inactive, image, image_bg
		FROM locations
		WHERE cell = true 
		AND deleted_at IS NULL
		ORDER BY slug
	`

	var locations []*domain.Location
	err := r.db.Select(&locations, query)
	if err != nil {
		return nil, err
	}

	return locations, nil
}

func (r *LocationRepository) FindAllCells() ([]*domain.Location, error) {
	query := `
		SELECT id, created_at, deleted_at, name, slug, cell, inactive, image, image_bg
		FROM locations
		WHERE cell = true 
		AND deleted_at IS NULL
		ORDER BY slug
	`

	var locations []*domain.Location
	err := r.db.Select(&locations, query)
	if err != nil {
		return nil, err
	}

	return locations, nil
}

func (r *LocationRepository) FindAllConnections() ([]*domain.LocationLocation, error) {
	query := `
		SELECT id, created_at, deleted_at, location_id, near_location_id
		FROM location_locations
		WHERE deleted_at IS NULL
	`

	var connections []*domain.LocationLocation
	err := r.db.Select(&connections, query)
	if err != nil {
		return nil, err
	}

	return connections, nil
}

func (r *LocationRepository) FindAll() ([]*domain.Location, error) {
	query := `
		SELECT id, created_at, deleted_at, name, slug, cell, inactive, image, image_bg
		FROM locations
		WHERE deleted_at IS NULL
	`

	var locations []*domain.Location
	err := r.db.Select(&locations, query)
	if err != nil {
		return nil, err
	}

	return locations, nil
}

func (r *LocationRepository) DefaultOutdoorLocation() (*domain.Location, error) {
	query := `
		SELECT id, created_at, deleted_at, name, slug, cell, inactive, image, image_bg
		FROM locations
		WHERE slug = $1 AND deleted_at IS NULL
	`

	location := &domain.Location{}
	err := r.db.Get(location, query, "29cell")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLocationNotFound
		}
		return nil, err
	}

	return location, nil
}
