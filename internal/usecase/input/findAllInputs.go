package input

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/input"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FindAllInputs struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (uc *FindAllInputs) Process() ([]domain.Input, error) {
	uc.Logger.Info("Processing find all inputs")

	var inputs []models.Input
	if err := uc.DB.Find(&inputs).Error; err != nil {
		uc.Logger.Error("Database error finding all inputs", zap.Error(err))
		return []domain.Input{}, err
	}

	uc.Logger.Info("Found inputs in database", zap.Int("count", len(inputs)))

	// Inicializa uma lista vazia
	domainInputs := make([]domain.Input, 0)

	for _, input := range inputs {
		domainInputs = append(domainInputs, domain.Input{
			ID:          input.ID,
			Name:        input.Name,
			Description: input.Description,
			Price:       input.Price,
			Quantity:    input.Quantity,
			InputType:   input.InputType,
			CreatedAt:   input.CreatedAt,
			UpdatedAt:   input.UpdatedAt,
		})
	}

	uc.Logger.Info("Successfully mapped inputs to domain", zap.Int("count", len(domainInputs)))

	// Sempre retorna uma lista (vazia se não encontrou inputs)
	return domainInputs, nil
}
