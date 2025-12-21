package queries

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"moonshine/internal/domain"
	"moonshine/internal/graphql/models"
)

type contextKey string

const userIDKey contextKey = "userID"

var errUnauthorized = errors.New("unauthorized")

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	v := ctx.Value(userIDKey)
	if v == nil {
		return uuid.Nil, errUnauthorized
	}

	switch id := v.(type) {
	case uuid.UUID:
		return id, nil
	case string:
		parsed, err := uuid.Parse(id)
		if err != nil {
			return uuid.Nil, errUnauthorized
		}
		return parsed, nil
	default:
		return uuid.Nil, errUnauthorized
	}
}

// DomainUserToGraphQL converts domain user to GraphQL model
func DomainUserToGraphQL(user *domain.User) *models.User {
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

