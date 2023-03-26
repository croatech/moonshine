package main

import (
	"github.com/joho/godotenv"
	"github.com/steinfletcher/apitest"
	"log"
	"moonshine/modules/database"
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
	services.CreateUser("test", "a@gmail.com", "password")

	apitest.New().
		Handler(appServer()).
		Post("/auth/sign_up").
		JSON(`{"username":"cro","password":"password","email":"a@gmail.com"}`).
		Expect(t).
		Body(`"Email or username already exists"`).
		Status(http.StatusInternalServerError).
		End()

	database.Clean()
}

func TestSignUp_FailNotUniqueUsername(t *testing.T) {
	services.CreateUser("cro", "an@gmail.com", "password")

	apitest.New().
		Handler(appServer()).
		Post("/auth/sign_up").
		JSON(`{"username":"cro","password":"password","email":"a@gmail.com"}`).
		Expect(t).
		Body(`"Email or username already exists"`).
		Status(http.StatusInternalServerError).
		End()

	database.Clean()
}
