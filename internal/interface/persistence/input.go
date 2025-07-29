package persistence

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/input"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
)

type InputPersistence struct{}

func (InputPersistence) ToEntity(model *models.Input) *domain.Input {
	if model == nil {
		return nil
	}
	return &domain.Input{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		Price:       model.Price,
		Quantity:    model.Quantity,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func (InputPersistence) ToModel(entity *domain.Input) *models.Input {
	if entity == nil {
		return nil
	}
	return &models.Input{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		Price:       entity.Price,
		Quantity:    entity.Quantity,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
