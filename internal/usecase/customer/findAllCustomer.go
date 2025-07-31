package customer

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/costumer"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FindAllCustomer struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (uc *FindAllCustomer) Process() ([]domain.Customer, error) {
	uc.Logger.Info("Processing find all customers")

	var customers []models.Customer
	if err := uc.DB.Find(&customers).Error; err != nil {
		uc.Logger.Error("Database error finding all customers", zap.Error(err))
		return nil, err
	}

	uc.Logger.Info("Found customers in database", zap.Int("count", len(customers)))

	var domainCustomers []domain.Customer
	for _, customer := range customers {
		domainCustomers = append(domainCustomers, domain.Customer{
			ID:             customer.ID,
			Name:           customer.Name,
			DocumentNumber: customer.DocumentNumber,
			CustomerType:   customer.CustomerType,
			CreatedAt:      customer.CreatedAt,
			UpdatedAt:      customer.UpdatedAt,
		})
	}

	uc.Logger.Info("Successfully mapped customers to domain", zap.Int("count", len(domainCustomers)))

	return domainCustomers, nil
}
