package supplie

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/supplie"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"gorm.io/gorm"
)

type CreateSupplie struct {
	DB *gorm.DB
}

func (uc *CreateSupplie) Process(entity *domain.Supplie) error {
	model := persistence.SuppliePersistence{}.ToModel(entity)
	return uc.DB.Create(model).Error
}
