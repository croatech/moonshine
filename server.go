package main

import (
	"net/http"
	"sunlight/config"
	"sunlight/handlers"
	"sunlight/modules/migrations"
	"sunlight/modules/seeds"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	migrations.Run()
	seeds.Load()
	config.Load()

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"response":"${latency_human}", time":"${time_rfc3339_nano}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status}}"` + "\n",
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Unauthenticated routes
	e.POST("/auth/sign_up", handlers.SignUp)
	e.POST("/auth/sign_in", handlers.SignIn)

	// Restricted routes
	r := e.Group("/")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("users/current", handlers.CurrentUser)

	e.Start(":1323")
}
