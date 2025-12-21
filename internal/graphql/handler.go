package graphql

import (
	"context"
	"reflect"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"moonshine/internal/graphql/generated"
	"moonshine/internal/repository"
)

func GraphQLHandler(db *sqlx.DB, isProduction bool) echo.HandlerFunc {
	userRepo := repository.NewUserRepository(db)
	avatarRepo := repository.NewAvatarRepository(db)
	locationRepo := repository.NewLocationRepository(db)
	resolver := newResolver(userRepo, avatarRepo, locationRepo)

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{
			Resolvers: resolver,
		}),
	)

	// Pass context with userID to GraphQL operations
	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		if isProduction {
			return blockIntrospection(ctx, next)
		}
		return next(ctx)
	})

	return func(c echo.Context) error {
		ctx := c.Request().Context()

		// Extract token from Echo context (set by echo-jwt middleware)
		tokenValue := c.Get("user")

		if tokenValue != nil {
			// Use reflection to extract claims since direct type assertion doesn't work
			rv := reflect.ValueOf(tokenValue)
			if rv.Kind() == reflect.Ptr {
				rv = rv.Elem()
			}

			// Try to find Claims field
			claimsField := rv.FieldByName("Claims")
			if claimsField.IsValid() && claimsField.CanInterface() {
				claims := claimsField.Interface()

				// Work with claims as a map using reflection to access values
				claimsValue := reflect.ValueOf(claims)
				if claimsValue.Kind() == reflect.Map {
					idValue := claimsValue.MapIndex(reflect.ValueOf("id"))
					if idValue.IsValid() && idValue.CanInterface() {
						idInterface := idValue.Interface()
						if idStr, ok := idInterface.(string); ok {
							if userID, err := uuid.Parse(idStr); err == nil {
								ctx = SetUserIDToContext(ctx, userID)
							}
						}
					}
				}
			} else {
				// Try via interface
				if tokenInterface, ok := tokenValue.(interface{ Claims() jwt.Claims }); ok {
					claims := tokenInterface.Claims()

					if mapClaims, ok := claims.(jwt.MapClaims); ok {
						if idStr, ok := mapClaims["id"].(string); ok {
							if userID, err := uuid.Parse(idStr); err == nil {
								ctx = SetUserIDToContext(ctx, userID)
							}
						}
					}
				} else {
					// Try direct type assertion
					if t, ok := tokenValue.(*jwt.Token); ok {
						if claims, ok := t.Claims.(jwt.MapClaims); ok {
							if idStr, ok := claims["id"].(string); ok {
								if userID, err := uuid.Parse(idStr); err == nil {
									ctx = SetUserIDToContext(ctx, userID)
								}
							}
						}
					}
				}
			}
		}

		// Create new request with updated context for GraphQL
		req := c.Request().WithContext(ctx)
		srv.ServeHTTP(c.Response(), req)
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
