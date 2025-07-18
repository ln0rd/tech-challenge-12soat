package customer

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/costumer"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"gorm.io/gorm"
)

type CreateCustomer struct {
	DB *gorm.DB
}

func (uc *CreateCustomer) Process(entity *domain.Customer) error {
	model := persistence.CustomerPersistence{}.ToModel(entity)
	return uc.DB.Create(model).Error
}
