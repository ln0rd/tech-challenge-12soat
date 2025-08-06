package input

import (
	"github.com/google/uuid"
	interfaces "github.com/ln0rd/tech_challenge_12soat/internal/domain/interfaces"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/repository"
	"go.uber.org/zap"
)

type DeleteByIdInput struct {
	InputRepository repository.InputRepository
	Logger          interfaces.Logger
}

// DeleteInputFromDB remove o input do banco de dados
func (uc *DeleteByIdInput) DeleteInputFromDB(id uuid.UUID) error {
	err := uc.InputRepository.Delete(id)
	if err != nil {
		uc.Logger.Error("Database error deleting input", zap.Error(err), zap.String("id", id.String()))
		return err
	}

	uc.Logger.Info("Input deleted successfully", zap.String("id", id.String()))
	return nil
}

func (uc *DeleteByIdInput) Process(id uuid.UUID) error {
	uc.Logger.Info("Processing delete input by ID", zap.String("id", id.String()))

	// Remove o input do banco
	err := uc.DeleteInputFromDB(id)
	if err != nil {
		return err
	}

	return nil
}
