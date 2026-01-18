package dto

import (
	"time"

	"moonshine/internal/domain"
)

type EquipmentItem struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Attack        int       `json:"attack"`
	Defense       int       `json:"defense"`
	Hp            int       `json:"hp"`
	RequiredLevel int       `json:"requiredLevel"`
	Price         int       `json:"price"`
	Artifact      bool      `json:"artifact"`
	Image         string    `json:"image"`
	CreatedAt     time.Time `json:"createdAt"`
}

func EquipmentItemFromDomain(item *domain.EquipmentItem) *EquipmentItem {
	if item == nil {
		return nil
	}

	return &EquipmentItem{
		ID:            item.ID.String(),
		Name:          item.Name,
		Slug:          item.Slug,
		Attack:        int(item.Attack),
		Defense:       int(item.Defense),
		Hp:            int(item.Hp),
		RequiredLevel: int(item.RequiredLevel),
		Price:         int(item.Price),
		Artifact:      item.Artifact,
		Image:         item.Image,
		CreatedAt:     item.CreatedAt,
	}
}

func EquipmentItemsFromDomain(items []*domain.EquipmentItem) []*EquipmentItem {
	result := make([]*EquipmentItem, len(items))
	for i, item := range items {
		result[i] = EquipmentItemFromDomain(item)
	}
	return result
}
