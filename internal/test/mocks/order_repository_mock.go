package mocks

import (
	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
)

// OrderRepositoryMock implementa OrderRepository para testes
type OrderRepositoryMock struct {
	CreateFunc   func(order *models.Order) error
	FindByIDFunc func(id uuid.UUID) (*models.Order, error)
	FindAllFunc  func() ([]models.Order, error)
	UpdateFunc   func(order *models.Order) error
	DeleteFunc   func(id uuid.UUID) error
}

// Create chama a função mock
func (m *OrderRepositoryMock) Create(order *models.Order) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(order)
	}
	return nil
}

// FindByID chama a função mock
func (m *OrderRepositoryMock) FindByID(id uuid.UUID) (*models.Order, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

// FindAll chama a função mock
func (m *OrderRepositoryMock) FindAll() ([]models.Order, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc()
	}
	return nil, nil
}

// Update chama a função mock
func (m *OrderRepositoryMock) Update(order *models.Order) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(order)
	}
	return nil
}

// Delete chama a função mock
func (m *OrderRepositoryMock) Delete(id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
