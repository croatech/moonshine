package repository

import (
	"database/sql"
	"errors"

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

	var items []*domain.EquipmentItem
	err := r.db.Select(&items, query, slug)
	if err != nil {
		return nil, err
	}

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

func (r *EquipmentItemRepository) Create(item *domain.EquipmentItem) error {
	query := `
		INSERT INTO equipment_items (name, slug, attack, defense, hp, required_level, price, equipment_category_id, image)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	err := r.db.QueryRow(query,
		item.Name, item.Slug, item.Attack, item.Defense, item.Hp,
		item.RequiredLevel, item.Price, item.EquipmentCategoryID, item.Image,
	).Scan(&item.ID)
	if err != nil {
		return err
	}

	return nil
}
