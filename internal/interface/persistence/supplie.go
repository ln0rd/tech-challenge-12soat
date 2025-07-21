package persistence

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/supplie"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
)

type SuppliePersistence struct{}

func (SuppliePersistence) ToEntity(model *models.Supplie) *domain.Supplie {
	if model == nil {
		return nil
	}
	return &domain.Supplie{
		ID:                model.ID,
		Name:              model.Name,
		Description:       model.Description,
		Price:             model.Price,
		QuantityAvailable: model.QuantityAvailable,
		CreatedAt:         model.CreatedAt,
		UpdatedAt:         model.UpdatedAt,
	}
}

func (SuppliePersistence) ToModel(entity *domain.Supplie) *models.Supplie {
	if entity == nil {
		return nil
	}
	return &models.Supplie{
		ID:                entity.ID,
		Name:              entity.Name,
		Description:       entity.Description,
		Price:             entity.Price,
		QuantityAvailable: entity.QuantityAvailable,
		CreatedAt:         entity.CreatedAt,
		UpdatedAt:         entity.UpdatedAt,
	}
}
