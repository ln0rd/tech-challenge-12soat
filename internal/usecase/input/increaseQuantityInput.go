package input

import (
	"errors"

	"github.com/google/uuid"
	interfaces "github.com/ln0rd/tech_challenge_12soat/internal/domain/interfaces"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/repository"
	"go.uber.org/zap"
)

type IncreaseQuantityInput struct {
	InputRepository repository.InputRepository
	Logger          interfaces.Logger
}

// FetchInputFromDB busca um input específico do banco de dados
func (uc *IncreaseQuantityInput) FetchInputFromDB(id uuid.UUID) (*models.Input, error) {
	input, err := uc.InputRepository.FindByID(id)
	if err != nil {
		uc.Logger.Error("Input not found", zap.String("id", id.String()))
		return nil, errors.New("input not found")
	}

	uc.Logger.Info("Found input",
		zap.String("id", input.ID.String()),
		zap.String("name", input.Name),
		zap.Int("currentQuantity", input.Quantity))

	return input, nil
}

// ValidateQuantityToIncrease valida se a quantidade a aumentar é válida
func (uc *IncreaseQuantityInput) ValidateQuantityToIncrease(quantity int) error {
	if quantity <= 0 {
		uc.Logger.Error("Invalid quantity to increase", zap.Int("quantity", quantity))
		return errors.New("quantity to increase must be greater than zero")
	}
	return nil
}

// CalculateNewQuantity calcula a nova quantidade após o aumento
func (uc *IncreaseQuantityInput) CalculateNewQuantity(currentQuantity, quantityToIncrease int) int {
	newQuantity := currentQuantity + quantityToIncrease

	uc.Logger.Info("Calculated new quantity",
		zap.Int("currentQuantity", currentQuantity),
		zap.Int("quantityToIncrease", quantityToIncrease),
		zap.Int("newQuantity", newQuantity))

	return newQuantity
}

// UpdateInputQuantity atualiza a quantidade do input no banco de dados
func (uc *IncreaseQuantityInput) UpdateInputQuantity(input *models.Input, newQuantity int) error {
	// Atualiza a quantidade no modelo
	input.Quantity = newQuantity

	err := uc.InputRepository.Update(input)
	if err != nil {
		uc.Logger.Error("Database error updating input quantity", zap.Error(err))
		return err
	}

	uc.Logger.Info("Input quantity increased successfully",
		zap.String("id", input.ID.String()),
		zap.String("name", input.Name),
		zap.Int("oldQuantity", input.Quantity),
		zap.Int("newQuantity", newQuantity))

	return nil
}

func (uc *IncreaseQuantityInput) Process(id uuid.UUID, quantity int) error {
	uc.Logger.Info("Processing increase quantity for input",
		zap.String("id", id.String()),
		zap.Int("quantityToIncrease", quantity))

	// Valida quantidade a aumentar
	if err := uc.ValidateQuantityToIncrease(quantity); err != nil {
		return err
	}

	// Busca o input
	input, err := uc.FetchInputFromDB(id)
	if err != nil {
		return err
	}

	// Calcula nova quantidade
	newQuantity := uc.CalculateNewQuantity(input.Quantity, quantity)

	// Atualiza a quantidade
	err = uc.UpdateInputQuantity(input, newQuantity)
	if err != nil {
		return err
	}

	return nil
}
