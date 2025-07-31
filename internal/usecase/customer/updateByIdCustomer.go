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

func (uc *UpdateByIdCustomer) Process(id uuid.UUID, entity *domain.Customer) error {
	uc.Logger.Info("Processing update customer by ID", zap.String("id", id.String()), zap.String("name", entity.Name))

	// Primeiro verifica se o customer existe
	var existingCustomer models.Customer
	if err := uc.DB.Where("id = ?", id).First(&existingCustomer).Error; err != nil {
		uc.Logger.Error("Database error finding customer to update", zap.Error(err), zap.String("id", id.String()))
		return err
	}

	uc.Logger.Info("Found existing customer", zap.String("id", existingCustomer.ID.String()), zap.String("name", existingCustomer.Name))

	// Atualiza os campos do customer existente
	existingCustomer.Name = entity.Name
	existingCustomer.DocumentNumber = entity.DocumentNumber
	existingCustomer.CustomerType = entity.CustomerType

	uc.Logger.Info("Updated customer fields", zap.String("name", existingCustomer.Name), zap.String("documentNumber", existingCustomer.DocumentNumber))

	// Salva as alterações
	result := uc.DB.Save(&existingCustomer)
	if result.Error != nil {
		uc.Logger.Error("Database error updating customer", zap.Error(result.Error))
		return result.Error
	}

	uc.Logger.Info("Customer updated successfully", zap.String("id", existingCustomer.ID.String()), zap.Int64("rowsAffected", result.RowsAffected))

	return nil
}
