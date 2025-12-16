package api

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"moonshine/internal/api/handler"
)

func SetupRoutes(e *echo.Echo) {
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	e.POST("/signup", handler.SignUp)
	e.POST("/signin", handler.SignIn)
	jwtConfig := echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_KEY")),
	}
	protected := e.Group("")
	protected.Use(echojwt.WithConfig(jwtConfig))
	protected.GET("/user", handler.CurrentUser)
}

