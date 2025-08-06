package input

import (
	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/input"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/logger"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/repository"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"go.uber.org/zap"
)

type FindByIdInput struct {
	InputRepository repository.InputRepository
	Logger          logger.Logger
}

// FetchInputFromDB busca um input específico do banco de dados
func (uc *FindByIdInput) FetchInputFromDB(id uuid.UUID) (*models.Input, error) {
	input, err := uc.InputRepository.FindByID(id)
	if err != nil {
		uc.Logger.Error("Database error finding input by ID", zap.Error(err), zap.String("id", id.String()))
		return nil, err
	}

	uc.Logger.Info("Found input in database",
		zap.String("id", input.ID.String()),
		zap.String("name", input.Name),
		zap.String("inputType", input.InputType),
		zap.Float64("price", input.Price),
		zap.Int("quantity", input.Quantity))

	return input, nil
}

func (uc *FindByIdInput) Process(id uuid.UUID) (*domain.Input, error) {
	uc.Logger.Info("Processing find input by ID", zap.String("id", id.String()))

	// Busca input do banco
	input, err := uc.FetchInputFromDB(id)
	if err != nil {
		return nil, err
	}

	// Mapeia para o domínio usando persistence
	domainInput := persistence.InputPersistence{}.ToEntity(input)
	uc.Logger.Info("Successfully mapped input to domain", zap.String("id", domainInput.ID.String()))

	return domainInput, nil
}
