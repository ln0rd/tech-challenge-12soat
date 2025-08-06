package order_input

import (
	"errors"

	"github.com/google/uuid"
	interfaces "github.com/ln0rd/tech_challenge_12soat/internal/domain/interfaces"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/repository"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/input"
	"go.uber.org/zap"
)

type AddInputToOrder struct {
	OrderRepository       repository.OrderRepository
	InputRepository       repository.InputRepository
	OrderInputRepository  repository.OrderInputRepository
	Logger                interfaces.Logger
	DecreaseQuantityInput *input.DecreaseQuantityInput
}

// FetchOrderFromDB busca um order específico do banco de dados
func (uc *AddInputToOrder) FetchOrderFromDB(orderID uuid.UUID) (*models.Order, error) {
	order, err := uc.OrderRepository.FindByID(orderID)
	if err != nil {
		uc.Logger.Error("Order not found", zap.String("orderID", orderID.String()))
		return nil, errors.New("order not found")
	}
	uc.Logger.Info("Order found", zap.String("orderID", orderID.String()), zap.String("status", order.Status))
	return order, nil
}

// FetchInputFromDB busca um input específico do banco de dados
func (uc *AddInputToOrder) FetchInputFromDB(inputID uuid.UUID) (*models.Input, error) {
	input, err := uc.InputRepository.FindByID(inputID)
	if err != nil {
		uc.Logger.Error("Input not found", zap.String("inputID", inputID.String()))
		return nil, errors.New("input not found")
	}
	uc.Logger.Info("Input found",
		zap.String("inputID", inputID.String()),
		zap.String("name", input.Name),
		zap.String("inputType", input.InputType),
		zap.Int("availableQuantity", input.Quantity))
	return input, nil
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
	// Busca todos os order inputs para este order
	orderInputs, err := uc.OrderInputRepository.FindByOrderID(orderID)
	if err != nil {
		return nil, err
	}

	// Procura por um order input com o input_id específico
	for _, orderInput := range orderInputs {
		if orderInput.InputID == inputID {
			uc.Logger.Info("Existing order input found",
				zap.String("orderInputID", orderInput.ID.String()),
				zap.String("orderID", orderInput.OrderID.String()),
				zap.String("inputID", orderInput.InputID.String()),
				zap.Int("currentQuantity", orderInput.Quantity),
				zap.Float64("currentTotalPrice", orderInput.TotalPrice))
			return &orderInput, nil
		}
	}

	uc.Logger.Info("No existing order input found",
		zap.String("orderID", orderID.String()),
		zap.String("inputID", inputID.String()))
	return nil, nil
}

// UpdateExistingOrderInput atualiza um order input existente
func (uc *AddInputToOrder) UpdateExistingOrderInput(existingOrderInput *models.OrderInput, quantity int, unitPrice float64) error {
	// Calcula novos valores
	newQuantity := existingOrderInput.Quantity + quantity
	newTotalPrice := float64(newQuantity) * unitPrice

	// Atualiza o order input
	existingOrderInput.Quantity = newQuantity
	existingOrderInput.TotalPrice = newTotalPrice

	err := uc.OrderInputRepository.Update(existingOrderInput)
	if err != nil {
		uc.Logger.Error("Database error updating existing order input", zap.Error(err))
		return err
	}

	uc.Logger.Info("Existing order input updated successfully",
		zap.String("orderInputID", existingOrderInput.ID.String()),
		zap.Int("oldQuantity", existingOrderInput.Quantity-quantity),
		zap.Int("newQuantity", newQuantity),
		zap.Float64("oldTotalPrice", existingOrderInput.TotalPrice-float64(quantity)*unitPrice),
		zap.Float64("newTotalPrice", newTotalPrice))

	return nil
}

// DecreaseInputQuantity diminui a quantidade do input
func (uc *AddInputToOrder) DecreaseInputQuantity(input *models.Input, quantity int) error {
	// Para inputs do tipo "service", não diminuímos a quantidade
	if input.InputType == "service" {
		uc.Logger.Info("Skipping quantity decrease for service type",
			zap.String("inputID", input.ID.String()),
			zap.String("name", input.Name),
			zap.String("inputType", input.InputType))
		return nil
	}

	// Usa o usecase de decrease quantity
	err := uc.DecreaseQuantityInput.Process(input.ID, quantity)
	if err != nil {
		uc.Logger.Error("Error decreasing input quantity", zap.Error(err))
		return err
	}

	uc.Logger.Info("Input quantity decreased successfully",
		zap.String("inputID", input.ID.String()),
		zap.String("name", input.Name),
		zap.Int("quantityDecreased", quantity))

	return nil
}

// CreateNewOrderInput cria um novo order input
func (uc *AddInputToOrder) CreateNewOrderInput(orderID, inputID uuid.UUID, quantity int, unitPrice float64) error {
	totalPrice := float64(quantity) * unitPrice

	newOrderInput := &models.OrderInput{
		ID:         uuid.New(),
		OrderID:    orderID,
		InputID:    inputID,
		Quantity:   quantity,
		UnitPrice:  unitPrice,
		TotalPrice: totalPrice,
	}

	err := uc.OrderInputRepository.Create(newOrderInput)
	if err != nil {
		uc.Logger.Error("Database error creating new order input", zap.Error(err))
		return err
	}

	uc.Logger.Info("New order input created successfully",
		zap.String("orderInputID", newOrderInput.ID.String()),
		zap.String("orderID", newOrderInput.OrderID.String()),
		zap.String("inputID", newOrderInput.InputID.String()),
		zap.Int("quantity", newOrderInput.Quantity),
		zap.Float64("unitPrice", newOrderInput.UnitPrice),
		zap.Float64("totalPrice", newOrderInput.TotalPrice))

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
