package services

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"moonshine/internal/repository"
)

var (
	ErrItemNotInInventory   = errors.New("item not in inventory")
	ErrInsufficientLevel    = errors.New("insufficient level")
	ErrInvalidEquipmentType = errors.New("invalid equipment type")
)

type EquipmentItemTakeOnService struct {
	db                    *sqlx.DB
	equipmentItemRepo     *repository.EquipmentItemRepository
	userEquipmentItemRepo *repository.UserEquipmentItemRepository
	userRepo              *repository.UserRepository
}

func NewEquipmentItemTakeOnService(
	db *sqlx.DB,
	equipmentItemRepo *repository.EquipmentItemRepository,
	userEquipmentItemRepo *repository.UserEquipmentItemRepository,
	userRepo *repository.UserRepository,
) *EquipmentItemTakeOnService {
	return &EquipmentItemTakeOnService{
		db:                    db,
		equipmentItemRepo:     equipmentItemRepo,
		userEquipmentItemRepo: userEquipmentItemRepo,
		userRepo:              userRepo,
	}
}

// Map equipment type to user field name
func getEquipmentFieldName(equipmentType string) (string, error) {
	fieldMap := map[string]string{
		"chest":  "chest_equipment_item_id",
		"belt":   "belt_equipment_item_id",
		"head":   "head_equipment_item_id",
		"neck":   "neck_equipment_item_id",
		"weapon": "weapon_equipment_item_id",
		"shield": "shield_equipment_item_id",
		"legs":   "legs_equipment_item_id",
		"feet":   "feet_equipment_item_id",
		"arms":   "arms_equipment_item_id",
		"hands":  "hands_equipment_item_id",
		// Rings are handled separately
	}

	if fieldName, ok := fieldMap[equipmentType]; ok {
		return fieldName, nil
	}

	if equipmentType == "ring" {
		// For rings, we'll need to find an empty slot
		return "ring", nil
	}

	return "", ErrInvalidEquipmentType
}

func (s *EquipmentItemTakeOnService) TakeOnEquipmentItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) error {
	// Start transaction
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("[EquipmentItemTakeOnService] Failed to begin transaction: %+v", err)
		return err
	}
	defer tx.Rollback()

	// Get user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		log.Printf("[EquipmentItemTakeOnService] User not found: %+v", err)
		return repository.ErrUserNotFound
	}

	// Check if user owns this item in inventory
	checkQuery := `
		SELECT COUNT(*) 
		FROM user_equipment_items 
		WHERE user_id = $1 AND equipment_item_id = $2 AND deleted_at IS NULL
	`
	var count int
	err = tx.Get(&count, checkQuery, userID, itemID)
	if err != nil {
		log.Printf("[EquipmentItemTakeOnService] Failed to check inventory: %+v", err)
		return err
	}
	if count == 0 {
		log.Printf("[EquipmentItemTakeOnService] Item not in inventory: user %s, item %s", userID, itemID)
		return ErrItemNotInInventory
	}

	// Get equipment item with category
	item, err := s.equipmentItemRepo.FindByID(itemID)
	if err != nil {
		log.Printf("[EquipmentItemTakeOnService] Equipment item not found: %+v", err)
		return ErrEquipmentItemNotFound
	}

	// Get category to determine type
	categoryQuery := `SELECT type FROM equipment_categories WHERE id = $1 AND deleted_at IS NULL`
	var equipmentType string
	err = tx.Get(&equipmentType, categoryQuery, item.EquipmentCategoryID)
	if err != nil {
		log.Printf("[EquipmentItemTakeOnService] Failed to get equipment category: %+v", err)
		return err
	}

	// Validate user level
	if user.Level < item.RequiredLevel {
		log.Printf("[EquipmentItemTakeOnService] Insufficient level: user level %d, required %d", user.Level, item.RequiredLevel)
		return ErrInsufficientLevel
	}

	// Get the field name for this equipment type
	fieldName, err := getEquipmentFieldName(equipmentType)
	if err != nil {
		log.Printf("[EquipmentItemTakeOnService] Invalid equipment type: %s", equipmentType)
		return err
	}

	// If there's already an equipped item, return it to inventory
	var oldItemID *uuid.UUID
	var getOldItemQuery string

	if equipmentType == "ring" {
		// For rings, find the first empty slot or use ring1
		getOldItemQuery = `
			SELECT 
				COALESCE(ring1_equipment_item_id, ring2_equipment_item_id, ring3_equipment_item_id, ring4_equipment_item_id) as old_item_id
			FROM users 
			WHERE id = $1 AND deleted_at IS NULL
		`
		err = tx.Get(&oldItemID, getOldItemQuery, userID)
		if err != nil {
			log.Printf("[EquipmentItemTakeOnService] Failed to get old ring: %+v", err)
			// Continue anyway, might be no ring equipped
		}

		// Find first empty ring slot
		if user.Ring1EquipmentItemID == nil {
			fieldName = "ring1_equipment_item_id"
		} else if user.Ring2EquipmentItemID == nil {
			fieldName = "ring2_equipment_item_id"
		} else if user.Ring3EquipmentItemID == nil {
			fieldName = "ring3_equipment_item_id"
		} else if user.Ring4EquipmentItemID == nil {
			fieldName = "ring4_equipment_item_id"
		} else {
			// All slots full, replace ring1
			fieldName = "ring1_equipment_item_id"
			oldItemID = user.Ring1EquipmentItemID
		}
	} else {
		// For other equipment, get the currently equipped item
		getOldItemQuery = fmt.Sprintf(`
			SELECT %s 
			FROM users 
			WHERE id = $1 AND deleted_at IS NULL
		`, fieldName)
		err = tx.Get(&oldItemID, getOldItemQuery, userID)
		if err != nil {
			log.Printf("[EquipmentItemTakeOnService] Failed to get old item (might be empty): %+v", err)
			// Continue anyway, might be no item equipped
		}
	}

	// If there was an old item, return it to inventory
	if oldItemID != nil {
		returnToInventoryQuery := `
			INSERT INTO user_equipment_items (id, user_id, equipment_item_id, created_at)
			VALUES ($1, $2, $3, NOW())
		`
		_, err = tx.Exec(returnToInventoryQuery, uuid.New(), userID, *oldItemID)
		if err != nil {
			log.Printf("[EquipmentItemTakeOnService] Failed to return old item to inventory: %+v", err)
			return err
		}
	}

	// Remove item from inventory
	deleteFromInventoryQuery := `
		DELETE FROM user_equipment_items 
		WHERE user_id = $1 AND equipment_item_id = $2
	`
	_, err = tx.Exec(deleteFromInventoryQuery, userID, itemID)
	if err != nil {
		log.Printf("[EquipmentItemTakeOnService] Failed to remove item from inventory: %+v", err)
		return err
	}

	// Equip the new item
	updateUserQuery := fmt.Sprintf(`
		UPDATE users 
		SET %s = $1 
		WHERE id = $2 AND deleted_at IS NULL
	`, fieldName)
	_, err = tx.Exec(updateUserQuery, itemID, userID)
	if err != nil {
		log.Printf("[EquipmentItemTakeOnService] Failed to equip item: %+v", err)
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("[EquipmentItemTakeOnService] Failed to commit transaction: %+v", err)
		return err
	}

	log.Printf("[EquipmentItemTakeOnService] Successfully equipped item %s for user %s", itemID, userID)
	return nil
}
