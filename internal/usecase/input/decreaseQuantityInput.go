package input

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DecreaseQuantityInput struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (uc *DecreaseQuantityInput) Process(id uuid.UUID, quantity int) error {
	uc.Logger.Info("Processing decrease quantity for input",
		zap.String("id", id.String()),
		zap.Int("quantityToDecrease", quantity))

	// Verifica se o input existe
	var input models.Input
	if err := uc.DB.Where("id = ?", id).First(&input).Error; err != nil {
		uc.Logger.Error("Input not found", zap.String("id", id.String()))
		return errors.New("input not found")
	}

	uc.Logger.Info("Found input",
		zap.String("id", input.ID.String()),
		zap.String("name", input.Name),
		zap.Int("currentQuantity", input.Quantity),
		zap.Int("quantityToDecrease", quantity))

	// Valida se a quantidade a diminuir é válida
	if quantity <= 0 {
		uc.Logger.Error("Invalid quantity to decrease", zap.Int("quantity", quantity))
		return errors.New("quantity to decrease must be greater than zero")
	}

	// Calcula a nova quantidade
	newQuantity := input.Quantity - quantity
	if newQuantity < 0 {
		uc.Logger.Error("Insufficient quantity",
			zap.Int("currentQuantity", input.Quantity),
			zap.Int("quantityToDecrease", quantity),
			zap.Int("newQuantity", newQuantity))
		return errors.New("insufficient quantity")
	}

	uc.Logger.Info("Calculated new quantity",
		zap.Int("currentQuantity", input.Quantity),
		zap.Int("quantityToDecrease", quantity),
		zap.Int("newQuantity", newQuantity))

	// Atualiza a quantidade
	result := uc.DB.Model(&input).Update("quantity", newQuantity)
	if result.Error != nil {
		uc.Logger.Error("Database error updating input quantity", zap.Error(result.Error))
		return result.Error
	}

	uc.Logger.Info("Input quantity decreased successfully",
		zap.String("id", input.ID.String()),
		zap.String("name", input.Name),
		zap.Int("oldQuantity", input.Quantity),
		zap.Int("newQuantity", newQuantity),
		zap.Int64("rowsAffected", result.RowsAffected))

	return nil
}
