package customer

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/costumer"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CreateCustomer struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (uc *CreateCustomer) Process(entity *domain.Customer) error {
	uc.Logger.Info("Processing customer creation", zap.String("name", entity.Name))

	model := persistence.CustomerPersistence{}.ToModel(entity)

	uc.Logger.Info("Model created", zap.String("name", model.Name), zap.String("documentNumber", model.DocumentNumber))

	result := uc.DB.Create(model)
	if result.Error != nil {
		uc.Logger.Error("Database error creating customer", zap.Error(result.Error))
		return result.Error
	}

	uc.Logger.Info("Customer created in database", zap.String("name", model.Name), zap.Int64("rowsAffected", result.RowsAffected))

	return nil
}
