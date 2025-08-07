package mocks

import (
	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
)

// UserRepositoryMock implementa UserRepository para testes
type UserRepositoryMock struct {
	CreateFunc      func(user *models.User) error
	FindByIDFunc    func(id uuid.UUID) (*models.User, error)
	FindByEmailFunc func(email string) (*models.User, error)
	UpdateFunc      func(user *models.User) error
	DeleteFunc      func(id uuid.UUID) error
}

// Create chama a função mock
func (m *UserRepositoryMock) Create(user *models.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(user)
	}
	return nil
}

// FindByID chama a função mock
func (m *UserRepositoryMock) FindByID(id uuid.UUID) (*models.User, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

// FindByEmail chama a função mock
func (m *UserRepositoryMock) FindByEmail(email string) (*models.User, error) {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(email)
	}
	return nil, nil
}

// Update chama a função mock
func (m *UserRepositoryMock) Update(user *models.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(user)
	}
	return nil
}

// Delete chama a função mock
func (m *UserRepositoryMock) Delete(id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
