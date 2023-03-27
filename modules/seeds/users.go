package seeds

import (
	"moonshine/models"
	"moonshine/modules/support"
	services "moonshine/services/users"
)

func SeedUsers() {
	user := models.User{
		Username: "cro",
		Email:    "admin@gmail.com",
		Password: support.HashPassword("password"),
	}

	services.CreateUser(&user)
}
