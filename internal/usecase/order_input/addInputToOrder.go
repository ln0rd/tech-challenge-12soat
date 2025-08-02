package order_input

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/input"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AddInputToOrder struct {
	DB                    *gorm.DB
	Logger                *zap.Logger
	DecreaseQuantityInput *input.DecreaseQuantityInput
}

func (uc *AddInputToOrder) Process(orderID uuid.UUID, inputID uuid.UUID, quantity int) error {
	uc.Logger.Info("Processing add input to order",
		zap.String("orderID", orderID.String()),
		zap.String("inputID", inputID.String()),
		zap.Int("quantity", quantity))

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
		zap.Int("availableQuantity", input.Quantity))

	// Para inputs do tipo "service", não fazemos controle de estoque
	if input.InputType == "service" {
		uc.Logger.Info("Input is service type, skipping stock control",
			zap.String("inputID", inputID.String()),
			zap.String("name", input.Name),
			zap.String("inputType", input.InputType))
	} else {
		// Verifica se há quantidade suficiente (apenas para inputs que não são service)
		if input.Quantity < quantity {
			uc.Logger.Error("Insufficient input quantity",
				zap.String("inputID", inputID.String()),
				zap.String("name", input.Name),
				zap.Int("requestedQuantity", quantity),
				zap.Int("availableQuantity", input.Quantity))
			return errors.New("insufficient input quantity")
		}
	}

	// Verifica se a quantidade é válida
	if quantity <= 0 {
		uc.Logger.Error("Invalid quantity", zap.Int("quantity", quantity))
		return errors.New("quantity must be greater than zero")
	}

	// Busca o preço unitário do input
	unitPrice := input.Price
	if unitPrice <= 0 {
		uc.Logger.Error("Input has invalid price",
			zap.String("inputID", inputID.String()),
			zap.String("name", input.Name),
			zap.Float64("price", unitPrice))
		return errors.New("input has invalid price")
	}

	uc.Logger.Info("Input price retrieved",
		zap.String("inputID", inputID.String()),
		zap.String("name", input.Name),
		zap.Float64("unitPrice", unitPrice))

	uc.Logger.Info("Validation passed, checking if order input already exists")

	// Verifica se já existe um order_input com o mesmo input_id para este order
	var existingOrderInput models.OrderInput
	if err := uc.DB.Where("order_id = ? AND input_id = ?", orderID, inputID).First(&existingOrderInput).Error; err == nil {
		// Já existe um registro, vamos atualizar a quantidade e o total_price
		uc.Logger.Info("Existing order input found, updating quantity and total price",
			zap.String("orderInputID", existingOrderInput.ID.String()),
			zap.Int("currentQuantity", existingOrderInput.Quantity),
			zap.Int("quantityToAdd", quantity),
			zap.Float64("currentTotalPrice", existingOrderInput.TotalPrice))

		// Calcula a nova quantidade
		newQuantity := existingOrderInput.Quantity + quantity
		newTotalPrice := float64(newQuantity) * unitPrice

		uc.Logger.Info("Calculated new values",
			zap.Int("newQuantity", newQuantity),
			zap.Float64("newTotalPrice", newTotalPrice))

		// Atualiza o order_input existente
		result := uc.DB.Model(&existingOrderInput).Updates(map[string]interface{}{
			"quantity":    newQuantity,
			"total_price": newTotalPrice,
		})
		if result.Error != nil {
			uc.Logger.Error("Database error updating order input", zap.Error(result.Error))
			return result.Error
		}

		uc.Logger.Info("OrderInput updated successfully in database",
			zap.String("id", existingOrderInput.ID.String()),
			zap.Int("oldQuantity", existingOrderInput.Quantity),
			zap.Int("newQuantity", newQuantity),
			zap.Float64("oldTotalPrice", existingOrderInput.TotalPrice),
			zap.Float64("newTotalPrice", newTotalPrice),
			zap.Int64("rowsAffected", result.RowsAffected))

		return nil
	} else if err != gorm.ErrRecordNotFound {
		uc.Logger.Error("Database error checking existing order input", zap.Error(err))
		return err
	}

	// Não existe registro, vamos criar um novo
	uc.Logger.Info("No existing order input found, creating new one")

	// Diminui a quantidade do input apenas se não for service
	if input.InputType != "service" {
		uc.Logger.Info("Decreasing input quantity (not service type)",
			zap.String("inputID", inputID.String()),
			zap.String("name", input.Name),
			zap.String("inputType", input.InputType))

		err := uc.DecreaseQuantityInput.Process(inputID, quantity)
		if err != nil {
			uc.Logger.Error("Error decreasing input quantity", zap.Error(err))
			return err
		}

		uc.Logger.Info("Input quantity decreased successfully")
	} else {
		uc.Logger.Info("Skipping quantity decrease for service type",
			zap.String("inputID", inputID.String()),
			zap.String("name", input.Name),
			zap.String("inputType", input.InputType))
	}

	// Calcula o preço total
	totalPrice := float64(quantity) * unitPrice

	uc.Logger.Info("Calculated total price",
		zap.Int("quantity", quantity),
		zap.Float64("unitPrice", unitPrice),
		zap.Float64("totalPrice", totalPrice))

	// Cria o vínculo order_input
	orderInput := &models.OrderInput{
		ID:         uuid.New(),
		OrderID:    orderID,
		InputID:    inputID,
		Quantity:   quantity,
		UnitPrice:  unitPrice,
		TotalPrice: totalPrice,
	}

	uc.Logger.Info("OrderInput model created",
		zap.String("id", orderInput.ID.String()),
		zap.String("orderID", orderInput.OrderID.String()),
		zap.String("inputID", orderInput.InputID.String()),
		zap.Int("quantity", orderInput.Quantity),
		zap.Float64("unitPrice", orderInput.UnitPrice),
		zap.Float64("totalPrice", orderInput.TotalPrice))

	// Salva o vínculo no banco
	result := uc.DB.Create(orderInput)
	if result.Error != nil {
		uc.Logger.Error("Database error creating order input", zap.Error(result.Error))
		return result.Error
	}

	uc.Logger.Info("OrderInput created successfully in database",
		zap.String("id", orderInput.ID.String()),
		zap.Int64("rowsAffected", result.RowsAffected))

	return nil
}
