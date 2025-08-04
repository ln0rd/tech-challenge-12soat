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

// FetchInputFromDB busca um input específico do banco de dados
func (uc *UpdateByIdInput) FetchInputFromDB(id uuid.UUID) (*models.Input, error) {
	var existingInput models.Input
	if err := uc.DB.Where("id = ?", id).First(&existingInput).Error; err != nil {
		uc.Logger.Error("Database error finding input to update", zap.Error(err), zap.String("id", id.String()))
		return nil, err
	}

	uc.Logger.Info("Found existing input",
		zap.String("id", existingInput.ID.String()),
		zap.String("name", existingInput.Name),
		zap.String("inputType", existingInput.InputType))

	return &existingInput, nil
}

// ValidateInputNameUniqueness verifica se o nome do input é único (para update)
func (uc *UpdateByIdInput) ValidateInputNameUniqueness(name string, inputID uuid.UUID) error {
	var inputWithSameName models.Input
	if err := uc.DB.Where("name = ? AND id != ?", name, inputID).First(&inputWithSameName).Error; err == nil {
		uc.Logger.Error("Input name already exists", zap.String("name", name))
		return errors.New("input name already exists")
	} else if err != gorm.ErrRecordNotFound {
		uc.Logger.Error("Error checking input name uniqueness", zap.Error(err))
		return err
	}

	uc.Logger.Info("Input name is unique", zap.String("name", name))
	return nil
}

// AdjustQuantityForInputType ajusta a quantidade baseado no tipo de input
func (uc *UpdateByIdInput) AdjustQuantityForInputType(quantity int, inputType string) int {
	if inputType == "service" {
		uc.Logger.Info("Forcing quantity to 1 for service type",
			zap.String("inputType", inputType),
			zap.Int("originalQuantity", quantity),
			zap.Int("finalQuantity", 1))
		return 1
	}
	return quantity
}

// UpdateInputFields atualiza os campos do input existente
func (uc *UpdateByIdInput) UpdateInputFields(existingInput *models.Input, entity *domain.Input) {
	finalQuantity := uc.AdjustQuantityForInputType(entity.Quantity, entity.InputType)

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
}

// SaveInputToDB salva as alterações do input no banco de dados
func (uc *UpdateByIdInput) SaveInputToDB(input *models.Input) error {
	result := uc.DB.Save(input)
	if result.Error != nil {
		uc.Logger.Error("Database error updating input", zap.Error(result.Error))
		return result.Error
	}

	uc.Logger.Info("Input updated successfully",
		zap.String("id", input.ID.String()),
		zap.Int64("rowsAffected", result.RowsAffected))
	return nil
}

func (uc *UpdateByIdInput) Process(id uuid.UUID, entity *domain.Input) error {
	uc.Logger.Info("Processing update input by ID",
		zap.String("id", id.String()),
		zap.String("name", entity.Name),
		zap.String("inputType", entity.InputType))

	// Busca o input existente
	existingInput, err := uc.FetchInputFromDB(id)
	if err != nil {
		return err
	}

	// Verifica se o novo nome já existe (se foi alterado)
	if entity.Name != existingInput.Name {
		if err := uc.ValidateInputNameUniqueness(entity.Name, id); err != nil {
			return err
		}
	}

	// Atualiza os campos do input
	uc.UpdateInputFields(existingInput, entity)

	// Salva as alterações
	err = uc.SaveInputToDB(existingInput)
	if err != nil {
		return err
	}

	return nil
}
