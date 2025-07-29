package persistence

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/order_input"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
)

type OrderInputPersistence struct{}

func (OrderInputPersistence) ToEntity(model *models.OrderInput) *domain.OrderInput {
	if model == nil {
		return nil
	}
	return &domain.OrderInput{
		ID:         model.ID,
		OrderID:    model.OrderID,
		InputID:    model.InputID,
		Quantity:   model.Quantity,
		UnitPrice:  model.UnitPrice,
		TotalPrice: model.TotalPrice,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}
}

func (OrderInputPersistence) ToModel(entity *domain.OrderInput) *models.OrderInput {
	if entity == nil {
		return nil
	}
	return &models.OrderInput{
		ID:         entity.ID,
		OrderID:    entity.OrderID,
		InputID:    entity.InputID,
		Quantity:   entity.Quantity,
		UnitPrice:  entity.UnitPrice,
		TotalPrice: entity.TotalPrice,
		CreatedAt:  entity.CreatedAt,
		UpdatedAt:  entity.UpdatedAt,
	}
}
