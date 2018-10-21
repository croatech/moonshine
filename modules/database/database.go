package database

import (
	"sunlight/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var connection *gorm.DB

func Prepare() error {
	conn := connect()

	migrate(conn)

	defer conn.Close()

	return nil
}

func Connection() *gorm.DB {
	conn := connect()
	conn.LogMode(true)

	return conn
}

func connect() *gorm.DB {
	conn, err := gorm.Open("postgres", "host=localhost port=5432 user=vitaliy dbname=go_moonlight password=password sslmode=disable")

	if err != nil {
		panic(err)
	}

	return conn
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
}
