package input

import (
	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/input"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FindByIdInput struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (uc *FindByIdInput) Process(id uuid.UUID) (*domain.Input, error) {
	uc.Logger.Info("Processing find input by ID", zap.String("id", id.String()))

	var input models.Input
	if err := uc.DB.Where("id = ?", id).First(&input).Error; err != nil {
		uc.Logger.Error("Database error finding input by ID", zap.Error(err), zap.String("id", id.String()))
		return nil, err
	}

	uc.Logger.Info("Found input in database",
		zap.String("id", input.ID.String()),
		zap.String("name", input.Name),
		zap.String("inputType", input.InputType),
		zap.Float64("price", input.Price),
		zap.Int("quantity", input.Quantity))

	domainInput := domain.Input{
		ID:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Quantity:    input.Quantity,
		InputType:   input.InputType,
		CreatedAt:   input.CreatedAt,
		UpdatedAt:   input.UpdatedAt,
	}

	uc.Logger.Info("Successfully mapped input to domain", zap.String("id", domainInput.ID.String()))

	return &domainInput, nil
}
