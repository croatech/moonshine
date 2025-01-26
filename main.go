package main

import (
	"github.com/joho/godotenv"
	"log"
	"moonshine/modules/database"
	"moonshine/modules/seeds"
	"moonshine/modules/server"
	"os"
)

func main() {
	// Load envs
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Drop()
	database.Migrate()
	seeds.Load()

	app := server.AppServer()
	app.Start(":" + os.Getenv("APP_PORT"))
	if err != nil {
		panic(err)
	}
}
