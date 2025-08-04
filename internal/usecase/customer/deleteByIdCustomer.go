package customer

import (
	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DeleteByIdCustomer struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

// DeleteCustomerFromDB remove o customer do banco de dados
func (uc *DeleteByIdCustomer) DeleteCustomerFromDB(id uuid.UUID) error {
	result := uc.DB.Where("id = ?", id).Delete(&models.Customer{})
	if result.Error != nil {
		uc.Logger.Error("Database error deleting customer", zap.Error(result.Error), zap.String("id", id.String()))
		return result.Error
	}

	if result.RowsAffected == 0 {
		uc.Logger.Warn("No customer found to delete", zap.String("id", id.String()))
		return gorm.ErrRecordNotFound
	}

	uc.Logger.Info("Customer deleted successfully", zap.String("id", id.String()), zap.Int64("rowsAffected", result.RowsAffected))
	return nil
}

func (uc *DeleteByIdCustomer) Process(id uuid.UUID) error {
	uc.Logger.Info("Processing delete customer by ID", zap.String("id", id.String()))

	// Remove o customer do banco
	err := uc.DeleteCustomerFromDB(id)
	if err != nil {
		return err
	}

	return nil
}
