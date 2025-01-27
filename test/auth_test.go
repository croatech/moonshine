package test

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"moonshine/modules/server"
)

func TestSignUp_Success(t *testing.T) {
	apitest.New().
		Handler(server.AppServer()).
		Post("/auth/sign_up").
		JSON(`{"username":"cro","password":"password","email":"a@gmail.com"}`).
		Expect(t).
		Body(`""`).
		Status(http.StatusOK).
		End()

	CleanDatabase()
}

func TestSignUp_FailNotUniqueEmail(t *testing.T) {
	SeedUsers()

	apitest.New().
		Handler(server.AppServer()).
		Post("/auth/sign_up").
		JSON(`{"username":"croa","password":"password","email":"admin@gmail.com"}`).
		Expect(t).
		Body(`"Email or username already exists"`).
		Status(http.StatusInternalServerError).
		End()

	CleanDatabase()
}

func TestSignUp_FailNotUniqueUsername(t *testing.T) {
	SeedUsers()

	apitest.New().
		Handler(server.AppServer()).
		Post("/auth/sign_up").
		JSON(`{"username":"cro","password":"password","email":"a@gmail.com"}`).
		Expect(t).
		Body(`"Email or username already exists"`).
		Status(http.StatusInternalServerError).
		End()

	CleanDatabase()
}

func TestSignIn_Success(t *testing.T) {
	SeedUsers()

	apitest.New().
		Handler(server.AppServer()).
		Post("/auth/sign_in").
		JSON(`{"username":"cro","password":"password"}`).
		Expect(t).
		Status(http.StatusOK).
		End()

	CleanDatabase()
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
