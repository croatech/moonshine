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

func setupTestData(db *sqlx.DB) (*domain.User, *domain.EquipmentItem, uuid.UUID, error) {
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
	userQuery := `INSERT INTO users (id, username, email, password, location_id, attack, defense, hp, current_hp, level)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	ts := time.Now().UnixNano()
	username := fmt.Sprintf("testuser%d", ts)
	_, err = db.Exec(userQuery, userID, username, fmt.Sprintf("test%d@example.com", ts), "password", locationID, 1, 1, 20, 20, 5)
	if err != nil {
		return nil, nil, uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	inventoryQuery := `INSERT INTO user_equipment_items (id, user_id, equipment_item_id) VALUES ($1, $2, $3)`
	_, err = db.Exec(inventoryQuery, uuid.New(), userID, itemID)
	if err != nil {
		return nil, nil, uuid.Nil, fmt.Errorf("failed to add item to inventory: %w", err)
	}

	user := &domain.User{
		Model: domain.Model{
			ID: userID,
		},
		Username: username,
		Level:    5,
		Attack:   1,
		Defense:  1,
		Hp:       20,
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

func TestEquipmentItemTakeOnService_TakeOnEquipmentItem(t *testing.T) {
	if testDB == nil {
		t.Skip("Test database not initialized")
	}

	db := testDB.DB()
	ctx := context.Background()

	user, item, categoryID, err := setupTestData(db)
	if err != nil {
		t.Fatalf("Failed to setup test data: %v", err)
	}

	equipmentItemRepo := repository.NewEquipmentItemRepository(db)
	userEquipmentItemRepo := repository.NewUserEquipmentItemRepository(db)
	userRepo := repository.NewUserRepository(db)
	service := NewEquipmentItemTakeOnService(db, equipmentItemRepo, userEquipmentItemRepo, userRepo)

	t.Run("successfully equip item", func(t *testing.T) {
		err := service.TakeOnEquipmentItem(ctx, user.ID, item.ID)
		if err != nil {
			t.Fatalf("Failed to equip item: %v", err)
		}

		var equippedItemID uuid.UUID
		query := `SELECT weapon_equipment_item_id FROM users WHERE id = $1`
		err = db.Get(&equippedItemID, query, user.ID)
		if err != nil {
			t.Fatalf("Failed to get equipped item: %v", err)
		}

		if equippedItemID != item.ID {
			t.Errorf("Expected equipped item ID %s, got %s", item.ID, equippedItemID)
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

		expectedAttack := uint(11)
		expectedDefense := uint(6)
		expectedHp := uint(40)

		if userStats.Attack != expectedAttack {
			t.Errorf("Expected attack %d, got %d", expectedAttack, userStats.Attack)
		}
		if userStats.Defense != expectedDefense {
			t.Errorf("Expected defense %d, got %d", expectedDefense, userStats.Defense)
		}
		if userStats.Hp != expectedHp {
			t.Errorf("Expected hp %d, got %d", expectedHp, userStats.Hp)
		}
	})

	t.Run("item not in inventory", func(t *testing.T) {
		newItemID := uuid.New()
		newItemQuery := `INSERT INTO equipment_items (id, name, slug, attack, defense, hp, required_level, price, equipment_category_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
		_, err := db.Exec(newItemQuery, newItemID, "New Sword", "new-sword", 10, 5, 20, 1, 100, categoryID)
		if err != nil {
			t.Fatalf("Failed to create new item: %v", err)
		}

		err = service.TakeOnEquipmentItem(ctx, user.ID, newItemID)
		if !errors.Is(err, ErrItemNotInInventory) {
			t.Errorf("Expected ErrItemNotInInventory, got %v", err)
		}
	})

	t.Run("insufficient level", func(t *testing.T) {
		highLevelItemID := uuid.New()
		highLevelItemQuery := `INSERT INTO equipment_items (id, name, slug, attack, defense, hp, required_level, price, equipment_category_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
		_, err := db.Exec(highLevelItemQuery, highLevelItemID, "High Level Sword", "high-level-sword", 10, 5, 20, 10, 100, categoryID)
		if err != nil {
			t.Fatalf("Failed to create high level item: %v", err)
		}

		inventoryQuery := `INSERT INTO user_equipment_items (id, user_id, equipment_item_id) VALUES ($1, $2, $3)`
		_, err = db.Exec(inventoryQuery, uuid.New(), user.ID, highLevelItemID)
		if err != nil {
			t.Fatalf("Failed to add item to inventory: %v", err)
		}

		err = service.TakeOnEquipmentItem(ctx, user.ID, highLevelItemID)
		if !errors.Is(err, ErrInsufficientLevel) {
			t.Errorf("Expected ErrInsufficientLevel, got %v", err)
		}
	})

	t.Run("replace existing equipment", func(t *testing.T) {
		newItemID := uuid.New()
		newItemQuery := `INSERT INTO equipment_items (id, name, slug, attack, defense, hp, required_level, price, equipment_category_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
		_, err := db.Exec(newItemQuery, newItemID, "New Sword 2", "new-sword-2", 15, 8, 25, 1, 100, categoryID)
		if err != nil {
			t.Fatalf("Failed to create new item: %v", err)
		}

		inventoryQuery := `INSERT INTO user_equipment_items (id, user_id, equipment_item_id) VALUES ($1, $2, $3)`
		_, err = db.Exec(inventoryQuery, uuid.New(), user.ID, newItemID)
		if err != nil {
			t.Fatalf("Failed to add item to inventory: %v", err)
		}

		err = service.TakeOnEquipmentItem(ctx, user.ID, newItemID)
		if err != nil {
			t.Fatalf("Failed to equip new item: %v", err)
		}

		var equippedItemID uuid.UUID
		query := `SELECT weapon_equipment_item_id FROM users WHERE id = $1`
		err = db.Get(&equippedItemID, query, user.ID)
		if err != nil {
			t.Fatalf("Failed to get equipped item: %v", err)
		}

		if equippedItemID != newItemID {
			t.Errorf("Expected equipped item ID %s, got %s", newItemID, equippedItemID)
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

		expectedAttack := uint(16)
		expectedDefense := uint(9)
		expectedHp := uint(45)

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
			t.Errorf("Expected old item to be in inventory, count: %d", inventoryCount)
		}
	})
}

