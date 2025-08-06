package mocks

import (
	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
)

// CustomerRepositoryMock implementa CustomerRepository para testes
type CustomerRepositoryMock struct {
	CreateFunc   func(customer *models.Customer) error
	FindByIDFunc func(id uuid.UUID) (*models.Customer, error)
	FindAllFunc  func() ([]models.Customer, error)
	UpdateFunc   func(customer *models.Customer) error
	DeleteFunc   func(id uuid.UUID) error
}

// Create chama a função mock
func (m *CustomerRepositoryMock) Create(customer *models.Customer) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(customer)
	}
	return nil
}

// FindByID chama a função mock
func (m *CustomerRepositoryMock) FindByID(id uuid.UUID) (*models.Customer, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

// FindAll chama a função mock
func (m *CustomerRepositoryMock) FindAll() ([]models.Customer, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc()
	}
	return nil, nil
}

// Update chama a função mock
func (m *CustomerRepositoryMock) Update(customer *models.Customer) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(customer)
	}
	return nil
}

// Delete chama a função mock
func (m *CustomerRepositoryMock) Delete(id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
