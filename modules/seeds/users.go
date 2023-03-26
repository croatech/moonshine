package seeds

import (
	"moonshine/models"
	"moonshine/modules/database"
	"moonshine/modules/support"
)

func seedUsers() {
	user := models.User{
		Username: "Croaton",
		Email:    "admin@gmail.com",
		Password: support.HashPassword("password"),
	}
	database.Connection().Create(&user)
}
