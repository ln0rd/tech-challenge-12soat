package customer

import (
	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/costumer"
	interfaces "github.com/ln0rd/tech_challenge_12soat/internal/domain/interfaces"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/repository"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FindByIdCustomer struct {
	CustomerRepository repository.CustomerRepository
	Logger             interfaces.Logger
}

// FetchCustomerFromDB busca um customer específico do banco
func (uc *FindByIdCustomer) FetchCustomerFromDB(id uuid.UUID) (*models.Customer, error) {
	customer, err := uc.CustomerRepository.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			uc.Logger.Error("Customer not found", zap.String("id", id.String()))
			return nil, err
		}
		uc.Logger.Error("Database error fetching customer", zap.Error(err))
		return nil, err
	}

	uc.Logger.Info("Successfully fetched customer from database", zap.String("id", customer.ID.String()))
	return customer, nil
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
