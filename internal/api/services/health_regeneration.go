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

func (s *HealthRegenerationService) RegenerateAllUsers(percent float64) ([]repository.HPUpdate, error) {
	return s.userRepo.RegenerateAllUsersHealth(percent)
}
