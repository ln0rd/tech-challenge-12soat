package customer

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/costumer"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CreateCustomer struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

// SaveCustomerToDB salva o customer no banco de dados
func (uc *CreateCustomer) SaveCustomerToDB(model *models.Customer) error {
	result := uc.DB.Create(model)
	if result.Error != nil {
		uc.Logger.Error("Database error creating customer", zap.Error(result.Error))
		return result.Error
	}

	uc.Logger.Info("Customer created in database", zap.String("name", model.Name), zap.Int64("rowsAffected", result.RowsAffected))
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
