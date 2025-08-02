package order

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UpdateOrderStatus struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (uc *UpdateOrderStatus) Process(orderID uuid.UUID, newStatus string) error {
	uc.Logger.Info("Processing update order status",
		zap.String("orderID", orderID.String()),
		zap.String("newStatus", newStatus))

	// Valida se o status é válido
	validStatuses := []string{
		"Received",
		"Undergoing diagnosis",
		"Awaiting approval",
		"In progress",
		"Completed",
		"Delivered",
	}

	isValidStatus := false
	for _, status := range validStatuses {
		if status == newStatus {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		uc.Logger.Error("Invalid order status",
			zap.String("newStatus", newStatus),
			zap.Strings("validStatuses", validStatuses))
		return errors.New("invalid order status")
	}

	uc.Logger.Info("Status validation passed", zap.String("newStatus", newStatus))

	// Verifica se a order existe
	var order models.Order
	if err := uc.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		uc.Logger.Error("Order not found", zap.String("orderID", orderID.String()))
		return errors.New("order not found")
	}

	uc.Logger.Info("Order found",
		zap.String("orderID", order.ID.String()),
		zap.String("currentStatus", order.Status),
		zap.String("newStatus", newStatus))

	// Atualiza o status da order
	result := uc.DB.Model(&order).Update("status", newStatus)
	if result.Error != nil {
		uc.Logger.Error("Database error updating order status", zap.Error(result.Error))
		return result.Error
	}

	uc.Logger.Info("Order status updated successfully",
		zap.String("orderID", order.ID.String()),
		zap.String("oldStatus", order.Status),
		zap.String("newStatus", newStatus),
		zap.Int64("rowsAffected", result.RowsAffected))

	return nil
}
