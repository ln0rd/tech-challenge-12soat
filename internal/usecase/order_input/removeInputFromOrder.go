package order_input

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/logger"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/repository"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/input"
	"go.uber.org/zap"
)

type RemoveInputFromOrder struct {
	OrderRepository       repository.OrderRepository
	InputRepository       repository.InputRepository
	OrderInputRepository  repository.OrderInputRepository
	Logger                logger.Logger
	IncreaseQuantityInput *input.IncreaseQuantityInput
}

// FetchOrderFromDB busca um order específico do banco de dados
func (uc *RemoveInputFromOrder) FetchOrderFromDB(orderID uuid.UUID) (*models.Order, error) {
	order, err := uc.OrderRepository.FindByID(orderID)
	if err != nil {
		uc.Logger.Error("Order not found", zap.String("orderID", orderID.String()))
		return nil, errors.New("order not found")
	}
	uc.Logger.Info("Order found", zap.String("orderID", orderID.String()), zap.String("status", order.Status))
	return order, nil
}

// FetchInputFromDB busca um input específico do banco de dados
func (uc *RemoveInputFromOrder) FetchInputFromDB(inputID uuid.UUID) (*models.Input, error) {
	input, err := uc.InputRepository.FindByID(inputID)
	if err != nil {
		uc.Logger.Error("Input not found", zap.String("inputID", inputID.String()))
		return nil, errors.New("input not found")
	}
	uc.Logger.Info("Input found",
		zap.String("inputID", inputID.String()),
		zap.String("name", input.Name),
		zap.String("inputType", input.InputType),
		zap.Int("currentQuantity", input.Quantity))
	return input, nil
}

// ValidateQuantityToRemove valida se a quantidade a remover é válida
func (uc *RemoveInputFromOrder) ValidateQuantityToRemove(quantityToRemove int) error {
	if quantityToRemove <= 0 {
		uc.Logger.Error("Invalid quantity to remove", zap.Int("quantityToRemove", quantityToRemove))
		return errors.New("quantity to remove must be greater than zero")
	}
	return nil
}

// FetchOrderInputFromDB busca um order input específico do banco de dados
func (uc *RemoveInputFromOrder) FetchOrderInputFromDB(orderID, inputID uuid.UUID) (*models.OrderInput, error) {
	// Busca todos os order inputs para este order
	orderInputs, err := uc.OrderInputRepository.FindByOrderID(orderID)
	if err != nil {
		return nil, err
	}

	// Procura por um order input com o input_id específico
	for _, orderInput := range orderInputs {
		if orderInput.InputID == inputID {
			uc.Logger.Info("Order input found",
				zap.String("orderInputID", orderInput.ID.String()),
				zap.String("orderID", orderInput.OrderID.String()),
				zap.String("inputID", orderInput.InputID.String()),
				zap.Int("quantity", orderInput.Quantity),
				zap.Float64("unitPrice", orderInput.UnitPrice),
				zap.Float64("totalPrice", orderInput.TotalPrice))
			return &orderInput, nil
		}
	}

	uc.Logger.Error("Order input not found",
		zap.String("orderID", orderID.String()),
		zap.String("inputID", inputID.String()))
	return nil, errors.New("order input not found")
}

// ValidateOrderInputQuantity valida se a quantidade no order input é válida
func (uc *RemoveInputFromOrder) ValidateOrderInputQuantity(orderInput *models.OrderInput, quantityToRemove int) error {
	if orderInput.Quantity <= 0 {
		uc.Logger.Error("Invalid quantity in order input",
			zap.String("orderInputID", orderInput.ID.String()),
			zap.Int("quantity", orderInput.Quantity))
		return errors.New("invalid quantity in order input")
	}

	// Verifica se há quantidade suficiente no order_input para remover
	if orderInput.Quantity < quantityToRemove {
		uc.Logger.Error("Insufficient quantity in order input",
			zap.String("orderInputID", orderInput.ID.String()),
			zap.Int("currentQuantity", orderInput.Quantity),
			zap.Int("quantityToRemove", quantityToRemove))
		return errors.New("insufficient quantity in order input")
	}

	return nil
}

// IncreaseInputQuantity aumenta a quantidade do input
func (uc *RemoveInputFromOrder) IncreaseInputQuantity(input *models.Input, quantityToRemove int) error {
	if input.InputType == "service" {
		uc.Logger.Info("Skipping quantity increase for service type",
			zap.String("inputID", input.ID.String()),
			zap.String("name", input.Name),
			zap.String("inputType", input.InputType))
		return nil
	}

	// Usa o usecase de increase quantity
	err := uc.IncreaseQuantityInput.Process(input.ID, quantityToRemove)
	if err != nil {
		uc.Logger.Error("Error increasing input quantity", zap.Error(err))
		return err
	}

	uc.Logger.Info("Input quantity increased successfully",
		zap.String("inputID", input.ID.String()),
		zap.String("name", input.Name),
		zap.Int("quantityIncreased", quantityToRemove))

	return nil
}

