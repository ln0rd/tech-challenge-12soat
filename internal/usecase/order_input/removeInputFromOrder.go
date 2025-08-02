package order_input

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/input"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RemoveInputFromOrder struct {
	DB                    *gorm.DB
	Logger                *zap.Logger
	IncreaseQuantityInput *input.IncreaseQuantityInput
}

func (uc *RemoveInputFromOrder) Process(orderID uuid.UUID, inputID uuid.UUID, quantityToRemove int) error {
	uc.Logger.Info("Processing remove input from order",
		zap.String("orderID", orderID.String()),
		zap.String("inputID", inputID.String()),
		zap.Int("quantityToRemove", quantityToRemove))

	// Verifica se o order existe
	var order models.Order
	if err := uc.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		uc.Logger.Error("Order not found", zap.String("orderID", orderID.String()))
		return errors.New("order not found")
	}
	uc.Logger.Info("Order found", zap.String("orderID", orderID.String()), zap.String("status", order.Status))

	// Verifica se o input existe
	var input models.Input
	if err := uc.DB.Where("id = ?", inputID).First(&input).Error; err != nil {
		uc.Logger.Error("Input not found", zap.String("inputID", inputID.String()))
		return errors.New("input not found")
	}
	uc.Logger.Info("Input found",
		zap.String("inputID", inputID.String()),
		zap.String("name", input.Name),
		zap.String("inputType", input.InputType),
		zap.Int("currentQuantity", input.Quantity))

	// Valida se a quantidade a remover é válida
	if quantityToRemove <= 0 {
		uc.Logger.Error("Invalid quantity to remove", zap.Int("quantityToRemove", quantityToRemove))
		return errors.New("quantity to remove must be greater than zero")
	}

	// Busca o order_input específico
	var orderInput models.OrderInput
	if err := uc.DB.Where("order_id = ? AND input_id = ?", orderID, inputID).First(&orderInput).Error; err != nil {
		uc.Logger.Error("Order input not found",
			zap.String("orderID", orderID.String()),
			zap.String("inputID", inputID.String()))
		return errors.New("order input not found")
	}

	uc.Logger.Info("Order input found",
		zap.String("orderInputID", orderInput.ID.String()),
		zap.String("orderID", orderInput.OrderID.String()),
		zap.String("inputID", orderInput.InputID.String()),
		zap.Int("quantity", orderInput.Quantity),
		zap.Float64("unitPrice", orderInput.UnitPrice),
		zap.Float64("totalPrice", orderInput.TotalPrice))

	// Valida se a quantidade a remover é válida
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

	uc.Logger.Info("Validation passed, proceeding with input quantity increase")

	// Aumenta a quantidade do input apenas se não for service
	if input.InputType != "service" {
		uc.Logger.Info("Increasing input quantity (not service type)",
			zap.String("inputID", inputID.String()),
			zap.String("name", input.Name),
			zap.String("inputType", input.InputType))

		err := uc.IncreaseQuantityInput.Process(inputID, quantityToRemove)
		if err != nil {
			uc.Logger.Error("Error increasing input quantity", zap.Error(err))
			return err
		}

		uc.Logger.Info("Input quantity increased successfully")
	} else {
		uc.Logger.Info("Skipping quantity increase for service type",
			zap.String("inputID", inputID.String()),
			zap.String("name", input.Name),
			zap.String("inputType", input.InputType))
	}

	// Calcula a nova quantidade e total_price
	newQuantity := orderInput.Quantity - quantityToRemove
	newTotalPrice := float64(newQuantity) * orderInput.UnitPrice

	uc.Logger.Info("Calculated new values",
		zap.Int("oldQuantity", orderInput.Quantity),
		zap.Int("newQuantity", newQuantity),
		zap.Float64("oldTotalPrice", orderInput.TotalPrice),
		zap.Float64("newTotalPrice", newTotalPrice))

	// Se a nova quantidade for 0, remove o registro
	if newQuantity == 0 {
		uc.Logger.Info("New quantity is 0, removing order input record")
		result := uc.DB.Where("order_id = ? AND input_id = ?", orderID, inputID).Delete(&models.OrderInput{})
		if result.Error != nil {
			uc.Logger.Error("Database error removing order input", zap.Error(result.Error))
			return result.Error
		}

		uc.Logger.Info("Order input removed successfully (quantity became 0)",
			zap.String("orderID", orderID.String()),
			zap.String("inputID", inputID.String()),
			zap.Int("quantityReturned", quantityToRemove),
			zap.Int64("rowsAffected", result.RowsAffected))
	} else {
		// Atualiza o order_input com a nova quantidade e total_price
		result := uc.DB.Model(&orderInput).Updates(map[string]interface{}{
			"quantity":    newQuantity,
			"total_price": newTotalPrice,
		})
		if result.Error != nil {
			uc.Logger.Error("Database error updating order input", zap.Error(result.Error))
			return result.Error
		}

		uc.Logger.Info("Order input updated successfully",
			zap.String("orderID", orderID.String()),
			zap.String("inputID", inputID.String()),
			zap.Int("oldQuantity", orderInput.Quantity),
			zap.Int("newQuantity", newQuantity),
			zap.Float64("oldTotalPrice", orderInput.TotalPrice),
			zap.Float64("newTotalPrice", newTotalPrice),
			zap.Int("quantityReturned", quantityToRemove),
			zap.Int64("rowsAffected", result.RowsAffected))
	}

	return nil
}
