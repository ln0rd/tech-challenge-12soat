package input

import (
	"errors"

	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/input"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CreateInput struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (uc *CreateInput) Process(entity *domain.Input) error {
	uc.Logger.Info("Processing input creation",
		zap.String("name", entity.Name),
		zap.Float64("price", entity.Price),
		zap.Int("quantity", entity.Quantity))

	// Verifica se o nome já existe
	var existingInput models.Input
	if err := uc.DB.Where("name = ?", entity.Name).First(&existingInput).Error; err == nil {
		uc.Logger.Error("Input name already exists", zap.String("name", entity.Name))
		return errors.New("input name already exists")
	} else if err != gorm.ErrRecordNotFound {
		uc.Logger.Error("Error checking input name uniqueness", zap.Error(err))
		return err
	}

	uc.Logger.Info("Input name is unique", zap.String("name", entity.Name))

	model := &models.Input{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		Price:       entity.Price,
		Quantity:    entity.Quantity,
		InputType:   entity.InputType,
	}

	uc.Logger.Info("Model created",
		zap.String("name", model.Name),
		zap.Float64("price", model.Price),
		zap.Int("quantity", model.Quantity))

	result := uc.DB.Create(model)
	if result.Error != nil {
		uc.Logger.Error("Database error creating input", zap.Error(result.Error))
		return result.Error
	}

	uc.Logger.Info("Input created in database",
		zap.String("id", model.ID.String()),
		zap.String("name", model.Name),
		zap.Int64("rowsAffected", result.RowsAffected))

	return nil
}
