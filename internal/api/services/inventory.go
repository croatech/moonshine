package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"moonshine/internal/domain"
	"moonshine/internal/repository"
)

type InventoryService struct {
	db *sqlx.DB
}

func NewInventoryService(db *sqlx.DB) *InventoryService {
	return &InventoryService{
		db: db,
	}
}

func (s *InventoryService) GetUserInventory(ctx context.Context, userID uuid.UUID) ([]*domain.EquipmentItem, error) {
	return repository.FindByUserID(s.db, userID)
}

