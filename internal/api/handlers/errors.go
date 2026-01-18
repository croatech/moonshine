package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func ErrUnauthorized(c echo.Context) error {
	return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
}

func ErrNotFound(c echo.Context, message string) error {
	if message == "" {
		message = "not found"
	}
	return c.JSON(http.StatusNotFound, map[string]string{"error": message})
}

func ErrBadRequest(c echo.Context, message string) error {
	if message == "" {
		message = "invalid request"
	}
	return c.JSON(http.StatusBadRequest, map[string]string{"error": message})
}

func ErrInternalServerError(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
}

func ErrConflict(c echo.Context, message string) error {
	if message == "" {
		message = "conflict"
	}
	return c.JSON(http.StatusConflict, map[string]string{"error": message})
}

func ErrUnauthorizedWithMessage(c echo.Context, message string) error {
	return c.JSON(http.StatusUnauthorized, map[string]string{"error": message})
}
