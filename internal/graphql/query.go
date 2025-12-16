package graphql

import (
	"context"
	"errors"

	"moonshine/internal/graphql/models"
)

type queryResolver struct{ *Resolver }

func (r *queryResolver) CurrentUser(ctx context.Context) (*models.User, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, errors.New("unauthorized: invalid or missing token")
	}

	user, err := r.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return domainUserToGraphQL(user), nil
}
