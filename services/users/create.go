package services

import (
	"moonshine/models"
	"moonshine/modules/database"
	"moonshine/modules/support"
)

func CreateUser(username string, email string, password string) error {
	user := models.User{
		Username: username,
		Email:    email,
		Password: support.HashPassword(password),
	}

	db := database.Connection()

	if err := db.Create(&user).Error; err != nil {
		db.Close()
		return err
	}

	defer db.Close()

	return nil
}
