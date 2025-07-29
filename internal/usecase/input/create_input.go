package input

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/input"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"gorm.io/gorm"
)

type CreateInput struct {
	DB *gorm.DB
}

func (uc *CreateInput) Process(entity *domain.Input) error {
	model := persistence.InputPersistence{}.ToModel(entity)
	return uc.DB.Create(model).Error
}
