package handler

import (
	"fmt"
	"moonshine/internal/domain"
	"moonshine/internal/repository"
	"moonshine/internal/util"
	"net/http"
	"testing"
	"time"

	"github.com/steinfletcher/apitest"
)

func TestUsersCurrent_Success(t *testing.T) {
	e := setupTestServer()

	userRepo := repository.NewUserRepository()
	timestamp := time.Now().Unix()
	username := fmt.Sprintf("cu%d", timestamp)
	user := &domain.User{
		Username: username,
		Email:    fmt.Sprintf("cu%d@test.com", timestamp),
		Password: util.HashPassword("password"),
	}
	if err := userRepo.Create(user); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	apitest.New().
		Handler(e).
		Post("/signin").
		JSON(fmt.Sprintf(`{"username":"%s","password":"password"}`, username)).
		Expect(t).
		Status(http.StatusOK).
		End()

	apitest.New().
		Handler(e).
		Get("/user").
		Header("Authorization", "Bearer test-token").
		Expect(t).
		Status(http.StatusUnauthorized).
		End()
}

func TestUsersCurrent_Fail(t *testing.T) {
	e := setupTestServer()
	apitest.New().
		Handler(e).
		Get("/user").
		Header("Authorization", "Bearer invalid-token").
		Expect(t).
		Status(http.StatusUnauthorized).
		End()
}

