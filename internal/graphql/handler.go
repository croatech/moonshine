package graphql

import (
	"context"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"moonshine/internal/graphql/generated"
	"moonshine/internal/repository"
)

func GraphQLHandler(db *gorm.DB, isProduction bool) echo.HandlerFunc {
	userRepo := repository.NewUserRepository(db)
	resolver := newResolver(userRepo)
	
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{
			Resolvers: resolver,
		}),
	)

	if isProduction {
		srv.AroundOperations(blockIntrospection)
	}

	return func(c echo.Context) error {
		ctx := c.Request().Context()

		if tokenValue := c.Get("user"); tokenValue != nil {
			if token, ok := tokenValue.(*jwt.Token); ok {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if idStr, ok := claims["id"].(string); ok {
						if userID, err := uuid.Parse(idStr); err == nil {
							ctx = setUserIDToContext(ctx, userID)
							c.SetRequest(c.Request().WithContext(ctx))
						}
					}
				}
			}
		}

		srv.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

func blockIntrospection(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	opCtx := graphql.GetOperationContext(ctx)
	if opCtx != nil {
		query := strings.ToLower(opCtx.RawQuery)
		if strings.Contains(query, "__schema") ||
			strings.Contains(query, "__type") ||
			opCtx.OperationName == "IntrospectionQuery" {
			return func(ctx context.Context) *graphql.Response {
				return &graphql.Response{
					Errors: graphql.GetErrors(ctx),
				}
			}
		}
	}
	return next(ctx)
}

func SchemaHandler(isProduction bool) echo.HandlerFunc {
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
	return func(c echo.Context) error {
		if isProduction {
			return c.String(404, "Not Found")
		}

		c.Response().Header().Set("Content-Type", "text/plain; charset=utf-8")
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request().Method == "OPTIONS" {
			return c.NoContent(204)
		}

		return c.String(200, schema)
	}
}
