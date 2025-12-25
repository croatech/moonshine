package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"moonshine/internal/domain"
)

var (
	ErrEquipmentItemNotFound = errors.New("equipment item not found")
)

type EquipmentItemRepository struct {
	db *sqlx.DB
}

func NewEquipmentItemRepository(db *sqlx.DB) *EquipmentItemRepository {
	return &EquipmentItemRepository{db: db}
}

func (r *EquipmentItemRepository) FindByCategorySlug(slug string) ([]*domain.EquipmentItem, error) {
	query := `
		SELECT ei.id, ei.created_at, ei.deleted_at, ei.name, ei.slug, ei.attack, ei.defense, ei.hp,
			ei.required_level, ei.price, ei.artifact, ei.equipment_category_id, ei.image
		FROM equipment_items ei
		INNER JOIN equipment_categories ec ON ei.equipment_category_id = ec.id
		WHERE ec.type = $1::equipment_category_type 
			AND ei.deleted_at IS NULL
			AND ec.deleted_at IS NULL
		ORDER BY ei.required_level ASC
	`

	log.Printf("[EquipmentItemRepository] Querying items for category type: %s", slug)
	var items []*domain.EquipmentItem
	err := r.db.Select(&items, query, slug)
	if err != nil {
		log.Printf("[EquipmentItemRepository] Error querying items for category %s: %+v", slug, err)
		return nil, err
	}

	// Log slug for each item to debug
	for i, item := range items {
		log.Printf("[EquipmentItemRepository] Item %d: name=%s, slug=%s", i, item.Name, item.Slug)
	}

	log.Printf("[EquipmentItemRepository] Found %d items for category: %s", len(items), slug)
	return items, nil
}

func (r *EquipmentItemRepository) FindByID(id uuid.UUID) (*domain.EquipmentItem, error) {
	query := `
		SELECT id, created_at, deleted_at, name, slug, attack, defense, hp,
			required_level, price, artifact, equipment_category_id, image
		FROM equipment_items
		WHERE id = $1 AND deleted_at IS NULL
	`

	item := &domain.EquipmentItem{}
	err := r.db.Get(item, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEquipmentItemNotFound
		}
		return nil, err
	}

	return item, nil
}

func (r *EquipmentItemRepository) FindBySlug(slug string) (*domain.EquipmentItem, error) {
	query := `
		SELECT id, created_at, deleted_at, name, slug, attack, defense, hp,
			required_level, price, artifact, equipment_category_id, image
		FROM equipment_items
		WHERE slug = $1 AND deleted_at IS NULL
	`

	item := &domain.EquipmentItem{}
	err := r.db.Get(item, query, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEquipmentItemNotFound
		}
		return nil, err
	}

	return item, nil
}

