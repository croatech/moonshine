package repository

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"moonshine/internal/domain"
)

func TestUserRepository_Create(t *testing.T) {
	repo := NewUserRepository(testDB.DB())
	locationRepo := NewLocationRepository(testDB.DB())
	ts := time.Now().UnixNano()

	location := &domain.Location{
		Name:     fmt.Sprintf("Test Location %d", ts),
		Slug:     fmt.Sprintf("test-location-%d", ts),
		Cell:     false,
		Inactive: false,
	}
	if err := locationRepo.Create(location); err != nil {
		t.Fatalf("Failed to create location: %v", err)
	}

	user := &domain.User{
		Username:   fmt.Sprintf("testuser%d", ts),
		Email:      fmt.Sprintf("test%d@example.com", ts),
		Password:   "hashedpassword",
		LocationID: location.ID,
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID == uuid.Nil {
		t.Error("User ID should be set after creation")
	}
}

func TestUserRepository_FindByUsername(t *testing.T) {
	repo := NewUserRepository(testDB.DB())
	locationRepo := NewLocationRepository(testDB.DB())
	ts := time.Now().UnixNano()

	location := &domain.Location{
		Name:     fmt.Sprintf("Test Location %d", ts),
		Slug:     fmt.Sprintf("test-location-%d", ts),
		Cell:     false,
		Inactive: false,
	}
	if err := locationRepo.Create(location); err != nil {
		t.Fatalf("Failed to create location: %v", err)
	}

	username := fmt.Sprintf("finduser%d", ts)
	user := &domain.User{
		Username:   username,
		Email:      fmt.Sprintf("find%d@example.com", ts),
		Password:   "hashedpassword",
		LocationID: location.ID,
	}
	if err := repo.Create(user); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	found, err := repo.FindByUsername(username)
	if err != nil {
		t.Fatalf("Failed to find user: %v", err)
	}

	if found.Username != username {
		t.Errorf("Expected username '%s', got '%s'", username, found.Username)
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	repo := NewUserRepository(testDB.DB())
	locationRepo := NewLocationRepository(testDB.DB())
	ts := time.Now().UnixNano()

	location := &domain.Location{
		Name:     fmt.Sprintf("Test Location %d", ts),
		Slug:     fmt.Sprintf("test-location-%d", ts),
		Cell:     false,
		Inactive: false,
	}
	if err := locationRepo.Create(location); err != nil {
		t.Fatalf("Failed to create location: %v", err)
	}

	user := &domain.User{
		Username:   fmt.Sprintf("iduser%d", ts),
		Email:      fmt.Sprintf("id%d@example.com", ts),
		Password:   "hashedpassword",
		LocationID: location.ID,
	}
	if err := repo.Create(user); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	found, err := repo.FindByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to find user: %v", err)
	}

	if found.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, found.ID)
	}
}

func TestUserRepository_FindByUsername_NotFound(t *testing.T) {
	repo := NewUserRepository(testDB.DB())

	_, err := repo.FindByUsername("nonexistent")
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	repo := NewUserRepository(testDB.DB())

	_, err := repo.FindByID(uuid.New())
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}
