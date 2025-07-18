package persistence

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/user"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
)

type UserPersistence struct{}

func (UserPersistence) ToEntity(model *models.User) *domain.User {
	if model == nil {
		return nil
	}
	return &domain.User{
		ID:         model.ID,
		Email:      model.Email,
		Password:   model.Password,
		Username:   model.Username,
		CustomerID: model.CustomerID,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}
}

func (UserPersistence) ToModel(entity *domain.User) *models.User {
	if entity == nil {
		return nil
	}
	return &models.User{
		ID:         entity.ID,
		Email:      entity.Email,
		Password:   entity.Password,
		Username:   entity.Username,
		CustomerID: entity.CustomerID,
		CreatedAt:  entity.CreatedAt,
		UpdatedAt:  entity.UpdatedAt,
	}
}
