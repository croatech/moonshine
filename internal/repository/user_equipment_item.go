package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"moonshine/internal/domain"
)

var (
	ErrUserEquipmentItemNotFound = errors.New("user equipment item not found")
)

type dbInterface interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Select(dest interface{}, query string, args ...interface{}) error
}

type UserEquipmentItemRepository struct {
	db dbInterface
}

func NewUserEquipmentItemRepository(db dbInterface) *UserEquipmentItemRepository {
	return &UserEquipmentItemRepository{db: db}
}

func (r *UserEquipmentItemRepository) Create(userEquipmentItem *domain.UserEquipmentItem) error {
	query := `
		INSERT INTO user_equipment_items (id, user_id, equipment_item_id)
		VALUES ($1, $2, $3)
	`

	if userEquipmentItem.ID == uuid.Nil {
		userEquipmentItem.ID = uuid.New()
	}

	_, err := r.db.Exec(query,
		userEquipmentItem.ID,
		userEquipmentItem.UserID,
		userEquipmentItem.EquipmentItemID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserEquipmentItemRepository) FindByUserID(userID uuid.UUID) ([]*domain.EquipmentItem, error) {
	query := `
		SELECT 
			ei.id, ei.created_at, ei.deleted_at, ei.name, ei.slug, ei.attack, ei.defense, ei.hp,
			ei.required_level, ei.price, ei.artifact, ei.equipment_category_id, ei.image
		FROM user_equipment_items uei
		INNER JOIN equipment_items ei ON uei.equipment_item_id = ei.id
		WHERE uei.user_id = $1 
			AND uei.deleted_at IS NULL
			AND ei.deleted_at IS NULL
		ORDER BY ei.name ASC
	`

	log.Printf("[UserEquipmentItemRepository] Querying inventory for user: %s", userID)
	var items []*domain.EquipmentItem
	
	err := r.db.Select(&items, query, userID)
	if err != nil {
		log.Printf("[UserEquipmentItemRepository] Error querying inventory for user %s: %+v", userID, err)
		return nil, err
	}

	log.Printf("[UserEquipmentItemRepository] Found %d items in inventory for user: %s", len(items), userID)
	return items, nil
}

