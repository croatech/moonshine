package repository

import (
	"fmt"
	"testing"
	"time"

	"moonshine/internal/domain"
)

func TestUserRepository_Create(t *testing.T) {
	repo := NewUserRepository()
	ts := time.Now().UnixNano()

	user := &domain.User{
		Username: fmt.Sprintf("testuser%d", ts),
		Email:    fmt.Sprintf("test%d@example.com", ts),
		Password: "hashedpassword",
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Error("User ID should be set after creation")
	}
}

func TestUserRepository_FindByUsername(t *testing.T) {
	repo := NewUserRepository()
	ts := time.Now().UnixNano()

	username := fmt.Sprintf("finduser%d", ts)
	user := &domain.User{
		Username: username,
		Email:    fmt.Sprintf("find%d@example.com", ts),
		Password: "hashedpassword",
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
	repo := NewUserRepository()
	ts := time.Now().UnixNano()

	user := &domain.User{
		Username: fmt.Sprintf("iduser%d", ts),
		Email:    fmt.Sprintf("id%d@example.com", ts),
		Password: "hashedpassword",
	}
	if err := repo.Create(user); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	found, err := repo.FindByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to find user: %v", err)
	}

	if found.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, found.ID)
	}
}
