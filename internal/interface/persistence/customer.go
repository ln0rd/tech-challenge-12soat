package persistence

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/costumer"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
)

type CustomerPersistence struct{}

func (CustomerPersistence) ToEntity(model *models.Customer) *domain.Customer {
	if model == nil {
		return nil
	}
	return &domain.Customer{
		ID:             model.ID,
		Name:           model.Name,
		DocumentNumber: model.DocumentNumber,
		CustomerType:   model.CustomerType,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
	}
}

func (CustomerPersistence) ToModel(entity *domain.Customer) *models.Customer {
	if entity == nil {
		return nil
	}
	return &models.Customer{
		ID:             entity.ID,
		Name:           entity.Name,
		DocumentNumber: entity.DocumentNumber,
		CustomerType:   entity.CustomerType,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
	}
}
