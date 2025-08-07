package mocks

import (
	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
)

// OrderInputRepositoryMock implementa OrderInputRepository para testes
type OrderInputRepositoryMock struct {
	CreateFunc                    func(orderInput *models.OrderInput) error
	FindByIDFunc                  func(id uuid.UUID) (*models.OrderInput, error)
	FindByOrderIDFunc             func(orderID uuid.UUID) ([]models.OrderInput, error)
	FindByOrderIDAndInputIDFunc   func(orderID uuid.UUID, inputID uuid.UUID) (*models.OrderInput, error)
	UpdateFunc                    func(orderInput *models.OrderInput) error
	DeleteFunc                    func(id uuid.UUID) error
	DeleteByOrderIDAndInputIDFunc func(orderID uuid.UUID, inputID uuid.UUID) error
}

// Create chama a função mock
func (m *OrderInputRepositoryMock) Create(orderInput *models.OrderInput) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(orderInput)
	}
	return nil
}

// FindByID chama a função mock
func (m *OrderInputRepositoryMock) FindByID(id uuid.UUID) (*models.OrderInput, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

// FindByOrderID chama a função mock
func (m *OrderInputRepositoryMock) FindByOrderID(orderID uuid.UUID) ([]models.OrderInput, error) {
	if m.FindByOrderIDFunc != nil {
		return m.FindByOrderIDFunc(orderID)
	}
	return nil, nil
}

// FindByOrderIDAndInputID chama a função mock
func (m *OrderInputRepositoryMock) FindByOrderIDAndInputID(orderID uuid.UUID, inputID uuid.UUID) (*models.OrderInput, error) {
	if m.FindByOrderIDAndInputIDFunc != nil {
		return m.FindByOrderIDAndInputIDFunc(orderID, inputID)
	}
	return nil, nil
}

// Update chama a função mock
func (m *OrderInputRepositoryMock) Update(orderInput *models.OrderInput) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(orderInput)
	}
	return nil
}

// Delete chama a função mock
func (m *OrderInputRepositoryMock) Delete(id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

// DeleteByOrderIDAndInputID chama a função mock
func (m *OrderInputRepositoryMock) DeleteByOrderIDAndInputID(orderID uuid.UUID, inputID uuid.UUID) error {
	if m.DeleteByOrderIDAndInputIDFunc != nil {
		return m.DeleteByOrderIDAndInputIDFunc(orderID, inputID)
	}
	return nil
}
