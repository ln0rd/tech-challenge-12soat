package customer

import (
	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/costumer"
	interfaces "github.com/ln0rd/tech_challenge_12soat/internal/domain/interfaces"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UpdateByIdCustomer struct {
	CustomerRepository repository.CustomerRepository
	Logger             interfaces.Logger
}

// FetchCustomerFromDB busca um customer específico do banco
func (uc *UpdateByIdCustomer) FetchCustomerFromDB(id uuid.UUID) (*models.Customer, error) {
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

// UpdateCustomerFields atualiza os campos do customer
func (uc *UpdateByIdCustomer) UpdateCustomerFields(existingCustomer *models.Customer, entity *domain.Customer) {
	existingCustomer.Name = entity.Name
	existingCustomer.DocumentNumber = entity.DocumentNumber
	existingCustomer.CustomerType = entity.CustomerType

	uc.Logger.Info("Customer fields updated",
		zap.String("id", existingCustomer.ID.String()),
		zap.String("name", existingCustomer.Name))
}

// SaveCustomerToDB salva as alterações do customer no banco
func (uc *UpdateByIdCustomer) SaveCustomerToDB(customer *models.Customer) error {
	err := uc.CustomerRepository.Update(customer)
	if err != nil {
		uc.Logger.Error("Database error updating customer", zap.Error(err))
		return err
	}

	uc.Logger.Info("Customer updated in database", zap.String("id", customer.ID.String()))
	return nil
}

func (uc *UpdateByIdCustomer) Process(id uuid.UUID, entity *domain.Customer) error {
	uc.Logger.Info("Processing update customer by ID", zap.String("id", id.String()), zap.String("name", entity.Name))

	// Busca o customer existente
	existingCustomer, err := uc.FetchCustomerFromDB(id)
	if err != nil {
		return err
	}

	// Atualiza os campos do customer
	uc.UpdateCustomerFields(existingCustomer, entity)

	// Salva as alterações
	err = uc.SaveCustomerToDB(existingCustomer)
	if err != nil {
		return err
	}

	return nil
}
