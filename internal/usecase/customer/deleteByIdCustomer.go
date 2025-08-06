package customer

import (
	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/logger"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DeleteByIdCustomer struct {
	CustomerRepository repository.CustomerRepository
	Logger             logger.Logger
}

// DeleteCustomerFromDB remove o customer do banco
func (uc *DeleteByIdCustomer) DeleteCustomerFromDB(id uuid.UUID) error {
	err := uc.CustomerRepository.Delete(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			uc.Logger.Error("Customer not found for deletion", zap.String("id", id.String()))
			return err
		}
		uc.Logger.Error("Database error deleting customer", zap.Error(err), zap.String("id", id.String()))
		return err
	}

	uc.Logger.Info("Customer deleted from database", zap.String("id", id.String()))
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
