package repository

import (
	"moonshine/internal/domain"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	repo := NewUserRepository()

	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Error("User ID should be set after creation")
	}
}

func TestUserRepository_FindByUsername(t *testing.T) {
	repo := NewUserRepository()

	// Create a test user first
	user := &domain.User{
		Username: "finduser",
		Email:    "find@example.com",
		Password: "hashedpassword",
	}
	if err := repo.Create(user); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Find the user
	found, err := repo.FindByUsername("finduser")
	if err != nil {
		t.Fatalf("Failed to find user: %v", err)
	}

	if found.Username != "finduser" {
		t.Errorf("Expected username 'finduser', got '%s'", found.Username)
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	repo := NewUserRepository()

	// Create a test user first
	user := &domain.User{
		Username: "iduser",
		Email:    "id@example.com",
		Password: "hashedpassword",
	}
	if err := repo.Create(user); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Find the user by ID
	found, err := repo.FindByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to find user: %v", err)
	}

	if found.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, found.ID)
	}
}

