package api

import (
	"net/http"
	"os"
	"strings"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"moonshine/internal/graphql"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB, isProduction bool) {
	e.GET("/health", healthCheck)

	if !isProduction {
		e.GET("/schema.graphql", graphql.SchemaHandler(isProduction))
		e.OPTIONS("/schema.graphql", graphql.SchemaHandler(isProduction))
	}

	graphqlGroup := e.Group("/graphql")
	graphqlGroup.Use(cacheRequestBody())

	jwtConfig := echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_KEY")),
		ErrorHandler: func(c echo.Context, err error) error {
			return nil
		},
		Skipper: func(c echo.Context) bool {
			if isProduction {
				return isPublicOperation(c)
			}
			return isIntrospectionQuery(c) || isPublicOperation(c)
		},
	}

	graphqlGroup.Use(echojwt.WithConfig(jwtConfig))
	graphqlGroup.POST("", graphql.GraphQLHandler(db, isProduction))
}

func isIntrospectionQuery(c echo.Context) bool {
	bodyBytes := getCachedBody(c)
	if bodyBytes == nil {
		return false
	}

	body := strings.ToLower(string(bodyBytes))
	patterns := []string{"__schema", "__type", "introspection", "query introspection"}

	for _, pattern := range patterns {
		if strings.Contains(body, pattern) {
			return true
		}
	}

	return false
}

func isPublicOperation(c echo.Context) bool {
	bodyBytes := getCachedBody(c)
	if bodyBytes == nil {
		return false
	}

	body := strings.ToLower(string(bodyBytes))
	publicOps := []string{"signup", "signin"}

	for _, op := range publicOps {
		if strings.Contains(body, op) {
			return true
		}
	}

	return false
}

func healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
