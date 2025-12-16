package service

import (
	"moonshine/internal/domain"
	"moonshine/internal/repository"
	"testing"
)

func TestUserService_CreateUser(t *testing.T) {
	// Initialize repository for testing
	if err := repository.Init(); err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repository.Close()

	service := NewUserService()

	user := &domain.User{
		Username: "serviceuser",
		Email:    "service@example.com",
		Password: "hashedpassword",
	}

	created, err := service.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if created.ID == 0 {
		t.Error("User ID should be set after creation")
	}
}

func TestUserService_GetUserByUsername(t *testing.T) {
	// Initialize repository for testing
	if err := repository.Init(); err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repository.Close()

	service := NewUserService()

	// Create a test user first
	user := &domain.User{
		Username: "getuser",
		Email:    "get@example.com",
		Password: "hashedpassword",
	}
	if _, err := service.CreateUser(user); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Get the user
	found, err := service.GetUserByUsername("getuser")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if found.Username != "getuser" {
		t.Errorf("Expected username 'getuser', got '%s'", found.Username)
	}
}

