package handlers

import (
	"net/http"
	"sunlight/models"

	"github.com/labstack/echo"
)

func CurrentUser(c echo.Context) error {
	user := models.CurrentUserByJwtToken(c)
	return c.JSON(http.StatusOK, user)
}
