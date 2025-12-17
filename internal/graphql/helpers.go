package graphql

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"moonshine/internal/domain"
	"moonshine/internal/graphql/models"
)

func domainUserToGraphQL(user *domain.User) *models.User {
	return &models.User{
		ID:        formatID(user.ID),
		Username:  user.Username,
		Email:     user.Email,
		Hp:        int(user.Hp),
		Level:     int(user.Level),
		Gold:      int(user.Gold),
		Exp:       int(user.Exp),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func formatID(id uuid.UUID) string {
	return id.String()
}

func generateJWTToken(id uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"id":  id.String(),
		"exp": time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_KEY")))
}

func getUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userIDValue := ctx.Value("userID")
	if userIDValue == nil {
		return uuid.Nil, errors.New("unauthorized")
	}

	var userID uuid.UUID
	switch v := userIDValue.(type) {
	case uuid.UUID:
		userID = v
	case string:
		var err error
		userID, err = uuid.Parse(v)
		if err != nil {
			return uuid.Nil, errors.New("unauthorized: invalid user ID")
		}
	default:
		return uuid.Nil, errors.New("unauthorized: invalid user ID type")
	}

	return userID, nil
}

