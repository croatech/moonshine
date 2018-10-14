package database

import (
	"sunlight/models"

	"github.com/jinzhu/gorm"
)

func Initialize() *gorm.DB {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=vitaliy dbname=go_moonlight password=password sslmode=disable")

	if err != nil {
		panic(err)
	}

	migrate(db)

	return db
}

func Connect() *gorm.DB {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=vitaliy dbname=go_moonlight password=password sslmode=disable")

	if err != nil {
		panic(err)
	}

	return db
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
}
