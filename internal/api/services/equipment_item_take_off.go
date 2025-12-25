package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"moonshine/internal/repository"
)

var (
	ErrNoItemEquipped = errors.New("no item equipped in this slot")
)

type EquipmentItemTakeOffService struct {
	db                    *sqlx.DB
	equipmentItemRepo     *repository.EquipmentItemRepository
	userEquipmentItemRepo *repository.UserEquipmentItemRepository
	userRepo              *repository.UserRepository
}

func NewEquipmentItemTakeOffService(
	db *sqlx.DB,
	equipmentItemRepo *repository.EquipmentItemRepository,
	userEquipmentItemRepo *repository.UserEquipmentItemRepository,
	userRepo *repository.UserRepository,
) *EquipmentItemTakeOffService {
	return &EquipmentItemTakeOffService{
		db:                    db,
		equipmentItemRepo:     equipmentItemRepo,
		userEquipmentItemRepo: userEquipmentItemRepo,
		userRepo:              userRepo,
	}
}

// Map slot name to user field name
func getFieldNameFromSlot(slotName string) (string, error) {
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
		"ring1":  "ring1_equipment_item_id",
		"ring2":  "ring2_equipment_item_id",
		"ring3":  "ring3_equipment_item_id",
		"ring4":  "ring4_equipment_item_id",
	}
	
	if fieldName, ok := fieldMap[slotName]; ok {
		return fieldName, nil
	}
	
	return "", ErrInvalidEquipmentType
}

func (s *EquipmentItemTakeOffService) TakeOffEquipmentItem(ctx context.Context, userID uuid.UUID, slotName string) error {
	// Start transaction
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("[EquipmentItemTakeOffService] Failed to begin transaction: %+v", err)
		return err
	}
	defer tx.Rollback()

	// Verify user exists
	_, err = s.userRepo.FindByID(userID)
	if err != nil {
		log.Printf("[EquipmentItemTakeOffService] User not found: %+v", err)
		return repository.ErrUserNotFound
	}

	// Get the field name for this slot
	fieldName, err := getFieldNameFromSlot(slotName)
	if err != nil {
		log.Printf("[EquipmentItemTakeOffService] Invalid slot name: %s", slotName)
		return err
	}

	// Get the currently equipped item ID from this slot
	getItemQuery := fmt.Sprintf(`
		SELECT %s 
		FROM users 
		WHERE id = $1 AND deleted_at IS NULL
	`, fieldName)
	
	var equippedItemIDStr sql.NullString
	err = tx.Get(&equippedItemIDStr, getItemQuery, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[EquipmentItemTakeOffService] User not found")
			return repository.ErrUserNotFound
		}
		log.Printf("[EquipmentItemTakeOffService] Failed to get equipped item: %+v", err)
		return err
	}

	if !equippedItemIDStr.Valid || equippedItemIDStr.String == "" {
		log.Printf("[EquipmentItemTakeOffService] No item equipped in slot %s for user %s", slotName, userID)
		return ErrNoItemEquipped
	}

	equippedItemID, err := uuid.Parse(equippedItemIDStr.String)
	if err != nil {
		log.Printf("[EquipmentItemTakeOffService] Invalid UUID format: %+v", err)
		return err
	}

	// Return item to inventory
	returnToInventoryQuery := `
		INSERT INTO user_equipment_items (id, user_id, equipment_item_id, created_at)
		VALUES ($1, $2, $3, NOW())
	`
	_, err = tx.Exec(returnToInventoryQuery, uuid.New(), userID, equippedItemID)
	if err != nil {
		log.Printf("[EquipmentItemTakeOffService] Failed to return item to inventory: %+v", err)
		return err
	}

	// Clear the slot (set to NULL)
	clearSlotQuery := fmt.Sprintf(`
		UPDATE users 
		SET %s = NULL 
		WHERE id = $1 AND deleted_at IS NULL
	`, fieldName)
	_, err = tx.Exec(clearSlotQuery, userID)
	if err != nil {
		log.Printf("[EquipmentItemTakeOffService] Failed to clear equipment slot: %+v", err)
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("[EquipmentItemTakeOffService] Failed to commit transaction: %+v", err)
		return err
	}

	log.Printf("[EquipmentItemTakeOffService] Successfully removed item %s from slot %s for user %s", equippedItemID, slotName, userID)
	return nil
}

