package main

import (
	"sunlight/config/database"
	"sunlight/handlers"

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

	e.POST("/auth/sign_up", handlers.SignUp)

	e.Start(":1323")
}
