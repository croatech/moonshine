package main

import (
	"github.com/joho/godotenv"
	"github.com/steinfletcher/apitest"
	"log"
	"moonshine/modules/database"
	"moonshine/modules/seeds"
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

	database.Migrate()

	exitVal := m.Run()

	database.Drop()

	os.Exit(exitVal)

	log.Println("Tests finished")
}

func TestSignUp_Success(t *testing.T) {
	apitest.New().
		Handler(appServer()).
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
		Handler(appServer()).
		Post("/auth/sign_up").
		JSON(`{"username":"Croa","password":"password","email":"admin@gmail.com"}`).
		Expect(t).
		Body(`"Email or username already exists"`).
		Status(http.StatusInternalServerError).
		End()

	database.Clean()
}

func TestSignUp_FailNotUniqueUsername(t *testing.T) {
	seeds.SeedUsers()

	apitest.New().
		Handler(appServer()).
		Post("/auth/sign_up").
		JSON(`{"username":"Cro","password":"password","email":"a@gmail.com"}`).
		Expect(t).
		Body(`"Email or username already exists"`).
		Status(http.StatusInternalServerError).
		End()

	database.Clean()
}

func TestSignIn_Success(t *testing.T) {
	apitest.New().
		Handler(appServer()).
		Post("/auth/sign_up").
		JSON(`{"username":"cro","password":"password","email":"a@gmail.com"}`).
		Expect(t).
		Body(`""`).
		Status(http.StatusOK).
		End()

	database.Clean()
}
