package customer

import (
	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/costumer"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FindByIdCustomer struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

// FetchCustomerFromDB busca um customer específico do banco de dados
func (uc *FindByIdCustomer) FetchCustomerFromDB(id uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	if err := uc.DB.Where("id = ?", id).First(&customer).Error; err != nil {
		uc.Logger.Error("Database error finding customer by ID", zap.Error(err), zap.String("id", id.String()))
		return nil, err
	}

	uc.Logger.Info("Found customer in database", zap.String("id", customer.ID.String()), zap.String("name", customer.Name))
	return &customer, nil
}

func (uc *FindByIdCustomer) Process(id uuid.UUID) (*domain.Customer, error) {
	uc.Logger.Info("Processing find customer by ID", zap.String("id", id.String()))

	// Busca customer do banco
	customer, err := uc.FetchCustomerFromDB(id)
	if err != nil {
		return nil, err
	}

	// Mapeia para o domínio usando persistence
	domainCustomer := persistence.CustomerPersistence{}.ToEntity(customer)
	uc.Logger.Info("Successfully mapped customer to domain", zap.String("id", domainCustomer.ID.String()))

	return domainCustomer, nil
}
