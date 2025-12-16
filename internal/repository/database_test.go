package repository

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	log.Println("Repository tests started")

	err := godotenv.Load("../../.env.test")
	if err != nil {
		log.Println("Warning: .env.test not found, using environment variables")
	}

	if err := Init(); err != nil {
		log.Fatalf("Failed to initialize test database: %v", err)
	}

	exitVal := m.Run()

	Close()

	os.Exit(exitVal)

	log.Println("Repository tests finished")
}
