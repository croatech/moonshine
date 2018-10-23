package seeds

import (
	"sunlight/models"
	"sunlight/modules/database"
)

func seedUsers() {
	user := models.User{
		Username: "Croaton",
		Email:    "admin@gmail.com",
		Password: models.HashPassword("password"),
	}
	database.Connection().Create(&user)
}
