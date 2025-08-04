package customer

import (
	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/costumer"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UpdateByIdCustomer struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

// FetchCustomerFromDB busca um customer específico do banco de dados
func (uc *UpdateByIdCustomer) FetchCustomerFromDB(id uuid.UUID) (*models.Customer, error) {
	var existingCustomer models.Customer
	if err := uc.DB.Where("id = ?", id).First(&existingCustomer).Error; err != nil {
		uc.Logger.Error("Database error finding customer to update", zap.Error(err), zap.String("id", id.String()))
		return nil, err
	}

	uc.Logger.Info("Found existing customer", zap.String("id", existingCustomer.ID.String()), zap.String("name", existingCustomer.Name))
	return &existingCustomer, nil
}

// UpdateCustomerFields atualiza os campos do customer existente
func (uc *UpdateByIdCustomer) UpdateCustomerFields(existingCustomer *models.Customer, entity *domain.Customer) {
	existingCustomer.Name = entity.Name
	existingCustomer.DocumentNumber = entity.DocumentNumber
	existingCustomer.CustomerType = entity.CustomerType

	uc.Logger.Info("Updated customer fields", zap.String("name", existingCustomer.Name), zap.String("documentNumber", existingCustomer.DocumentNumber))
}

// SaveCustomerToDB salva as alterações do customer no banco de dados
func (uc *UpdateByIdCustomer) SaveCustomerToDB(customer *models.Customer) error {
	result := uc.DB.Save(customer)
	if result.Error != nil {
		uc.Logger.Error("Database error updating customer", zap.Error(result.Error))
		return result.Error
	}

	uc.Logger.Info("Customer updated successfully", zap.String("id", customer.ID.String()), zap.Int64("rowsAffected", result.RowsAffected))
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
