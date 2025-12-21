package graphql

import (
	"context"

	"moonshine/internal/graphql/models"
	"moonshine/internal/graphql/queries"
)

// CurrentUser returns the currently authenticated user
func (r *queryResolver) CurrentUser(ctx context.Context) (*models.User, error) {
	return queries.CurrentUser(ctx, r.UserRepo)
}
