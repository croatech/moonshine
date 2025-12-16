package graphql

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"

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

func formatID(id uint) string {
	return fmt.Sprintf("%d", id)
}

func generateJWTToken(id uint) (string, error) {
	claims := jwt.MapClaims{
		"id":  float64(id),
		"exp": time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_KEY")))
}

func getUserIDFromContext(ctx context.Context) (uint, error) {
	userID, ok := ctx.Value("userID").(uint)
	if !ok {
		return 0, errors.New("unauthorized")
	}
	return userID, nil
}

