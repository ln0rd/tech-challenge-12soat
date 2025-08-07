package mocks

import (
	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
)

// InputRepositoryMock implementa InputRepository para testes
type InputRepositoryMock struct {
	CreateFunc     func(input *models.Input) error
	FindByIDFunc   func(id uuid.UUID) (*models.Input, error)
	FindAllFunc    func() ([]models.Input, error)
	FindByNameFunc func(name string) (*models.Input, error)
	UpdateFunc     func(input *models.Input) error
	DeleteFunc     func(id uuid.UUID) error
}

// Create chama a função mock
func (m *InputRepositoryMock) Create(input *models.Input) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(input)
	}
	return nil
}

// FindByID chama a função mock
func (m *InputRepositoryMock) FindByID(id uuid.UUID) (*models.Input, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

// FindAll chama a função mock
func (m *InputRepositoryMock) FindAll() ([]models.Input, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc()
	}
	return nil, nil
}

// FindByName chama a função mock
func (m *InputRepositoryMock) FindByName(name string) (*models.Input, error) {
	if m.FindByNameFunc != nil {
		return m.FindByNameFunc(name)
	}
	return nil, nil
}

// Update chama a função mock
func (m *InputRepositoryMock) Update(input *models.Input) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(input)
	}
	return nil
}

// Delete chama a função mock
func (m *InputRepositoryMock) Delete(id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
