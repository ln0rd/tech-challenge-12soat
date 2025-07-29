package customer

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/costumer"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"gorm.io/gorm"
)

type FindAllCustomer struct {
	DB *gorm.DB
}

func (uc *FindAllCustomer) Process() ([]domain.Customer, error) {
	var customers []models.Customer
	if err := uc.DB.Find(&customers).Error; err != nil {
		return nil, err
	}

	var domainCustomers []domain.Customer
	for _, customer := range customers {
		domainCustomers = append(domainCustomers, domain.Customer{
			ID:             customer.ID,
			Name:           customer.Name,
			Email:          customer.Email,
			UserID:         customer.UserID,
			DocumentNumber: customer.DocumentNumber,
			CustomerType:   customer.CustomerType,
			CreatedAt:      customer.CreatedAt,
			UpdatedAt:      customer.UpdatedAt,
		})
	}

	return domainCustomers, nil
}
