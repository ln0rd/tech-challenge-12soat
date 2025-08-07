package mocks

import (
	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
)

// OrderStatusHistoryRepositoryMock implementa OrderStatusHistoryRepository para testes
type OrderStatusHistoryRepositoryMock struct {
	CreateFunc               func(statusHistory *models.OrderStatusHistory) error
	FindByIDFunc             func(id uuid.UUID) (*models.OrderStatusHistory, error)
	FindByOrderIDFunc        func(orderID uuid.UUID) ([]models.OrderStatusHistory, error)
	FindCurrentByOrderIDFunc func(orderID uuid.UUID) (*models.OrderStatusHistory, error)
	UpdateFunc               func(statusHistory *models.OrderStatusHistory) error
	DeleteFunc               func(id uuid.UUID) error
}

// Create chama a função mock
func (m *OrderStatusHistoryRepositoryMock) Create(statusHistory *models.OrderStatusHistory) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(statusHistory)
	}
	return nil
}

// FindByID chama a função mock
func (m *OrderStatusHistoryRepositoryMock) FindByID(id uuid.UUID) (*models.OrderStatusHistory, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

// FindByOrderID chama a função mock
func (m *OrderStatusHistoryRepositoryMock) FindByOrderID(orderID uuid.UUID) ([]models.OrderStatusHistory, error) {
	if m.FindByOrderIDFunc != nil {
		return m.FindByOrderIDFunc(orderID)
	}
	return nil, nil
}

// FindCurrentByOrderID chama a função mock
func (m *OrderStatusHistoryRepositoryMock) FindCurrentByOrderID(orderID uuid.UUID) (*models.OrderStatusHistory, error) {
	if m.FindCurrentByOrderIDFunc != nil {
		return m.FindCurrentByOrderIDFunc(orderID)
	}
	return nil, nil
}

// Update chama a função mock
func (m *OrderStatusHistoryRepositoryMock) Update(statusHistory *models.OrderStatusHistory) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(statusHistory)
	}
	return nil
}

// Delete chama a função mock
func (m *OrderStatusHistoryRepositoryMock) Delete(id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
