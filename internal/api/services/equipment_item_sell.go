package services

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"moonshine/internal/repository"
)

var (
	ErrItemNotOwned = errors.New("item not owned by user")
)

type EquipmentItemSellService struct {
	db                    *sqlx.DB
	equipmentItemRepo     *repository.EquipmentItemRepository
	userEquipmentItemRepo *repository.UserEquipmentItemRepository
	userRepo              *repository.UserRepository
}

func NewEquipmentItemSellService(
	db *sqlx.DB,
	equipmentItemRepo *repository.EquipmentItemRepository,
	userEquipmentItemRepo *repository.UserEquipmentItemRepository,
	userRepo *repository.UserRepository,
) *EquipmentItemSellService {
	return &EquipmentItemSellService{
		db:                    db,
		equipmentItemRepo:     equipmentItemRepo,
		userEquipmentItemRepo: userEquipmentItemRepo,
		userRepo:              userRepo,
	}
}

func (s *EquipmentItemSellService) SellEquipmentItem(ctx context.Context, userID uuid.UUID, itemSlug string) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("[EquipmentItemSellService] Failed to begin transaction: %+v", err)
		return err
	}
	defer tx.Rollback()

	// Find equipment item by slug
	item, err := s.equipmentItemRepo.FindBySlug(itemSlug)
	if err != nil {
		log.Printf("[EquipmentItemSellService] Equipment item not found by slug %s: %+v", itemSlug, err)
		return ErrEquipmentItemNotFound
	}

	// Check if user owns this item (in inventory, not equipped)
	var count int
	checkOwnershipQuery := `SELECT COUNT(*) FROM user_equipment_items WHERE user_id = $1 AND equipment_item_id = $2 AND deleted_at IS NULL`
	err = tx.Get(&count, checkOwnershipQuery, userID, item.ID)
	if err != nil {
		log.Printf("[EquipmentItemSellService] Failed to check item ownership: %+v", err)
		return err
	}

	if count == 0 {
		log.Printf("[EquipmentItemSellService] User %s does not own item %s in inventory", userID, itemSlug)
		return ErrItemNotOwned
	}

	// Get user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		log.Printf("[EquipmentItemSellService] User not found: %+v", err)
		return repository.ErrUserNotFound
	}

	// Add gold to user
	newGold := user.Gold + item.Price
	updateGoldQuery := `UPDATE users SET gold = $1 WHERE id = $2`
	_, err = tx.Exec(updateGoldQuery, newGold, userID)
	if err != nil {
		log.Printf("[EquipmentItemSellService] Failed to update user gold: %+v", err)
		return err
	}
	log.Printf("[EquipmentItemSellService] Added %d gold to user %s (new balance: %d)", item.Price, userID, newGold)

	// Remove item from user's inventory
	deleteItemQuery := `DELETE FROM user_equipment_items WHERE user_id = $1 AND equipment_item_id = $2 AND deleted_at IS NULL`
	_, err = tx.Exec(deleteItemQuery, userID, item.ID)
	if err != nil {
		log.Printf("[EquipmentItemSellService] Failed to remove item from inventory: %+v", err)
		return err
	}
	log.Printf("[EquipmentItemSellService] Removed item %s from inventory for user %s", item.ID, userID)

	if err := tx.Commit(); err != nil {
		log.Printf("[EquipmentItemSellService] Failed to commit transaction: %+v", err)
		return err
	}

	log.Printf("[EquipmentItemSellService] Successfully sold item %s (slug: %s) by user %s for %d gold", item.ID, itemSlug, userID, item.Price)
	return nil
}

