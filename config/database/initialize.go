package database

import (
	"github.com/jinzhu/gorm"
	"feed/models"
)

func Initialize() *gorm.DB {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=vitaliy dbname=go_moonlight password=password sslmode=disable")

	if err != nil {
		panic(err)
	}

	migrate(db)

	return db
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
}
