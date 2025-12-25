package services

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"moonshine/internal/domain"
	"moonshine/internal/repository"
)

var (
	ErrInsufficientGold      = errors.New("insufficient gold")
	ErrEquipmentItemNotFound = errors.New("equipment item not found")
	ErrAlreadyOwned          = errors.New("item already owned")
)

type EquipmentItemBuyService struct {
	db                    *sqlx.DB
	equipmentItemRepo     *repository.EquipmentItemRepository
	userEquipmentItemRepo *repository.UserEquipmentItemRepository
	userRepo              *repository.UserRepository
}

func NewEquipmentItemBuyService(
	db *sqlx.DB,
	equipmentItemRepo *repository.EquipmentItemRepository,
	userEquipmentItemRepo *repository.UserEquipmentItemRepository,
	userRepo *repository.UserRepository,
) *EquipmentItemBuyService {
	return &EquipmentItemBuyService{
		db:                    db,
		equipmentItemRepo:     equipmentItemRepo,
		userEquipmentItemRepo: userEquipmentItemRepo,
		userRepo:              userRepo,
	}
}

func (s *EquipmentItemBuyService) BuyEquipmentItem(ctx context.Context, userID uuid.UUID, itemSlug string) error {
	// Start transaction
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("[EquipmentItemBuyService] Failed to begin transaction: %+v", err)
		return err
	}
	defer tx.Rollback()

	// Get equipment item by slug
	item, err := s.equipmentItemRepo.FindBySlug(itemSlug)
	if err != nil {
		log.Printf("[EquipmentItemBuyService] Equipment item not found by slug %s: %+v", itemSlug, err)
		return ErrEquipmentItemNotFound
	}

	// Get user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		log.Printf("[EquipmentItemBuyService] User not found: %+v", err)
		return repository.ErrUserNotFound
	}

	// Check if user has enough gold
	if user.Gold < item.Price {
		log.Printf("[EquipmentItemBuyService] Insufficient gold: user has %d, item costs %d", user.Gold, item.Price)
		return ErrInsufficientGold
	}

	// Check if user already owns this item
	checkQuery := `SELECT COUNT(*) FROM user_equipment_items WHERE user_id = $1 AND equipment_item_id = $2 AND deleted_at IS NULL`
	var count int
	err = tx.Get(&count, checkQuery, userID, item.ID)
	if err != nil {
		log.Printf("[EquipmentItemBuyService] Failed to check if item already owned: %+v", err)
		return err
	}
	if count > 0 {
		log.Printf("[EquipmentItemBuyService] Item already owned: user %s, item %s", userID, item.ID)
		return ErrAlreadyOwned
	}

	// Create user equipment item
	userEquipmentItem := &domain.UserEquipmentItem{
		UserID:          userID,
		EquipmentItemID: item.ID,
	}

	// Use transaction-aware repository
	userEquipmentItemRepo := repository.NewUserEquipmentItemRepository(tx)
	if err := userEquipmentItemRepo.Create(userEquipmentItem); err != nil {
		log.Printf("[EquipmentItemBuyService] Failed to create user equipment item: %+v", err)
		return err
	}

	// Update user gold within transaction
	newGold := user.Gold - item.Price
	updateQuery := `UPDATE users SET gold = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err = tx.Exec(updateQuery, newGold, userID)
	if err != nil {
		log.Printf("[EquipmentItemBuyService] Failed to update user gold: %+v", err)
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("[EquipmentItemBuyService] Failed to commit transaction: %+v", err)
		return err
	}

	log.Printf("[EquipmentItemBuyService] Successfully bought item %s (slug: %s) for user %s", item.ID, itemSlug, userID)
	return nil
}
