package input

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/input"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/logger"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/repository"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"go.uber.org/zap"
)

type FindAllInputs struct {
	InputRepository repository.InputRepository
	Logger          logger.Logger
}

// FetchInputsFromDB busca todos os inputs do banco de dados
func (uc *FindAllInputs) FetchInputsFromDB() ([]models.Input, error) {
	inputs, err := uc.InputRepository.FindAll()
	if err != nil {
		uc.Logger.Error("Database error finding all inputs", zap.Error(err))
		return []models.Input{}, err
	}

	uc.Logger.Info("Found inputs in database", zap.Int("count", len(inputs)))
	return inputs, nil
}

func (uc *FindAllInputs) Process() ([]domain.Input, error) {
	uc.Logger.Info("Processing find all inputs")

	// Busca inputs do banco
	inputs, err := uc.FetchInputsFromDB()
	if err != nil {
		return []domain.Input{}, err
	}

	// Mapeia para o domínio usando persistence
	domainInputs := make([]domain.Input, 0)
	for _, input := range inputs {
		domainInput := persistence.InputPersistence{}.ToEntity(&input)
		domainInputs = append(domainInputs, *domainInput)
	}

	uc.Logger.Info("Successfully mapped inputs to domain", zap.Int("count", len(domainInputs)))

	// Sempre retorna uma lista (vazia se não encontrou inputs)
	return domainInputs, nil
}
