package input

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IncreaseQuantityInput struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (uc *IncreaseQuantityInput) Process(id uuid.UUID, quantity int) error {
	uc.Logger.Info("Processing increase quantity for input",
		zap.String("id", id.String()),
		zap.Int("quantityToIncrease", quantity))

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
		zap.Int("quantityToIncrease", quantity))

	// Valida se a quantidade a aumentar é válida
	if quantity <= 0 {
		uc.Logger.Error("Invalid quantity to increase", zap.Int("quantity", quantity))
		return errors.New("quantity to increase must be greater than zero")
	}

	// Calcula a nova quantidade
	newQuantity := input.Quantity + quantity

	uc.Logger.Info("Calculated new quantity",
		zap.Int("currentQuantity", input.Quantity),
		zap.Int("quantityToIncrease", quantity),
		zap.Int("newQuantity", newQuantity))

	// Atualiza a quantidade
	result := uc.DB.Model(&input).Update("quantity", newQuantity)
	if result.Error != nil {
		uc.Logger.Error("Database error updating input quantity", zap.Error(result.Error))
		return result.Error
	}

	uc.Logger.Info("Input quantity increased successfully",
		zap.String("id", input.ID.String()),
		zap.String("name", input.Name),
		zap.Int("oldQuantity", input.Quantity),
		zap.Int("newQuantity", newQuantity),
		zap.Int64("rowsAffected", result.RowsAffected))

	return nil
}
