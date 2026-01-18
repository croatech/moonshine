package services

import (
	"os"
	"testing"

	"github.com/joho/godotenv"

	"moonshine/internal/repository"
)

var testDB *repository.Database

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../.env.test")
	if err != nil {
	}

	db, err := repository.New()
	if err != nil {
		os.Exit(0)
	}
	testDB = db

	code := m.Run()

	testDB.Close()
	os.Exit(code)
}

