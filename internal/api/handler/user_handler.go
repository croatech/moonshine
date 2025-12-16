package handler

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"moonshine/internal/domain"
	"moonshine/internal/repository"
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

func currentUserByJwtToken(c echo.Context) (user domain.User, err error) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return user, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return user, errors.New("invalid claims")
	}

	idFloat, ok := claims["id"].(float64)
	if !ok {
		return user, errors.New("invalid user id")
	}

	userRepo := repository.NewUserRepository()
	foundUser, err := userRepo.FindByID(uint(idFloat))
	if err != nil {
		return user, errors.New("user not found")
	}
	user = *foundUser

	return user, nil
}
