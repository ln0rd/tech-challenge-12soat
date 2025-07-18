package persistence

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/order_supplie"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
)

type OrderSuppliesPersistence struct{}

func (OrderSuppliesPersistence) ToEntity(model *models.OrderSupplie) *domain.OrderSupplie {
	if model == nil {
		return nil
	}
	return &domain.OrderSupplie{
		ID:         model.ID,
		OrderID:    model.OrderID,
		SupplyID:   model.SupplyID,
		Quantity:   model.Quantity,
		TotalValue: model.TotalValue,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}
}

func (OrderSuppliesPersistence) ToModel(entity *domain.OrderSupplie) *models.OrderSupplie {
	if entity == nil {
		return nil
	}
	return &models.OrderSupplie{
		ID:         entity.ID,
		OrderID:    entity.OrderID,
		SupplyID:   entity.SupplyID,
		Quantity:   entity.Quantity,
		TotalValue: entity.TotalValue,
		CreatedAt:  entity.CreatedAt,
		UpdatedAt:  entity.UpdatedAt,
	}
}
