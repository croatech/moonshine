package main

import (
	"net/http"

	"github.com/labstack/echo"
	_ "github.com/jinzhu/gorm/dialects/postgres"
		"feed/config/database"
)

func main() {
	database.Initialize()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
