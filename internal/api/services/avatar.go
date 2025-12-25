package services

import (
	"context"
	"log"

	"moonshine/internal/domain"
	"moonshine/internal/repository"
)

type AvatarService struct {
	avatarRepo *repository.AvatarRepository
}

func NewAvatarService(avatarRepo *repository.AvatarRepository) *AvatarService {
	return &AvatarService{
		avatarRepo: avatarRepo,
	}
}

func (s *AvatarService) GetAllAvatars(ctx context.Context) ([]*domain.Avatar, error) {
	log.Printf("[AvatarService] Fetching all avatars")
	avatars, err := s.avatarRepo.FindAll()
	if err != nil {
		log.Printf("[AvatarService] Failed to fetch avatars: %v", err)
		return nil, err
	}
	log.Printf("[AvatarService] Found %d avatars", len(avatars))
	return avatars, nil
}

