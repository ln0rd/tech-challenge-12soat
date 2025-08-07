package order_status_history

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/logger"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ManageOrderStatusHistory struct {
	OrderStatusHistoryRepository repository.OrderStatusHistoryRepository
	Logger                       logger.Logger
}

// IsFinalStatus verifica se o status é final
func (uc *ManageOrderStatusHistory) IsFinalStatus(status string) bool {
	return status == "Delivered" || status == "Canceled"
}

// FetchCurrentStatusFromDB busca o status atual do banco de dados
func (uc *ManageOrderStatusHistory) FetchCurrentStatusFromDB(orderID uuid.UUID) (*models.OrderStatusHistory, error) {
	currentStatus, err := uc.OrderStatusHistoryRepository.FindCurrentByOrderID(orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			uc.Logger.Info("No current status found",
				zap.String("orderID", orderID.String()))
			return nil, nil
		}
		uc.Logger.Error("Error fetching current status", zap.Error(err))
		return nil, err
	}

	if currentStatus != nil {
		uc.Logger.Info("Current status found",
			zap.String("orderID", orderID.String()),
			zap.String("status", currentStatus.Status),
			zap.Time("startedAt", currentStatus.StartedAt))
	}

	return currentStatus, nil
}

// FinalizeCurrentStatus finaliza o status atual
func (uc *ManageOrderStatusHistory) FinalizeCurrentStatus(orderID uuid.UUID) error {
	currentStatus, err := uc.FetchCurrentStatusFromDB(orderID)
	if err != nil {
		return err
	}

	if currentStatus == nil {
		// Não há status atual, isso é normal para a primeira mudança
		uc.Logger.Info("No current status found, this is normal for first status change",
			zap.String("orderID", orderID.String()))
		return nil
	}

	// Calcula a duração
	now := time.Now()

	// Atualiza o registro atual
	currentStatus.EndedAt = &now
	err = uc.OrderStatusHistoryRepository.Update(currentStatus)
	if err != nil {
		uc.Logger.Error("Error updating current status", zap.Error(err))
		return err
	}

	uc.Logger.Info("Current status finalized",
		zap.String("orderID", orderID.String()),
		zap.String("status", currentStatus.Status),
		zap.Time("startedAt", currentStatus.StartedAt),
		zap.Time("endedAt", now))

	return nil
}

// CreateNewStatus cria um novo status
func (uc *ManageOrderStatusHistory) CreateNewStatus(orderID uuid.UUID, status string) error {
	now := time.Now()

	newStatusHistory := &models.OrderStatusHistory{
		ID:        uuid.New(),
		OrderID:   orderID,
		Status:    status,
		StartedAt: now,
	}

	err := uc.OrderStatusHistoryRepository.Create(newStatusHistory)
	if err != nil {
		uc.Logger.Error("Error creating new status history", zap.Error(err))
		return err
	}

	uc.Logger.Info("New status started",
		zap.String("orderID", orderID.String()),
		zap.String("status", status),
		zap.Time("startedAt", now))

	return nil
}

// UpdateCurrentStatusToFinal atualiza o status atual para um status final
func (uc *ManageOrderStatusHistory) UpdateCurrentStatusToFinal(orderID uuid.UUID, finalStatus string) error {
	currentStatus, err := uc.FetchCurrentStatusFromDB(orderID)
	if err != nil {
		return err
	}

	if currentStatus == nil {
		uc.Logger.Error("No current status found for final status update",
			zap.String("orderID", orderID.String()),
			zap.String("finalStatus", finalStatus))
		return errors.New("no current status found")
	}

	uc.Logger.Info("Current status found for final update",
		zap.String("orderID", orderID.String()),
		zap.String("currentStatus", currentStatus.Status),
		zap.String("finalStatus", finalStatus),
		zap.Time("startedAt", currentStatus.StartedAt))

	// Para status finais, ended_at será igual ao started_at (duração = 0)
	now := time.Now()

	// Atualiza o registro atual com o status final
	currentStatus.Status = finalStatus
	currentStatus.EndedAt = &now
	err = uc.OrderStatusHistoryRepository.Update(currentStatus)
	if err != nil {
		uc.Logger.Error("Error updating current status to final", zap.Error(err))
		return err
	}

	uc.Logger.Info("Current status updated to final successfully",
		zap.String("orderID", orderID.String()),
		zap.String("oldStatus", currentStatus.Status),
		zap.String("finalStatus", finalStatus),
		zap.Time("startedAt", currentStatus.StartedAt),
		zap.Time("endedAt", now))

	return nil
}

// FetchOrderHistoryFromDB busca o histórico completo de uma order do banco de dados
func (uc *ManageOrderStatusHistory) FetchOrderHistoryFromDB(orderID uuid.UUID) ([]models.OrderStatusHistory, error) {
	history, err := uc.OrderStatusHistoryRepository.FindByOrderID(orderID)
	if err != nil {
		uc.Logger.Error("Error fetching order history", zap.Error(err))
		return nil, err
	}

	uc.Logger.Info("Order history retrieved",
		zap.String("orderID", orderID.String()),
		zap.Int("historyCount", len(history)))

	return history, nil
}

// StartNewStatus inicia um novo status
func (uc *ManageOrderStatusHistory) StartNewStatus(orderID uuid.UUID, status string) error {
	return uc.CreateNewStatus(orderID, status)
}

// UpdateStatus finaliza o status atual e inicia um novo
func (uc *ManageOrderStatusHistory) UpdateStatus(orderID uuid.UUID, newStatus string) error {
	uc.Logger.Info("Managing order status history",
		zap.String("orderID", orderID.String()),
		zap.String("newStatus", newStatus))

	// Verifica se é um status final
	isFinalStatus := uc.IsFinalStatus(newStatus)

	if isFinalStatus {
		uc.Logger.Info("Status is final, updating current status instead of creating new one",
			zap.String("orderID", orderID.String()),
			zap.String("newStatus", newStatus))

		// Para status finais, apenas atualiza o status atual
		err := uc.UpdateCurrentStatusToFinal(orderID, newStatus)
		if err != nil {
			uc.Logger.Error("Error updating current status to final", zap.Error(err))
			return err
		}
	} else {
		// Finaliza o status atual (se existir)
		err := uc.FinalizeCurrentStatus(orderID)
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
	}

	uc.Logger.Info("Order status history updated successfully",
		zap.String("orderID", orderID.String()),
		zap.String("newStatus", newStatus),
		zap.Bool("isFinalStatus", isFinalStatus))

	return nil
}

// GetOrderHistory busca o histórico completo de uma order
func (uc *ManageOrderStatusHistory) GetOrderHistory(orderID uuid.UUID) ([]models.OrderStatusHistory, error) {
	return uc.FetchOrderHistoryFromDB(orderID)
}
