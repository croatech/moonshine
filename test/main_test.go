package test

import (
	"log"
	"moonshine/modules/database"
	"moonshine/modules/seeds"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	log.Println("Tests started")

	err := godotenv.Load("../.env.test")
	if err != nil {
		log.Fatal(err)
	}

	database.Drop()
	database.Migrate()

	exitVal := m.Run()

	os.Exit(exitVal)

	log.Println("Tests finished")
}

func CleanDatabase() {
	database.Clean()
}

func SeedUsers() {
	seeds.SeedUsers()
}
