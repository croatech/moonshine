package services

import (
	"moonshine/internal/repository"
)

type HealthRegenerationService struct {
	userRepo *repository.UserRepository
}

func NewHealthRegenerationService(userRepo *repository.UserRepository) *HealthRegenerationService {
	return &HealthRegenerationService{
		userRepo: userRepo,
	}
}

func (s *HealthRegenerationService) RegenerateAllUsers(percent float64) error {
	_, err := s.userRepo.RegenerateAllUsersHealth(percent)
	if err != nil {
		return err
	}
	return nil
}
