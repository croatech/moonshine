package services

import (
	"log"

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
	rowsAffected, err := s.userRepo.RegenerateAllUsersHealth(percent)
	if err != nil {
		log.Printf("[HealthRegenerationService] Failed to regenerate health: %v", err)
		return err
	}
	log.Printf("[HealthRegenerationService] Regenerated health for %d users", rowsAffected)
	return nil
}
