package order_input

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/order_input"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"gorm.io/gorm"
)

type CreateOrderInput struct {
	DB *gorm.DB
}

func (uc *CreateOrderInput) Process(entity *domain.OrderInput) error {
	model := persistence.OrderInputPersistence{}.ToModel(entity)
	return uc.DB.Create(model).Error
}
