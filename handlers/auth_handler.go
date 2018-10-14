package handlers

import (
	"net/http"
	"sunlight/config/database"
	"sunlight/models"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c echo.Context) error {
	db := database.Connect()

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(c.FormValue("password")), bcrypt.DefaultCost)
	user := &models.User{
		Username: c.FormValue("username"),
		Email:    c.FormValue("email"),
		Password: string(passwordHash),
	}
	err := db.Create(user).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	} else {
		return c.JSON(http.StatusOK, user)
	}
}
