package customer

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/costumer"
	interfaces "github.com/ln0rd/tech_challenge_12soat/internal/domain/interfaces"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/repository"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"go.uber.org/zap"
)

type CreateCustomer struct {
	CustomerRepository repository.CustomerRepository
	Logger             interfaces.Logger
}

// SaveCustomerToDB salva o customer no banco de dados
func (uc *CreateCustomer) SaveCustomerToDB(model *models.Customer) error {
	err := uc.CustomerRepository.Create(model)
	if err != nil {
		uc.Logger.Error("Database error creating customer", zap.Error(err))
		return err
	}

	uc.Logger.Info("Customer created in database", zap.String("name", model.Name))
	return nil
}

func (uc *CreateCustomer) Process(entity *domain.Customer) error {
	uc.Logger.Info("Processing customer creation", zap.String("name", entity.Name))

	// Mapeia entidade para modelo usando persistence
	model := persistence.CustomerPersistence{}.ToModel(entity)
	uc.Logger.Info("Model created", zap.String("name", model.Name), zap.String("documentNumber", model.DocumentNumber))

	// Salva no banco
	err := uc.SaveCustomerToDB(model)
	if err != nil {
		return err
	}

	return nil
}
