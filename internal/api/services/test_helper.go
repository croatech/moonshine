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
	_ = godotenv.Load("../../../.env.test")

	db, err := repository.New()
	if err != nil {
		log.Fatalf("Failed to initialize test database: %v", err)
	}
	testDB = db

	code := m.Run()

	testDB.Close()
	os.Exit(code)
}

