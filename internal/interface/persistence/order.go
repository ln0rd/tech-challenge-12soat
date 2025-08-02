package persistence

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/order"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
)

type OrderPersistence struct{}

func (OrderPersistence) ToEntity(model *models.Order) *domain.Order {
	if model == nil {
		return nil
	}
	return &domain.Order{
		ID:         model.ID,
		CustomerID: model.CustomerID,
		VehicleID:  model.VehicleID,
		Status:     model.Status,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}
}

func (OrderPersistence) ToModel(entity *domain.Order) *models.Order {
	if entity == nil {
		return nil
	}
	return &models.Order{
		ID:         entity.ID,
		CustomerID: entity.CustomerID,
		VehicleID:  entity.VehicleID,
		Status:     entity.Status,
		CreatedAt:  entity.CreatedAt,
		UpdatedAt:  entity.UpdatedAt,
	}
}
