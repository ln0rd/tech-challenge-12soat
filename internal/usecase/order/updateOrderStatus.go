package order

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/order_status_history"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UpdateOrderStatus struct {
	DB                   *gorm.DB
	Logger               *zap.Logger
	StatusHistoryManager *order_status_history.ManageOrderStatusHistory
}

// GetValidStatuses retorna a lista de status válidos
func (uc *UpdateOrderStatus) GetValidStatuses() []string {
	return []string{
		"Received",
		"Undergoing diagnosis",
		"Awaiting approval",
		"In progress",
		"Completed",
		"Delivered",
		"Canceled",
	}
}

// ValidateOrderStatus valida se o status é válido
func (uc *UpdateOrderStatus) ValidateOrderStatus(newStatus string) error {
	validStatuses := uc.GetValidStatuses()

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
	return nil
}

// FetchOrderFromDB busca um order específico do banco de dados
func (uc *UpdateOrderStatus) FetchOrderFromDB(orderID uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := uc.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		uc.Logger.Error("Order not found", zap.String("orderID", orderID.String()))
		return nil, errors.New("order not found")
	}

	uc.Logger.Info("Order found",
		zap.String("orderID", order.ID.String()),
		zap.String("currentStatus", order.Status))

	return &order, nil
}

// UpdateOrderStatusInDB atualiza o status da order no banco de dados
func (uc *UpdateOrderStatus) UpdateOrderStatusInDB(order *models.Order, newStatus string) error {
	result := uc.DB.Model(order).Update("status", newStatus)
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

// UpdateStatusHistory atualiza o histórico de status
func (uc *UpdateOrderStatus) UpdateStatusHistory(orderID uuid.UUID, newStatus string) error {
	historyErr := uc.StatusHistoryManager.UpdateStatus(orderID, newStatus)
	if historyErr != nil {
		uc.Logger.Error("Error updating status history", zap.Error(historyErr))
		// Não retorna erro aqui, pois a order já foi atualizada
		// Apenas loga o erro para monitoramento
		return nil
	}

	uc.Logger.Info("Status history updated successfully",
		zap.String("orderID", orderID.String()),
		zap.String("newStatus", newStatus))

	return nil
}

func (uc *UpdateOrderStatus) Process(orderID uuid.UUID, newStatus string) error {
	uc.Logger.Info("Processing update order status",
		zap.String("orderID", orderID.String()),
		zap.String("newStatus", newStatus))

	// Valida se o status é válido
	if err := uc.ValidateOrderStatus(newStatus); err != nil {
		return err
	}

	// Busca a order
	order, err := uc.FetchOrderFromDB(orderID)
	if err != nil {
		return err
	}

	uc.Logger.Info("Order found",
		zap.String("orderID", order.ID.String()),
		zap.String("currentStatus", order.Status),
		zap.String("newStatus", newStatus))

	// Atualiza o status da order
	err = uc.UpdateOrderStatusInDB(order, newStatus)
	if err != nil {
		return err
	}

	// Atualiza o histórico de status
	uc.UpdateStatusHistory(orderID, newStatus)

	return nil
}
