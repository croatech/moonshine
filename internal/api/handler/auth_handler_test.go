package handler

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/steinfletcher/apitest"
)

func TestSignUp_Success(t *testing.T) {
	e := setupTestServer()
	timestamp := time.Now().Unix()
	username := fmt.Sprintf("u%d", timestamp)
	email := fmt.Sprintf("e%d@test.com", timestamp)
	
	apitest.New().
		Handler(e).
		Post("/signup").
		JSON(fmt.Sprintf(`{"username":"%s","password":"password","email":"%s"}`, username, email)).
		Expect(t).
		Body(`""`).
		Status(http.StatusOK).
		End()
}

func TestSignUp_FailNotUniqueEmail(t *testing.T) {
	e := setupTestServer()
	timestamp := time.Now().Unix()
	email := fmt.Sprintf("dup%d@test.com", timestamp)
	username1 := fmt.Sprintf("u1%d", timestamp)
	username2 := fmt.Sprintf("u2%d", timestamp)
	
	apitest.New().
		Handler(e).
		Post("/signup").
		JSON(fmt.Sprintf(`{"username":"%s","password":"password","email":"%s"}`, username1, email)).
		Expect(t).
		Status(http.StatusOK).
		End()

	apitest.New().
		Handler(e).
		Post("/signup").
		JSON(fmt.Sprintf(`{"username":"%s","password":"password","email":"%s"}`, username2, email)).
		Expect(t).
		Body(`"Email or username already exists"`).
		Status(http.StatusInternalServerError).
		End()
}

func TestSignUp_FailNotUniqueUsername(t *testing.T) {
	e := setupTestServer()
	timestamp := time.Now().Unix()
	username := fmt.Sprintf("dup%d", timestamp)
	email1 := fmt.Sprintf("e1%d@test.com", timestamp)
	email2 := fmt.Sprintf("e2%d@test.com", timestamp)
	
	apitest.New().
		Handler(e).
		Post("/signup").
		JSON(fmt.Sprintf(`{"username":"%s","password":"password","email":"%s"}`, username, email1)).
		Expect(t).
		Status(http.StatusOK).
		End()

	apitest.New().
		Handler(e).
		Post("/signup").
		JSON(fmt.Sprintf(`{"username":"%s","password":"password","email":"%s"}`, username, email2)).
		Expect(t).
		Body(`"Email or username already exists"`).
		Status(http.StatusInternalServerError).
		End()
}

func TestSignIn_Success(t *testing.T) {
	e := setupTestServer()
	timestamp := time.Now().Unix()
	username := fmt.Sprintf("si%d", timestamp)
	email := fmt.Sprintf("si%d@test.com", timestamp)
	
	apitest.New().
		Handler(e).
		Post("/signup").
		JSON(fmt.Sprintf(`{"username":"%s","password":"password","email":"%s"}`, username, email)).
		Expect(t).
		Status(http.StatusOK).
		End()

	apitest.New().
		Handler(e).
		Post("/signin").
		JSON(fmt.Sprintf(`{"username":"%s","password":"password"}`, username)).
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestSignIn_Fail(t *testing.T) {
	e := setupTestServer()
	timestamp := time.Now().Unix()
	username := fmt.Sprintf("nx%d", timestamp)
	
	apitest.New().
		Handler(e).
		Post("/signin").
		JSON(fmt.Sprintf(`{"username":"%s","password":"password"}`, username)).
		Expect(t).
		Body(`"User not found"`).
		Status(http.StatusUnauthorized).
		End()
}

