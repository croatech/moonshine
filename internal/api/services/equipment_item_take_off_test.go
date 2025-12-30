package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"moonshine/internal/domain"
	"moonshine/internal/repository"
)


func setupTestDataForTakeOff(db *sqlx.DB) (*domain.User, *domain.EquipmentItem, uuid.UUID, error) {
	locationID := uuid.New()
	locationQuery := `INSERT INTO locations (id, name, slug, cell, inactive) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(locationQuery, locationID, "Test Location", "test_location", false, false)
	if err != nil {
		return nil, nil, uuid.Nil, fmt.Errorf("failed to create location: %w", err)
	}

	categoryID := uuid.New()
	categoryQuery := `INSERT INTO equipment_categories (id, name, type) VALUES ($1, $2, $3::equipment_category_type)`
	_, err = db.Exec(categoryQuery, categoryID, "Weapon", "weapon")
	if err != nil {
		return nil, nil, uuid.Nil, fmt.Errorf("failed to create category: %w", err)
	}

	itemID := uuid.New()
	itemQuery := `INSERT INTO equipment_items (id, name, slug, attack, defense, hp, required_level, price, equipment_category_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err = db.Exec(itemQuery, itemID, "Test Sword", "test-sword", 10, 5, 20, 1, 100, categoryID)
	if err != nil {
		return nil, nil, uuid.Nil, fmt.Errorf("failed to create item: %w", err)
	}

	userID := uuid.New()
	userQuery := `INSERT INTO users (id, username, email, password, location_id, attack, defense, hp, current_hp, level, weapon_equipment_item_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	ts := time.Now().UnixNano()
	username := fmt.Sprintf("testuser%d", ts)
	_, err = db.Exec(userQuery, userID, username, fmt.Sprintf("test%d@example.com", ts), "password", locationID, 11, 6, 40, 40, 5, itemID)
	if err != nil {
		return nil, nil, uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	user := &domain.User{
		Model: domain.Model{
			ID: userID,
		},
		Username: username,
		Level:    5,
		Attack:   11,
		Defense:  6,
		Hp:       40,
	}

	item := &domain.EquipmentItem{
		Model: domain.Model{
			ID: itemID,
		},
		Name:              "Test Sword",
		Slug:              "test-sword",
		Attack:            10,
		Defense:           5,
		Hp:                20,
		RequiredLevel:     1,
		Price:             100,
		EquipmentCategoryID: categoryID,
	}

	return user, item, categoryID, nil
}

func TestEquipmentItemTakeOffService_TakeOffEquipmentItem(t *testing.T) {
	if testDB == nil {
		t.Skip("Test database not initialized")
	}

	db := testDB.DB()
	ctx := context.Background()

	user, item, _, err := setupTestDataForTakeOff(db)
	if err != nil {
		t.Fatalf("Failed to setup test data: %v", err)
	}

	equipmentItemRepo := repository.NewEquipmentItemRepository(db)
	userEquipmentItemRepo := repository.NewUserEquipmentItemRepository(db)
	userRepo := repository.NewUserRepository(db)
	service := NewEquipmentItemTakeOffService(db, equipmentItemRepo, userEquipmentItemRepo, userRepo)

	t.Run("successfully unequip item", func(t *testing.T) {
		err := service.TakeOffEquipmentItem(ctx, user.ID, "weapon")
		if err != nil {
			t.Fatalf("Failed to unequip item: %v", err)
		}

		var equippedItemID *uuid.UUID
		query := `SELECT weapon_equipment_item_id FROM users WHERE id = $1`
		err = db.Get(&equippedItemID, query, user.ID)
		if err != nil {
			t.Fatalf("Failed to get equipped item: %v", err)
		}

		if equippedItemID != nil {
			t.Errorf("Expected no equipped item, got %s", *equippedItemID)
		}

		type stats struct {
			Attack  uint `db:"attack"`
			Defense uint `db:"defense"`
			Hp      uint `db:"hp"`
		}
		var userStats stats
		statsQuery := `SELECT attack, defense, hp FROM users WHERE id = $1`
		err = db.Get(&userStats, statsQuery, user.ID)
		if err != nil {
			t.Fatalf("Failed to get user stats: %v", err)
		}

		expectedAttack := uint(1)
		expectedDefense := uint(1)
		expectedHp := uint(20)

		if userStats.Attack != expectedAttack {
			t.Errorf("Expected attack %d, got %d", expectedAttack, userStats.Attack)
		}
		if userStats.Defense != expectedDefense {
			t.Errorf("Expected defense %d, got %d", expectedDefense, userStats.Defense)
		}
		if userStats.Hp != expectedHp {
			t.Errorf("Expected hp %d, got %d", expectedHp, userStats.Hp)
		}

		var inventoryCount int
		inventoryCountQuery := `SELECT COUNT(*) FROM user_equipment_items WHERE user_id = $1 AND equipment_item_id = $2`
		err = db.Get(&inventoryCount, inventoryCountQuery, user.ID, item.ID)
		if err != nil {
			t.Fatalf("Failed to check inventory: %v", err)
		}
		if inventoryCount != 1 {
			t.Errorf("Expected item to be in inventory, count: %d", inventoryCount)
		}
	})

	t.Run("no item equipped in slot", func(t *testing.T) {
		newUserID := uuid.New()
		ts := time.Now().UnixNano()
		username := fmt.Sprintf("testuser%d", ts)
		userQuery := `INSERT INTO users (id, username, email, password, location_id, attack, defense, hp, current_hp, level)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
		_, err := db.Exec(userQuery, newUserID, username, fmt.Sprintf("test%d@example.com", ts), "password", user.ID, 1, 1, 20, 20, 5)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		err = service.TakeOffEquipmentItem(ctx, newUserID, "weapon")
		if !errors.Is(err, ErrNoItemEquipped) {
			t.Errorf("Expected ErrNoItemEquipped, got %v", err)
		}
	})

	t.Run("invalid slot name", func(t *testing.T) {
		err := service.TakeOffEquipmentItem(ctx, user.ID, "invalid_slot")
		if !errors.Is(err, ErrInvalidEquipmentType) {
			t.Errorf("Expected ErrInvalidEquipmentType, got %v", err)
		}
	})

	t.Run("unequip one item with multiple equipped", func(t *testing.T) {
		locationID := uuid.New()
		locationQuery := `INSERT INTO locations (id, name, slug, cell, inactive) VALUES ($1, $2, $3, $4, $5)`
		_, err := db.Exec(locationQuery, locationID, "Test Location 2", "test_location_2", false, false)
		if err != nil {
			t.Fatalf("Failed to create location: %v", err)
		}

		weaponCatID := uuid.New()
		categoryQuery := `INSERT INTO equipment_categories (id, name, type) VALUES ($1, $2, $3::equipment_category_type)`
		_, err = db.Exec(categoryQuery, weaponCatID, "Weapon", "weapon")
		if err != nil {
			t.Fatalf("Failed to create weapon category: %v", err)
		}

		chestCatID := uuid.New()
		_, err = db.Exec(categoryQuery, chestCatID, "Chest", "chest")
		if err != nil {
			t.Fatalf("Failed to create chest category: %v", err)
		}

		weaponID := uuid.New()
		itemQuery := `INSERT INTO equipment_items (id, name, slug, attack, defense, hp, required_level, price, equipment_category_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
		_, err = db.Exec(itemQuery, weaponID, "Test Weapon", "test-weapon", 10, 0, 0, 1, 100, weaponCatID)
		if err != nil {
			t.Fatalf("Failed to create weapon: %v", err)
		}

		chestID := uuid.New()
		_, err = db.Exec(itemQuery, chestID, "Test Chest", "test-chest", 0, 15, 30, 1, 100, chestCatID)
		if err != nil {
			t.Fatalf("Failed to create chest: %v", err)
		}

		multiUserID := uuid.New()
		ts := time.Now().UnixNano()
		username := fmt.Sprintf("multiuser%d", ts)
		userQuery := `INSERT INTO users (id, username, email, password, location_id, attack, defense, hp, current_hp, level, weapon_equipment_item_id, chest_equipment_item_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err = db.Exec(userQuery, multiUserID, username, fmt.Sprintf("multi%d@example.com", ts), "password", locationID, 11, 16, 50, 50, 5, weaponID, chestID)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		err = service.TakeOffEquipmentItem(ctx, multiUserID, "weapon")
		if err != nil {
			t.Fatalf("Failed to unequip weapon: %v", err)
		}

		type stats struct {
			Attack  uint `db:"attack"`
			Defense uint `db:"defense"`
			Hp      uint `db:"hp"`
		}
		var userStats stats
		statsQuery := `SELECT attack, defense, hp FROM users WHERE id = $1`
		err = db.Get(&userStats, statsQuery, multiUserID)
		if err != nil {
			t.Fatalf("Failed to get user stats: %v", err)
		}

		expectedAttack := uint(1)
		expectedDefense := uint(16)
		expectedHp := uint(50)

		if userStats.Attack != expectedAttack {
			t.Errorf("Expected attack %d, got %d", expectedAttack, userStats.Attack)
		}
		if userStats.Defense != expectedDefense {
			t.Errorf("Expected defense %d, got %d", expectedDefense, userStats.Defense)
		}
		if userStats.Hp != expectedHp {
			t.Errorf("Expected hp %d, got %d", expectedHp, userStats.Hp)
		}

		var chestEquippedID uuid.UUID
		chestQuery := `SELECT chest_equipment_item_id FROM users WHERE id = $1`
		err = db.Get(&chestEquippedID, chestQuery, multiUserID)
		if err != nil {
			t.Fatalf("Failed to get chest equipment: %v", err)
		}
		if chestEquippedID != chestID {
			t.Errorf("Expected chest to still be equipped, got %s", chestEquippedID)
		}
	})
}

