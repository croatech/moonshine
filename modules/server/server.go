package server

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	"moonshine/handlers"
)

func AppServer() *echo.Echo {
	app := echo.New()

	// Middleware
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

	// Public routes
	app.POST("/auth/sign_up", handlers.SignUp)
	app.POST("/auth/sign_in", handlers.SignIn)

	// Protected routes
	r := app.Group("/")
	r.Use(echojwt.JWT([]byte(os.Getenv("JWT_KEY"))))
	r.GET("users/current", handlers.CurrentUser)

	return app
}
