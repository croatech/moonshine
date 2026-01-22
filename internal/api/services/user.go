package services

import (
	"context"

	"github.com/google/uuid"

	"moonshine/internal/domain"
	"moonshine/internal/repository"
)

type UserService struct {
	userRepo     *repository.UserRepository
	avatarRepo   *repository.AvatarRepository
	locationRepo *repository.LocationRepository
}

func NewUserService(userRepo *repository.UserRepository, avatarRepo *repository.AvatarRepository, locationRepo *repository.LocationRepository) *UserService {
	return &UserService{
		userRepo:     userRepo,
		avatarRepo:   avatarRepo,
		locationRepo: locationRepo,
	}
}

func (s *UserService) GetCurrentUser(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, repository.ErrUserNotFound
	}

	return user, nil
}

func (s *UserService) GetCurrentUserWithRelations(ctx context.Context, userID uuid.UUID) (*domain.User, *domain.Location, bool, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, nil, false, repository.ErrUserNotFound
	}

	var location *domain.Location
	if s.locationRepo != nil && user.LocationID != uuid.Nil {
		location, _ = s.locationRepo.FindByID(user.LocationID)
	}

	inFight, _ := s.userRepo.InFight(userID)

	return user, location, inFight, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID uuid.UUID, avatarID *uuid.UUID) (*domain.User, error) {

	if avatarID != nil {
		_, err := s.avatarRepo.FindByID(*avatarID)
		if err != nil {
			return nil, repository.ErrAvatarNotFound
		}
	}

	err := s.userRepo.UpdateAvatarID(userID, avatarID)
	if err != nil {
		return nil, err
	}

	return s.GetCurrentUser(ctx, userID)
}
