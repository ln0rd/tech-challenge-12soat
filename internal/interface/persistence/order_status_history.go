package persistence

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/order_status_history"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
)

type OrderStatusHistoryPersistence struct{}

func (OrderStatusHistoryPersistence) ToEntity(model *models.OrderStatusHistory) *domain.OrderStatusHistory {
	if model == nil {
		return nil
	}
	return &domain.OrderStatusHistory{
		ID:        model.ID,
		OrderID:   model.OrderID,
		Status:    model.Status,
		StartedAt: model.StartedAt,
		EndedAt:   model.EndedAt,
		CreatedAt: model.CreatedAt,
	}
}

func (OrderStatusHistoryPersistence) ToModel(entity *domain.OrderStatusHistory) *models.OrderStatusHistory {
	if entity == nil {
		return nil
	}
	return &models.OrderStatusHistory{
		ID:        entity.ID,
		OrderID:   entity.OrderID,
		Status:    entity.Status,
		StartedAt: entity.StartedAt,
		EndedAt:   entity.EndedAt,
		CreatedAt: entity.CreatedAt,
	}
}
