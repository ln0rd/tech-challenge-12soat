package customer

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/costumer"
	interfaces "github.com/ln0rd/tech_challenge_12soat/internal/domain/interfaces"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/repository"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"go.uber.org/zap"
)

type FindAllCustomer struct {
	CustomerRepository repository.CustomerRepository
	Logger             interfaces.Logger
}

// FetchCustomersFromDB busca todos os customers do banco
func (uc *FindAllCustomer) FetchCustomersFromDB() ([]models.Customer, error) {
	customers, err := uc.CustomerRepository.FindAll()
	if err != nil {
		uc.Logger.Error("Database error fetching customers", zap.Error(err))
		return nil, err
	}

	uc.Logger.Info("Successfully fetched customers from database", zap.Int("count", len(customers)))
	return customers, nil
}

func (uc *FindAllCustomer) Process() ([]domain.Customer, error) {
	uc.Logger.Info("Processing find all customers")

	// Busca customers do banco
	customers, err := uc.FetchCustomersFromDB()
	if err != nil {
		return nil, err
	}

	// Mapeia para o domínio usando persistence
	var domainCustomers []domain.Customer
	for _, customer := range customers {
		domainCustomer := persistence.CustomerPersistence{}.ToEntity(&customer)
		domainCustomers = append(domainCustomers, *domainCustomer)
	}

	uc.Logger.Info("Successfully mapped customers to domain", zap.Int("count", len(domainCustomers)))

	return domainCustomers, nil
}
