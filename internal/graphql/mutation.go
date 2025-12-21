package graphql

import (
	"context"

	"moonshine/internal/graphql/models"
	"moonshine/internal/graphql/mutations"
)

// SignUp creates a new user account
func (r *mutationResolver) SignUp(ctx context.Context, input models.SignUpInput) (*models.AuthPayload, error) {
	return mutations.SignUp(ctx, r.UserRepo, input)
}

// SignIn authenticates an existing user
func (r *mutationResolver) SignIn(ctx context.Context, input models.SignInInput) (*models.AuthPayload, error) {
	return mutations.SignIn(ctx, r.UserRepo, input)
}


