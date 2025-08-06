package input

import (
	"errors"

	"github.com/google/uuid"
	interfaces "github.com/ln0rd/tech_challenge_12soat/internal/domain/interfaces"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/repository"
	"go.uber.org/zap"
)

type DecreaseQuantityInput struct {
	InputRepository repository.InputRepository
	Logger          interfaces.Logger
}

// FetchInputFromDB busca um input específico do banco de dados
func (uc *DecreaseQuantityInput) FetchInputFromDB(id uuid.UUID) (*models.Input, error) {
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

// ValidateQuantityToDecrease valida se a quantidade a diminuir é válida
func (uc *DecreaseQuantityInput) ValidateQuantityToDecrease(quantity int) error {
	if quantity <= 0 {
		uc.Logger.Error("Invalid quantity to decrease", zap.Int("quantity", quantity))
		return errors.New("quantity to decrease must be greater than zero")
	}
	return nil
}

// CalculateNewQuantity calcula a nova quantidade após a diminuição
func (uc *DecreaseQuantityInput) CalculateNewQuantity(currentQuantity, quantityToDecrease int) (int, error) {
	newQuantity := currentQuantity - quantityToDecrease
	if newQuantity < 0 {
		uc.Logger.Error("Insufficient quantity",
			zap.Int("currentQuantity", currentQuantity),
			zap.Int("quantityToDecrease", quantityToDecrease),
			zap.Int("newQuantity", newQuantity))
		return 0, errors.New("insufficient quantity")
	}

	uc.Logger.Info("Calculated new quantity",
		zap.Int("currentQuantity", currentQuantity),
		zap.Int("quantityToDecrease", quantityToDecrease),
		zap.Int("newQuantity", newQuantity))

	return newQuantity, nil
}

// UpdateInputQuantity atualiza a quantidade do input no banco de dados
func (uc *DecreaseQuantityInput) UpdateInputQuantity(input *models.Input, newQuantity int) error {
	// Atualiza a quantidade no modelo
	input.Quantity = newQuantity

	err := uc.InputRepository.Update(input)
	if err != nil {
		uc.Logger.Error("Database error updating input quantity", zap.Error(err))
		return err
	}

	uc.Logger.Info("Input quantity decreased successfully",
		zap.String("id", input.ID.String()),
		zap.String("name", input.Name),
		zap.Int("oldQuantity", input.Quantity),
		zap.Int("newQuantity", newQuantity))

	return nil
}

func (uc *DecreaseQuantityInput) Process(id uuid.UUID, quantity int) error {
	uc.Logger.Info("Processing decrease quantity for input",
		zap.String("id", id.String()),
		zap.Int("quantityToDecrease", quantity))

	// Valida quantidade a diminuir
	if err := uc.ValidateQuantityToDecrease(quantity); err != nil {
		return err
	}

	// Busca o input
	input, err := uc.FetchInputFromDB(id)
	if err != nil {
		return err
	}

	// Calcula nova quantidade
	newQuantity, err := uc.CalculateNewQuantity(input.Quantity, quantity)
	if err != nil {
		return err
	}

	// Atualiza a quantidade
	err = uc.UpdateInputQuantity(input, newQuantity)
	if err != nil {
		return err
	}

	return nil
}
