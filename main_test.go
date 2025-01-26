package main

import (
	"github.com/joho/godotenv"
	"github.com/steinfletcher/apitest"
	"log"
	"moonshine/handlers"
	"moonshine/models"
	"moonshine/modules/database"
	"moonshine/modules/seeds"
	"moonshine/modules/server"
	"moonshine/modules/support"
	services "moonshine/services/users"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.Println("Tests started")

	err := godotenv.Load(".env.test")
	if err != nil {
		log.Fatal(err)
	}

	database.Drop()
	database.Migrate()

	exitVal := m.Run()

	os.Exit(exitVal)

	log.Println("Tests finished")
}

// /auth/sign_up

func TestSignUp_Success(t *testing.T) {
	apitest.New().
		Handler(server.AppServer()).
		Post("/auth/sign_up").
		JSON(`{"username":"cro","password":"password","email":"a@gmail.com"}`).
		Expect(t).
		Body(`""`).
		Status(http.StatusOK).
		End()

	database.Clean()
}

func TestSignUp_FailNotUniqueEmail(t *testing.T) {
	seeds.SeedUsers()

	apitest.New().
		Handler(server.AppServer()).
		Post("/auth/sign_up").
		JSON(`{"username":"croa","password":"password","email":"admin@gmail.com"}`).
		Expect(t).
		Body(`"Email or username already exists"`).
		Status(http.StatusInternalServerError).
		End()

	database.Clean()
}

func TestSignUp_FailNotUniqueUsername(t *testing.T) {
	seeds.SeedUsers()

	apitest.New().
		Handler(server.AppServer()).
		Post("/auth/sign_up").
		JSON(`{"username":"cro","password":"password","email":"a@gmail.com"}`).
		Expect(t).
		Body(`"Email or username already exists"`).
		Status(http.StatusInternalServerError).
		End()

	database.Clean()
}

// /auth/sign_in

func TestSignIn_Success(t *testing.T) {
	seeds.SeedUsers()

	apitest.New().
		Handler(server.AppServer()).
		Post("/auth/sign_in").
		JSON(`{"username":"cro","password":"password"}`).
		Expect(t).
		Status(http.StatusOK).
		End()

	database.Clean()
}

func TestSignIn_Fail(t *testing.T) {
	apitest.New().
		Handler(server.AppServer()).
		Post("/auth/sign_in").
		JSON(`{"username":"cro","password":"password"}`).
		Expect(t).
		Body(`"User not found"`).
		Status(http.StatusUnauthorized).
		End()
}

// users/current
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

	database.Clean()
}

func TestUsersCurrent_Fail(t *testing.T) {
	seeds.SeedUsers()

	token, _ := handlers.GenerateJwtPayload(0)

	apitest.New().
		Handler(server.AppServer()).
		Get("/users/current").
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Status(http.StatusUnauthorized).
		End()

	database.Clean()
}
