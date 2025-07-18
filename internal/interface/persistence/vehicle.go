package persistence

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/vehicle"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
)

type VehiclePersistence struct{}

func (VehiclePersistence) ToEntity(model *models.Vehicle) *domain.Vehicle {
	if model == nil {
		return nil
	}
	return &domain.Vehicle{
		ID:                          model.ID,
		Model:                       model.Model,
		Brand:                       model.Brand,
		ReleaseYear:                 model.ReleaseYear,
		VehicleIdentificationNumber: model.VehicleIdentificationNumber,
		NumberPlate:                 model.NumberPlate,
		Color:                       model.Color,
		CustomerID:                  model.CustomerID,
		CreatedAt:                   model.CreatedAt,
		UpdatedAt:                   model.UpdatedAt,
	}
}

func (VehiclePersistence) ToModel(entity *domain.Vehicle) *models.Vehicle {
	if entity == nil {
		return nil
	}
	return &models.Vehicle{
		ID:                          entity.ID,
		Model:                       entity.Model,
		Brand:                       entity.Brand,
		ReleaseYear:                 entity.ReleaseYear,
		VehicleIdentificationNumber: entity.VehicleIdentificationNumber,
		NumberPlate:                 entity.NumberPlate,
		Color:                       entity.Color,
		CustomerID:                  entity.CustomerID,
		CreatedAt:                   entity.CreatedAt,
		UpdatedAt:                   entity.UpdatedAt,
	}
}
