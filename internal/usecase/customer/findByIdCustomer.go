package customer

import (
	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/costumer"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FindByIdCustomer struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (uc *FindByIdCustomer) Process(id uuid.UUID) (*domain.Customer, error) {
	uc.Logger.Info("Processing find customer by ID", zap.String("id", id.String()))

	var customer models.Customer
	if err := uc.DB.Where("id = ?", id).First(&customer).Error; err != nil {
		uc.Logger.Error("Database error finding customer by ID", zap.Error(err), zap.String("id", id.String()))
		return nil, err
	}

	uc.Logger.Info("Found customer in database", zap.String("id", customer.ID.String()), zap.String("name", customer.Name))

	domainCustomer := domain.Customer{
		ID:             customer.ID,
		Name:           customer.Name,
		DocumentNumber: customer.DocumentNumber,
		CustomerType:   customer.CustomerType,
		CreatedAt:      customer.CreatedAt,
		UpdatedAt:      customer.UpdatedAt,
	}

	uc.Logger.Info("Successfully mapped customer to domain", zap.String("id", domainCustomer.ID.String()))

	return &domainCustomer, nil
}
