package handler

import (
	"moonshine/internal/repository"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
)

func setupTestServer() *echo.Echo {
	if os.Getenv("JWT_KEY") == "" {
		os.Setenv("JWT_KEY", "test-secret-key")
	}

	if repository.DB == nil {
		if err := repository.Init(); err != nil {
			panic(err)
		}
	}

	e := echo.New()
	
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	e.POST("/signup", SignUp)
	e.POST("/signin", SignIn)
	
	jwtConfig := echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_KEY")),
	}
	protected := e.Group("")
	protected.Use(echojwt.WithConfig(jwtConfig))
	protected.GET("/user", CurrentUser)
	
	return e
}

