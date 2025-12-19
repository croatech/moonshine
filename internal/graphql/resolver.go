package graphql

import (
	"moonshine/internal/graphql/generated"
	"moonshine/internal/repository"
)

type Resolver struct {
	userRepo *repository.UserRepository
}

func newResolver(userRepo *repository.UserRepository) *Resolver {
	return &Resolver{
		userRepo: userRepo,
	}
}

func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}
