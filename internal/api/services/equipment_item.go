package services

import (
	"context"

	"moonshine/internal/domain"
	"moonshine/internal/repository"
)

type EquipmentItemService struct {
	equipmentItemRepo *repository.EquipmentItemRepository
}

func NewEquipmentItemService(equipmentItemRepo *repository.EquipmentItemRepository) *EquipmentItemService {
	return &EquipmentItemService{
		equipmentItemRepo: equipmentItemRepo,
	}
}

func (s *EquipmentItemService) GetByCategorySlug(ctx context.Context, slug string) ([]*domain.EquipmentItem, error) {
	return s.equipmentItemRepo.FindByCategorySlug(slug)
}

