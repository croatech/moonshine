package main

import (
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	"log"
	"moonshine/handlers"
	"moonshine/modules/database"
	"moonshine/modules/seeds"
	"net/http"
	"os"
)

func appServer() *echo.Echo {
	app := echo.New()
	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"response":"${latency_human}", time":"${time_rfc3339_nano}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status}}"` + "\n",
	}))
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Swagger
	app.GET("/swagger/*", echoSwagger.WrapHandler)

	// Unauthenticated routes
	app.POST("/auth/sign_up", handlers.SignUp)
	app.POST("/auth/sign_in", handlers.SignIn)

	// Restricted routes
	r := app.Group("/")
	r.Use(echojwt.JWT([]byte(os.Getenv("JWT_KEY"))))
	r.GET("users/current", handlers.CurrentUser)

	return app
}

func main() {
	// Load envs
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Drop()
	database.Migrate()
	seeds.Load()

	app := appServer()
	app.Start(":" + os.Getenv("APP_PORT"))
	if err != nil {
		panic(err)
	}
}
