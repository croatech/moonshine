package tests

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"moonshine/handlers"
	"moonshine/modules/database"
	"net/http"
	"os"
	"testing"
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

	// Unauthenticated routes
	app.POST("/auth/sign_up", handlers.SignUp)
	app.POST("/auth/sign_in", handlers.SignIn)

	// Restricted routes
	r := app.Group("/")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("users/current", handlers.CurrentUser)

	return app
}

func TestMain(m *testing.M) {
	loadTestConfig()
	createTables()

	print("#################################")

	code := m.Run()

	print("#################################")

	dropTables()

	os.Exit(code)
}

func createTables() {
	database.Migrate()
}

func dropTables() {
	database.Drop()
}

func loadTestConfig() {
	viper.SetConfigFile("config_test.yml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
