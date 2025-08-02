package input

import (
	"errors"

	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/input"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UpdateByIdInput struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (uc *UpdateByIdInput) Process(id uuid.UUID, entity *domain.Input) error {
	uc.Logger.Info("Processing update input by ID",
		zap.String("id", id.String()),
		zap.String("name", entity.Name),
		zap.String("inputType", entity.InputType))

	// Primeiro verifica se o input existe
	var existingInput models.Input
	if err := uc.DB.Where("id = ?", id).First(&existingInput).Error; err != nil {
		uc.Logger.Error("Database error finding input to update", zap.Error(err), zap.String("id", id.String()))
		return err
	}

	uc.Logger.Info("Found existing input",
		zap.String("id", existingInput.ID.String()),
		zap.String("name", existingInput.Name),
		zap.String("inputType", existingInput.InputType))

	// Verifica se o novo nome já existe (se foi alterado)
	if entity.Name != existingInput.Name {
		var inputWithSameName models.Input
		if err := uc.DB.Where("name = ? AND id != ?", entity.Name, id).First(&inputWithSameName).Error; err == nil {
			uc.Logger.Error("Input name already exists", zap.String("name", entity.Name))
			return errors.New("input name already exists")
		} else if err != gorm.ErrRecordNotFound {
			uc.Logger.Error("Error checking input name uniqueness", zap.Error(err))
			return err
		}
		uc.Logger.Info("Input name is unique", zap.String("name", entity.Name))
	}

	// Ajusta a quantidade baseado no tipo
	finalQuantity := entity.Quantity
	if entity.InputType == "service" {
		finalQuantity = 1
		uc.Logger.Info("Forcing quantity to 1 for service type",
			zap.String("inputType", entity.InputType),
			zap.Int("originalQuantity", entity.Quantity),
			zap.Int("finalQuantity", finalQuantity))
	}

	// Atualiza os campos do input existente
	existingInput.Name = entity.Name
	existingInput.Description = entity.Description
	existingInput.Price = entity.Price
	existingInput.Quantity = finalQuantity
	existingInput.InputType = entity.InputType

	uc.Logger.Info("Updated input fields",
		zap.String("name", existingInput.Name),
		zap.String("inputType", existingInput.InputType),
		zap.Float64("price", existingInput.Price),
		zap.Int("quantity", existingInput.Quantity))

	// Salva as alterações
	result := uc.DB.Save(&existingInput)
	if result.Error != nil {
		uc.Logger.Error("Database error updating input", zap.Error(result.Error))
		return result.Error
	}

	uc.Logger.Info("Input updated successfully",
		zap.String("id", existingInput.ID.String()),
		zap.Int64("rowsAffected", result.RowsAffected))

	return nil
}
