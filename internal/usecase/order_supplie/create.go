package order_supplie

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/order_supplie"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"gorm.io/gorm"
)

type CreateOrderSupplie struct {
	DB *gorm.DB
}

func (uc *CreateOrderSupplie) Process(entity *domain.OrderSupplie) error {
	model := persistence.OrderSuppliesPersistence{}.ToModel(entity)
	return uc.DB.Create(model).Error
}
