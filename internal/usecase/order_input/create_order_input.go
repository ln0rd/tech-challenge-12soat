package order_input

import (
	interfaces "github.com/ln0rd/tech_challenge_12soat/internal/domain/interfaces"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/order_input"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/repository"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"go.uber.org/zap"
)

type CreateOrderInput struct {
	OrderInputRepository repository.OrderInputRepository
	Logger               interfaces.Logger
}

// SaveOrderInputToDB salva o order input no banco de dados
func (uc *CreateOrderInput) SaveOrderInputToDB(model *models.OrderInput) error {
	err := uc.OrderInputRepository.Create(model)
	if err != nil {
		uc.Logger.Error("Database error creating order input", zap.Error(err))
		return err
	}

	uc.Logger.Info("Order input created in database",
		zap.String("id", model.ID.String()),
		zap.String("orderID", model.OrderID.String()),
		zap.String("inputID", model.InputID.String()))
	return nil
}

func (uc *CreateOrderInput) Process(entity *domain.OrderInput) error {
	uc.Logger.Info("Processing order input creation",
		zap.String("orderID", entity.OrderID.String()),
		zap.String("inputID", entity.InputID.String()),
		zap.Int("quantity", entity.Quantity))

	// Mapeia entidade para modelo usando persistence
	model := persistence.OrderInputPersistence{}.ToModel(entity)
	uc.Logger.Info("Model created",
		zap.String("orderID", model.OrderID.String()),
		zap.String("inputID", model.InputID.String()),
		zap.Int("quantity", model.Quantity),
		zap.Float64("unitPrice", model.UnitPrice),
		zap.Float64("totalPrice", model.TotalPrice))

	// Salva no banco
	err := uc.SaveOrderInputToDB(model)
	if err != nil {
		return err
	}

	return nil
}
