package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"moonshine/models"
	"moonshine/modules/database"
	"net/http"

	"github.com/labstack/echo"
)

func CurrentUser(c echo.Context) error {
	user := currentUserByJwtToken(c)
	return c.JSON(http.StatusOK, user)
}

func currentUserByJwtToken(c echo.Context) (user models.User) {
	// Get username by Jwt token
	result := c.Get("user").(*jwt.Token)
	claims := result.Claims.(jwt.MapClaims)
	username := claims["username"].(string)

	database.Connection().Where("username = ?", username).First(&user)
	return user
}
