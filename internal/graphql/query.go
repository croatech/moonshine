package graphql

import (
	"context"

	"moonshine/internal/graphql/models"
	"moonshine/internal/repository"
)

type queryResolver struct{ *Resolver }

func (r *queryResolver) CurrentUser(ctx context.Context) (*models.User, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, errUnauthorized
	}

	user, err := r.userRepo.FindByID(userID)
	if err != nil {
		return nil, repository.ErrUserNotFound
	}

	return domainUserToGraphQL(user), nil
}
