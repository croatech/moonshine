package queries

import (
	"context"

	"moonshine/internal/graphql/models"
	"moonshine/internal/repository"
)

// CurrentUser returns the currently authenticated user
func CurrentUser(ctx context.Context, userRepo *repository.UserRepository) (*models.User, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := userRepo.FindByID(userID)
	if err != nil {
		return nil, repository.ErrUserNotFound
	}

	return DomainUserToGraphQL(user), nil
}

