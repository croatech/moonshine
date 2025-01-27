package test

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"moonshine/handlers"
	"moonshine/models"
	"moonshine/modules/server"
	"moonshine/modules/support"
	services "moonshine/services/users"
)

func TestUsersCurrent_Success(t *testing.T) {
	user := models.User{
		Username: "cro",
		Email:    "admin@gmail.com",
		Password: support.HashPassword("password"),
	}

	createdUser, _ := services.CreateUser(&user)

	token, _ := handlers.GenerateJwtPayload(createdUser.ID)

	apitest.New().
		Handler(server.AppServer()).
		Get("/users/current").
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Status(http.StatusOK).
		End()

	CleanDatabase()
}

func TestUsersCurrent_Fail(t *testing.T) {
	SeedUsers()

	token, _ := handlers.GenerateJwtPayload(0)

	apitest.New().
		Handler(server.AppServer()).
		Get("/users/current").
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Status(http.StatusUnauthorized).
		End()

	CleanDatabase()
}
