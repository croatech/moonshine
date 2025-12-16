package service

import (
	"fmt"
	"testing"
	"time"

	"moonshine/internal/domain"
	"moonshine/internal/repository"
)

func TestMain(m *testing.M) {
	if err := repository.Init(); err != nil {
		panic(err)
	}
	defer repository.Close()
	m.Run()
}

func TestUserService_Create(t *testing.T) {
	svc := NewUserService()
	ts := time.Now().UnixNano()

	user := &domain.User{
		Username: fmt.Sprintf("serviceuser%d", ts),
		Email:    fmt.Sprintf("service%d@example.com", ts),
		Password: "hashedpassword",
	}

	if err := svc.Create(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Error("User ID should be set after creation")
	}
}

func TestUserService_GetByUsername(t *testing.T) {
	svc := NewUserService()
	ts := time.Now().UnixNano()

	username := fmt.Sprintf("getuser%d", ts)
	user := &domain.User{
		Username: username,
		Email:    fmt.Sprintf("get%d@example.com", ts),
		Password: "hashedpassword",
	}
	if err := svc.Create(user); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	found, err := svc.GetByUsername(username)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if found.Username != username {
		t.Errorf("Expected username '%s', got '%s'", username, found.Username)
	}
}

func TestUserService_GetByID(t *testing.T) {
	svc := NewUserService()
	ts := time.Now().UnixNano()

	user := &domain.User{
		Username: fmt.Sprintf("getbyid%d", ts),
		Email:    fmt.Sprintf("getbyid%d@example.com", ts),
		Password: "hashedpassword",
	}
	if err := svc.Create(user); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	found, err := svc.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if found.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, found.ID)
	}
}
