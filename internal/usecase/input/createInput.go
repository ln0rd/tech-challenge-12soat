package input

import (
	"errors"

	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/input"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/logger"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/repository"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CreateInput struct {
	InputRepository repository.InputRepository
	Logger          logger.Logger
}

// ValidateInputNameUniqueness verifica se o nome do input é único
func (uc *CreateInput) ValidateInputNameUniqueness(name string) error {
	_, err := uc.InputRepository.FindByName(name)
	if err == nil {
		uc.Logger.Error("Input name already exists", zap.String("name", name))
		return errors.New("input name already exists")
	} else if err != gorm.ErrRecordNotFound {
		uc.Logger.Error("Error checking input name uniqueness", zap.Error(err))
		return err
	}

	uc.Logger.Info("Input name is unique", zap.String("name", name))
	return nil
}

// SaveInputToDB salva o input no banco de dados
func (uc *CreateInput) SaveInputToDB(model *models.Input) error {
	err := uc.InputRepository.Create(model)
	if err != nil {
		uc.Logger.Error("Database error creating input", zap.Error(err))
		return err
	}

	uc.Logger.Info("Input created in database",
		zap.String("id", model.ID.String()),
		zap.String("name", model.Name))
	return nil
}

func (uc *CreateInput) Process(entity *domain.Input) error {
	uc.Logger.Info("Processing input creation",
		zap.String("name", entity.Name),
		zap.Float64("price", entity.Price),
		zap.Int("quantity", entity.Quantity))

	// Valida unicidade do nome
	if err := uc.ValidateInputNameUniqueness(entity.Name); err != nil {
		return err
	}

	// Mapeia entidade para modelo usando persistence
	model := persistence.InputPersistence{}.ToModel(entity)
	uc.Logger.Info("Model created",
		zap.String("name", model.Name),
		zap.Float64("price", model.Price),
		zap.Int("quantity", model.Quantity))

	// Salva no banco
	err := uc.SaveInputToDB(model)
	if err != nil {
		return err
	}

	return nil
}
