package mocks

import (
	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
)

// OrderStatusHistoryManagerMock implementa ManageOrderStatusHistory para testes
type OrderStatusHistoryManagerMock struct {
	IsFinalStatusFunc              func(status string) bool
	FetchCurrentStatusFromDBFunc   func(orderID uuid.UUID) (*models.OrderStatusHistory, error)
	FinalizeCurrentStatusFunc      func(orderID uuid.UUID) error
	CreateNewStatusFunc            func(orderID uuid.UUID, status string) error
	UpdateCurrentStatusToFinalFunc func(orderID uuid.UUID, finalStatus string) error
	FetchOrderHistoryFromDBFunc    func(orderID uuid.UUID) ([]models.OrderStatusHistory, error)
	StartNewStatusFunc             func(orderID uuid.UUID, status string) error
	UpdateStatusFunc               func(orderID uuid.UUID, newStatus string) error
	GetOrderHistoryFunc            func(orderID uuid.UUID) ([]models.OrderStatusHistory, error)
}

// IsFinalStatus chama a função mock
func (m *OrderStatusHistoryManagerMock) IsFinalStatus(status string) bool {
	if m.IsFinalStatusFunc != nil {
		return m.IsFinalStatusFunc(status)
	}
	return false
}

// FetchCurrentStatusFromDB chama a função mock
func (m *OrderStatusHistoryManagerMock) FetchCurrentStatusFromDB(orderID uuid.UUID) (*models.OrderStatusHistory, error) {
	if m.FetchCurrentStatusFromDBFunc != nil {
		return m.FetchCurrentStatusFromDBFunc(orderID)
	}
	return nil, nil
}

// FinalizeCurrentStatus chama a função mock
func (m *OrderStatusHistoryManagerMock) FinalizeCurrentStatus(orderID uuid.UUID) error {
	if m.FinalizeCurrentStatusFunc != nil {
		return m.FinalizeCurrentStatusFunc(orderID)
	}
	return nil
}

// CreateNewStatus chama a função mock
func (m *OrderStatusHistoryManagerMock) CreateNewStatus(orderID uuid.UUID, status string) error {
	if m.CreateNewStatusFunc != nil {
		return m.CreateNewStatusFunc(orderID, status)
	}
	return nil
}

// UpdateCurrentStatusToFinal chama a função mock
func (m *OrderStatusHistoryManagerMock) UpdateCurrentStatusToFinal(orderID uuid.UUID, finalStatus string) error {
	if m.UpdateCurrentStatusToFinalFunc != nil {
		return m.UpdateCurrentStatusToFinalFunc(orderID, finalStatus)
	}
	return nil
}

// FetchOrderHistoryFromDB chama a função mock
func (m *OrderStatusHistoryManagerMock) FetchOrderHistoryFromDB(orderID uuid.UUID) ([]models.OrderStatusHistory, error) {
	if m.FetchOrderHistoryFromDBFunc != nil {
		return m.FetchOrderHistoryFromDBFunc(orderID)
	}
	return nil, nil
}

// StartNewStatus chama a função mock
func (m *OrderStatusHistoryManagerMock) StartNewStatus(orderID uuid.UUID, status string) error {
	if m.StartNewStatusFunc != nil {
		return m.StartNewStatusFunc(orderID, status)
	}
	return nil
}

// UpdateStatus chama a função mock
func (m *OrderStatusHistoryManagerMock) UpdateStatus(orderID uuid.UUID, newStatus string) error {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(orderID, newStatus)
	}
	return nil
}

// GetOrderHistory chama a função mock
func (m *OrderStatusHistoryManagerMock) GetOrderHistory(orderID uuid.UUID) ([]models.OrderStatusHistory, error) {
	if m.GetOrderHistoryFunc != nil {
		return m.GetOrderHistoryFunc(orderID)
	}
	return nil, nil
}
