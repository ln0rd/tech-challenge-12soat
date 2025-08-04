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

// FetchOrderFromDB busca um order específico do banco de dados
func (uc *AddInputToOrder) FetchOrderFromDB(orderID uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := uc.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		uc.Logger.Error("Order not found", zap.String("orderID", orderID.String()))
		return nil, errors.New("order not found")
	}
	uc.Logger.Info("Order found", zap.String("orderID", orderID.String()), zap.String("status", order.Status))
	return &order, nil
}

// FetchInputFromDB busca um input específico do banco de dados
func (uc *AddInputToOrder) FetchInputFromDB(inputID uuid.UUID) (*models.Input, error) {
	var input models.Input
	if err := uc.DB.Where("id = ?", inputID).First(&input).Error; err != nil {
		uc.Logger.Error("Input not found", zap.String("inputID", inputID.String()))
		return nil, errors.New("input not found")
	}
	uc.Logger.Info("Input found",
		zap.String("inputID", inputID.String()),
		zap.String("name", input.Name),
		zap.String("inputType", input.InputType),
		zap.Int("availableQuantity", input.Quantity))
	return &input, nil
}

// ValidateQuantity valida se a quantidade é válida
func (uc *AddInputToOrder) ValidateQuantity(quantity int) error {
	if quantity <= 0 {
		uc.Logger.Error("Invalid quantity", zap.Int("quantity", quantity))
		return errors.New("quantity must be greater than zero")
	}
	return nil
}

// ValidateInputAvailability valida se há quantidade suficiente do input
func (uc *AddInputToOrder) ValidateInputAvailability(input *models.Input, quantity int) error {
	// Para inputs do tipo "service", não fazemos controle de estoque
	if input.InputType == "service" {
		uc.Logger.Info("Input is service type, skipping stock control",
			zap.String("inputID", input.ID.String()),
			zap.String("name", input.Name),
			zap.String("inputType", input.InputType))
		return nil
	}

	// Verifica se há quantidade suficiente (apenas para inputs que não são service)
	if input.Quantity < quantity {
		uc.Logger.Error("Insufficient input quantity",
			zap.String("inputID", input.ID.String()),
			zap.String("name", input.Name),
			zap.Int("requestedQuantity", quantity),
			zap.Int("availableQuantity", input.Quantity))
		return errors.New("insufficient input quantity")
	}

	return nil
}

// ValidateInputPrice valida se o preço do input é válido
func (uc *AddInputToOrder) ValidateInputPrice(input *models.Input) (float64, error) {
	unitPrice := input.Price
	if unitPrice <= 0 {
		uc.Logger.Error("Input has invalid price",
			zap.String("inputID", input.ID.String()),
			zap.String("name", input.Name),
			zap.Float64("price", unitPrice))
		return 0, errors.New("input has invalid price")
	}

	uc.Logger.Info("Input price retrieved",
		zap.String("inputID", input.ID.String()),
		zap.String("name", input.Name),
		zap.Float64("unitPrice", unitPrice))

	return unitPrice, nil
}

// FetchExistingOrderInput busca um order input existente
func (uc *AddInputToOrder) FetchExistingOrderInput(orderID, inputID uuid.UUID) (*models.OrderInput, error) {
	var existingOrderInput models.OrderInput
	if err := uc.DB.Where("order_id = ? AND input_id = ?", orderID, inputID).First(&existingOrderInput).Error; err == nil {
		uc.Logger.Info("Existing order input found",
			zap.String("orderInputID", existingOrderInput.ID.String()),
			zap.Int("currentQuantity", existingOrderInput.Quantity),
			zap.Float64("currentTotalPrice", existingOrderInput.TotalPrice))
		return &existingOrderInput, nil
	} else if err != gorm.ErrRecordNotFound {
		uc.Logger.Error("Database error checking existing order input", zap.Error(err))
		return nil, err
	}

	uc.Logger.Info("No existing order input found")
	return nil, nil
}

// UpdateExistingOrderInput atualiza um order input existente
func (uc *AddInputToOrder) UpdateExistingOrderInput(existingOrderInput *models.OrderInput, quantity int, unitPrice float64) error {
	// Calcula a nova quantidade
	newQuantity := existingOrderInput.Quantity + quantity
	newTotalPrice := float64(newQuantity) * unitPrice

	uc.Logger.Info("Calculated new values",
		zap.Int("newQuantity", newQuantity),
		zap.Float64("newTotalPrice", newTotalPrice))

	// Atualiza o order_input existente
	result := uc.DB.Model(existingOrderInput).Updates(map[string]interface{}{
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
}

// DecreaseInputQuantity diminui a quantidade do input
func (uc *AddInputToOrder) DecreaseInputQuantity(input *models.Input, quantity int) error {
	if input.InputType == "service" {
		uc.Logger.Info("Skipping quantity decrease for service type",
			zap.String("inputID", input.ID.String()),
			zap.String("name", input.Name),
			zap.String("inputType", input.InputType))
		return nil
	}

	uc.Logger.Info("Decreasing input quantity (not service type)",
		zap.String("inputID", input.ID.String()),
		zap.String("name", input.Name),
		zap.String("inputType", input.InputType))

	err := uc.DecreaseQuantityInput.Process(input.ID, quantity)
	if err != nil {
		uc.Logger.Error("Error decreasing input quantity", zap.Error(err))
		return err
	}

	uc.Logger.Info("Input quantity decreased successfully")
	return nil
}

// CreateNewOrderInput cria um novo order input
func (uc *AddInputToOrder) CreateNewOrderInput(orderID, inputID uuid.UUID, quantity int, unitPrice float64) error {
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

func (uc *AddInputToOrder) Process(orderID uuid.UUID, inputID uuid.UUID, quantity int) error {
	uc.Logger.Info("Processing add input to order",
		zap.String("orderID", orderID.String()),
		zap.String("inputID", inputID.String()),
		zap.Int("quantity", quantity))

	// Valida quantidade
	if err := uc.ValidateQuantity(quantity); err != nil {
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

	// Valida disponibilidade do input
	if err := uc.ValidateInputAvailability(input, quantity); err != nil {
		return err
	}

	// Valida e obtém preço do input
	unitPrice, err := uc.ValidateInputPrice(input)
	if err != nil {
		return err
	}

	uc.Logger.Info("Validation passed, checking if order input already exists")

	// Verifica se já existe um order_input com o mesmo input_id para este order
	existingOrderInput, err := uc.FetchExistingOrderInput(orderID, inputID)
	if err != nil {
		return err
	}

	if existingOrderInput != nil {
		// Já existe um registro, vamos atualizar a quantidade e o total_price
		uc.Logger.Info("Existing order input found, updating quantity and total price",
			zap.String("orderInputID", existingOrderInput.ID.String()),
			zap.Int("currentQuantity", existingOrderInput.Quantity),
			zap.Int("quantityToAdd", quantity),
			zap.Float64("currentTotalPrice", existingOrderInput.TotalPrice))

		return uc.UpdateExistingOrderInput(existingOrderInput, quantity, unitPrice)
	}

	// Não existe registro, vamos criar um novo
	uc.Logger.Info("No existing order input found, creating new one")

	// Diminui a quantidade do input
	if err := uc.DecreaseInputQuantity(input, quantity); err != nil {
		return err
	}

	// Cria novo order input
	return uc.CreateNewOrderInput(orderID, inputID, quantity, unitPrice)
}
