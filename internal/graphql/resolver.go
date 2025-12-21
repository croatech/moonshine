package graphql

import (
	"moonshine/internal/graphql/generated"
	"moonshine/internal/repository"
)

type Resolver struct {
	UserRepo     *repository.UserRepository
	AvatarRepo   *repository.AvatarRepository
	LocationRepo *repository.LocationRepository
}

type queryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }

func newResolver(userRepo *repository.UserRepository, avatarRepo *repository.AvatarRepository, locationRepo *repository.LocationRepository) *Resolver {
	return &Resolver{
		UserRepo:     userRepo,
		AvatarRepo:   avatarRepo,
		LocationRepo: locationRepo,
	}
}

func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{Resolver: r}
}
