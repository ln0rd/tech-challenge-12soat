package customer

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/costumer"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FindAllCustomer struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

// FetchCustomersFromDB busca todos os customers do banco de dados
func (uc *FindAllCustomer) FetchCustomersFromDB() ([]models.Customer, error) {
	var customers []models.Customer
	if err := uc.DB.Find(&customers).Error; err != nil {
		uc.Logger.Error("Database error finding all customers", zap.Error(err))
		return nil, err
	}

	uc.Logger.Info("Found customers in database", zap.Int("count", len(customers)))
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
