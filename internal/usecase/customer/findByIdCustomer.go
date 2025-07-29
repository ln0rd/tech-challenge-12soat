package customer

import (
	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/costumer"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"gorm.io/gorm"
)

type FindByIdCustomer struct {
	DB *gorm.DB
}

func (uc *FindByIdCustomer) Process(id uuid.UUID) (*domain.Customer, error) {
	var customer models.Customer
	if err := uc.DB.Where("id = ?", id).First(&customer).Error; err != nil {
		return nil, err
	}

	domainCustomer := domain.Customer{
		ID:             customer.ID,
		Name:           customer.Name,
		Email:          customer.Email,
		UserID:         customer.UserID,
		DocumentNumber: customer.DocumentNumber,
		CustomerType:   customer.CustomerType,
		CreatedAt:      customer.CreatedAt,
		UpdatedAt:      customer.UpdatedAt,
	}

	return &domainCustomer, nil
}
