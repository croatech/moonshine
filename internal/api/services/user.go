package services

import (
	"context"
	"log"

	"github.com/google/uuid"

	"moonshine/internal/domain"
	"moonshine/internal/repository"
)

type UserService struct {
	userRepo   *repository.UserRepository
	avatarRepo *repository.AvatarRepository
}

func NewUserService(userRepo *repository.UserRepository, avatarRepo *repository.AvatarRepository) *UserService {
	return &UserService{
		userRepo:   userRepo,
		avatarRepo: avatarRepo,
	}
}

func (s *UserService) GetCurrentUser(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	log.Printf("[UserService] Fetching user with ID: %s", userID)
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		log.Printf("[UserService] Failed to find user by ID %s: %v", userID, err)
		return nil, repository.ErrUserNotFound
	}

	// Load avatar if avatar_id is present
	if user.AvatarID != nil {
		log.Printf("[UserService] Loading avatar with ID: %s", *user.AvatarID)
		avatar, err := s.avatarRepo.FindByID(*user.AvatarID)
		if err != nil {
			log.Printf("[UserService] Failed to load avatar: %v", err)
		} else if avatar != nil {
			user.Avatar = avatar
			log.Printf("[UserService] Avatar loaded: ID=%s, Image=%s", avatar.ID, avatar.Image)
		}
	} else {
		log.Printf("[UserService] User has no avatar_id")
	}

	log.Printf("[UserService] Successfully fetched user: %s, Avatar=%v", user.Username, user.Avatar != nil)
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID uuid.UUID, avatarID *uuid.UUID) (*domain.User, error) {
	log.Printf("[UserService] Updating user %s, avatar_id: %v", userID, avatarID)
	
	// Validate avatar exists if provided
	if avatarID != nil {
		_, err := s.avatarRepo.FindByID(*avatarID)
		if err != nil {
			log.Printf("[UserService] Avatar not found: %v", err)
			return nil, repository.ErrAvatarNotFound
		}
	}
	
	// Update avatar_id
	err := s.userRepo.UpdateAvatarID(userID, avatarID)
	if err != nil {
		log.Printf("[UserService] Failed to update user avatar_id: %v", err)
		return nil, err
	}
	
	// Return updated user
	return s.GetCurrentUser(ctx, userID)
}
