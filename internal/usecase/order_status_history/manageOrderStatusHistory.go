package order_status_history

import (
	"time"

	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ManageOrderStatusHistory struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

// Finaliza o status atual e inicia um novo
func (uc *ManageOrderStatusHistory) UpdateStatus(orderID uuid.UUID, newStatus string) error {
	uc.Logger.Info("Managing order status history",
		zap.String("orderID", orderID.String()),
		zap.String("newStatus", newStatus))

	// Finaliza o status atual (se existir)
	err := uc.finalizeCurrentStatus(orderID)
	if err != nil {
		uc.Logger.Error("Error finalizing current status", zap.Error(err))
		return err
	}

	// Inicia o novo status
	err = uc.StartNewStatus(orderID, newStatus)
	if err != nil {
		uc.Logger.Error("Error starting new status", zap.Error(err))
		return err
	}

	uc.Logger.Info("Order status history updated successfully",
		zap.String("orderID", orderID.String()),
		zap.String("newStatus", newStatus))

	return nil
}

// Finaliza o status atual
func (uc *ManageOrderStatusHistory) finalizeCurrentStatus(orderID uuid.UUID) error {
	// Busca o status atual (sem ended_at)
	var currentStatus models.OrderStatusHistory
	if err := uc.DB.Where("order_id = ? AND ended_at IS NULL", orderID).First(&currentStatus).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Não há status atual, isso é normal para a primeira mudança
			uc.Logger.Info("No current status found, this is normal for first status change",
				zap.String("orderID", orderID.String()))
			return nil
		}
		return err
	}

	// Calcula a duração
	now := time.Now()
	duration := int(now.Sub(currentStatus.StartedAt).Minutes())

	// Atualiza o registro atual
	result := uc.DB.Model(&currentStatus).Updates(map[string]interface{}{
		"ended_at":         now,
		"duration_minutes": duration,
	})

	if result.Error != nil {
		uc.Logger.Error("Error updating current status", zap.Error(result.Error))
		return result.Error
	}

	uc.Logger.Info("Current status finalized",
		zap.String("orderID", orderID.String()),
		zap.String("status", currentStatus.Status),
		zap.Time("startedAt", currentStatus.StartedAt),
		zap.Time("endedAt", now),
		zap.Int("durationMinutes", duration))

	return nil
}

// Inicia um novo status
func (uc *ManageOrderStatusHistory) StartNewStatus(orderID uuid.UUID, status string) error {
	now := time.Now()

	newStatusHistory := &models.OrderStatusHistory{
		ID:        uuid.New(),
		OrderID:   orderID,
		Status:    status,
		StartedAt: now,
	}

	result := uc.DB.Create(newStatusHistory)
	if result.Error != nil {
		uc.Logger.Error("Error creating new status history", zap.Error(result.Error))
		return result.Error
	}

	uc.Logger.Info("New status started",
		zap.String("orderID", orderID.String()),
		zap.String("status", status),
		zap.Time("startedAt", now))

	return nil
}

// Busca o histórico completo de uma order
func (uc *ManageOrderStatusHistory) GetOrderHistory(orderID uuid.UUID) ([]models.OrderStatusHistory, error) {
	var history []models.OrderStatusHistory

	if err := uc.DB.Where("order_id = ? ORDER BY started_at ASC", orderID).Find(&history).Error; err != nil {
		uc.Logger.Error("Error fetching order history", zap.Error(err))
		return nil, err
	}

	uc.Logger.Info("Order history retrieved",
		zap.String("orderID", orderID.String()),
		zap.Int("historyCount", len(history)))

	return history, nil
}