// CalculateNewOrderInputValues calcula os novos valores do order input
func (uc *RemoveInputFromOrder) CalculateNewOrderInputValues(orderInput *models.OrderInput, quantityToRemove int) (int, float64) {
	newQuantity := orderInput.Quantity - quantityToRemove
	newTotalPrice := float64(newQuantity) * orderInput.UnitPrice

	uc.Logger.Info("Calculated new order input values",
		zap.String("orderInputID", orderInput.ID.String()),
		zap.Int("oldQuantity", orderInput.Quantity),
		zap.Int("newQuantity", newQuantity),
		zap.Float64("oldTotalPrice", orderInput.TotalPrice),
		zap.Float64("newTotalPrice", newTotalPrice))

	return newQuantity, newTotalPrice
}

// RemoveOrderInputFromDB remove o order input do banco de dados
func (uc *RemoveInputFromOrder) RemoveOrderInputFromDB(orderID, inputID uuid.UUID, quantityToRemove int) error {
	// Busca todos os order inputs para este order
	orderInputs, err := uc.OrderInputRepository.FindByOrderID(orderID)
	if err != nil {
		return err
	}

	// Procura por um order input com o input_id específico
	for _, orderInput := range orderInputs {
		if orderInput.InputID == inputID {
			err := uc.OrderInputRepository.Delete(orderInput.ID)
			if err != nil {
				uc.Logger.Error("Database error removing order input", zap.Error(err))
				return err
			}

			uc.Logger.Info("Order input removed successfully",
				zap.String("orderInputID", orderInput.ID.String()),
				zap.String("orderID", orderID.String()),
				zap.String("inputID", inputID.String()),
				zap.Int("quantityRemoved", quantityToRemove))
			return nil
		}
	}

	return errors.New("order input not found")
}

// UpdateOrderInputInDB atualiza o order input no banco de dados
func (uc *RemoveInputFromOrder) UpdateOrderInputInDB(orderInput *models.OrderInput, newQuantity int, newTotalPrice float64, quantityToRemove int) error {
	// Atualiza os valores
	orderInput.Quantity = newQuantity
	orderInput.TotalPrice = newTotalPrice

	err := uc.OrderInputRepository.Update(orderInput)
	if err != nil {
		uc.Logger.Error("Database error updating order input", zap.Error(err))
		return err
	}

	uc.Logger.Info("Order input updated successfully",
		zap.String("orderInputID", orderInput.ID.String()),
		zap.Int("oldQuantity", orderInput.Quantity+quantityToRemove),
		zap.Int("newQuantity", newQuantity),
		zap.Float64("oldTotalPrice", orderInput.TotalPrice+float64(quantityToRemove)*orderInput.UnitPrice),
		zap.Float64("newTotalPrice", newTotalPrice))

	return nil
}

func (uc *RemoveInputFromOrder) Process(orderID uuid.UUID, inputID uuid.UUID, quantityToRemove int) error {
	uc.Logger.Info("Processing remove input from order",
		zap.String("orderID", orderID.String()),
		zap.String("inputID", inputID.String()),
		zap.Int("quantityToRemove", quantityToRemove))

	// Valida quantidade a remover
	if err := uc.ValidateQuantityToRemove(quantityToRemove); err != nil {
		return err
	}

	// Busca order
	_, err := uc.FetchOrderFromDB(orderID)
	if err != nil {
		return err
	}

	// Busca input
	input, err := uc.FetchInputFromDB(inputID)
	if err != nil {
		return err
	}

	// Busca order input
	orderInput, err := uc.FetchOrderInputFromDB(orderID, inputID)
	if err != nil {
		return err
	}

	// Valida quantidade no order input
	if err := uc.ValidateOrderInputQuantity(orderInput, quantityToRemove); err != nil {
		return err
	}

	uc.Logger.Info("Validation passed, proceeding with input quantity increase")

	// Aumenta a quantidade do input
	if err := uc.IncreaseInputQuantity(input, quantityToRemove); err != nil {
		return err
	}

	// Calcula novos valores
	newQuantity, newTotalPrice := uc.CalculateNewOrderInputValues(orderInput, quantityToRemove)

	// Se a nova quantidade for 0, remove o registro
	if newQuantity == 0 {
		uc.Logger.Info("New quantity is 0, removing order input record")
		return uc.RemoveOrderInputFromDB(orderID, inputID, quantityToRemove)
	}

	// Atualiza o order_input com a nova quantidade e total_price
	return uc.UpdateOrderInputInDB(orderInput, newQuantity, newTotalPrice, quantityToRemove)
}
