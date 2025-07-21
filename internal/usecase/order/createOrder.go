package order

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/order"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"gorm.io/gorm"
)

type CreateOrder struct {
	DB *gorm.DB
}

func (uc *CreateOrder) Process(entity *domain.Order) error {
	model := persistence.OrderPersistence{}.ToModel(entity)
	return uc.DB.Create(model).Error
}
