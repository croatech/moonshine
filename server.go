package main

import (
	"net/http"
	"sunlight/config/database"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	database.Initialize()

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status}}"` + "\n",
	}))

	e.POST("/hello", hello)

	e.Start(":1323")
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
