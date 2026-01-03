package services

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"

	"moonshine/internal/repository"
)

var testDB *repository.Database

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../.env.test")
	if err != nil {
		log.Printf("Warning: Failed to load .env.test: %v", err)
	}

	db, err := repository.New()
	if err != nil {
		log.Printf("Warning: Failed to initialize test database: %v", err)
		log.Printf("Tests will be skipped")
		os.Exit(0)
	}
	testDB = db

	code := m.Run()

	testDB.Close()
	os.Exit(code)
}

