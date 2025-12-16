package graphql

import (
	"context"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"

	"moonshine/internal/graphql/generated"
)

func GraphQLHandler() echo.HandlerFunc {
	resolver := NewResolver()
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{
			Resolvers: resolver,
		}),
	)

	if isProduction() {
		srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
			opCtx := graphql.GetOperationContext(ctx)
			if opCtx != nil {
				query := opCtx.RawQuery
				if strings.Contains(strings.ToLower(query), "__schema") ||
					strings.Contains(strings.ToLower(query), "__type") ||
					opCtx.OperationName == "IntrospectionQuery" {
					return func(ctx context.Context) *graphql.Response {
						return &graphql.Response{
							Errors: graphql.GetErrors(ctx),
						}
					}
				}
			}
			return next(ctx)
		})
	}

	return func(c echo.Context) error {
		ctx := c.Request().Context()

		if tokenValue := c.Get("user"); tokenValue != nil {
			if token, ok := tokenValue.(*jwt.Token); ok {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if idFloat, ok := claims["id"].(float64); ok {
						ctx = context.WithValue(ctx, "userID", uint(idFloat))
						c.SetRequest(c.Request().WithContext(ctx))
					}
				}
			}
		}

		srv.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

func SchemaHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		if isProduction() {
			return c.String(404, "Not Found")
		}

		c.Response().Header().Set("Content-Type", "text/plain; charset=utf-8")
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request().Method == "OPTIONS" {
			return c.NoContent(204)
		}

		schema := `scalar Time

type User {
  id: ID!
  username: String!
  email: String!
  hp: Int!
  level: Int!
  gold: Int!
  exp: Int!
  createdAt: Time!
  updatedAt: Time!
}

type AuthPayload {
  token: String!
  user: User!
}

type Query {
  currentUser: User
}

type Mutation {
  signUp(input: SignUpInput!): AuthPayload!
  signIn(input: SignInInput!): AuthPayload!
}

input SignUpInput {
  username: String!
  email: String!
  password: String!
}

input SignInInput {
  username: String!
  password: String!
}
`
		return c.String(200, schema)
	}
}

func isProduction() bool {
	env := os.Getenv("ENV")
	return env == "production" || env == "prod"
}
