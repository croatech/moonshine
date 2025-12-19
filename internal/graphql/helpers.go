package graphql

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"moonshine/internal/domain"
	"moonshine/internal/graphql/models"
)

func domainUserToGraphQL(user *domain.User) *models.User {
	return &models.User{
		ID:        user.ID.String(),
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

func generateJWTToken(id uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"id":  id.String(),
		"exp": time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_KEY")))
}
