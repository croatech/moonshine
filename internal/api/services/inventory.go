package services

import (
	"context"

	"github.com/google/uuid"

	"moonshine/internal/domain"
	"moonshine/internal/repository"
)

type InventoryService struct {
	userEquipmentRepo *repository.UserEquipmentItemRepository
}

func NewInventoryService(userEquipmentRepo *repository.UserEquipmentItemRepository) *InventoryService {
	return &InventoryService{
		userEquipmentRepo: userEquipmentRepo,
	}
}

func (s *InventoryService) GetUserInventory(ctx context.Context, userID uuid.UUID) ([]*domain.EquipmentItem, error) {
	return s.userEquipmentRepo.FindByUserID(userID)
}

