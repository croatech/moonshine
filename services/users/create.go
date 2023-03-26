package services

import (
	"moonshine/models"
	"moonshine/modules/database"
)

func CreateUser(user *models.User) error {
	db := database.Connection()

	if err := db.Create(&user).Error; err != nil {
		db.Close()
		return err
	}

	defer db.Close()

	return nil
}
