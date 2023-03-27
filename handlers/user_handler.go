package handlers

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"moonshine/models"
	"moonshine/modules/database"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserResponse struct {
	Username string `json:"username"`
	Hp       uint   `json:"hp"`
}

func CurrentUser(c echo.Context) error {
	user, err := currentUserByJwtToken(c)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, "User not found")
	}

	response := UserResponse{
		Username: user.Username,
		Hp:       user.Hp,
	}

	return c.JSON(http.StatusOK, response)
}

func currentUserByJwtToken(c echo.Context) (user models.User, err error) {
	// Get username by Jwt token
	result := c.Get("user").(*jwt.Token)
	claims := result.Claims.(jwt.MapClaims)
	id := claims["id"].(float64)

	if database.Connection().Where("id = ?", uint(id)).First(&user).RecordNotFound() {
		return user, errors.New("user not found")
	}

	return user, nil
}
