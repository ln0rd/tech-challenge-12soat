package mocks

import (
	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
)

// VehicleRepositoryMock implementa VehicleRepository para testes
type VehicleRepositoryMock struct {
	CreateFunc            func(vehicle *models.Vehicle) error
	FindByIDFunc          func(id uuid.UUID) (*models.Vehicle, error)
	FindByCustomerIDFunc  func(customerID uuid.UUID) ([]models.Vehicle, error)
	FindByNumberPlateFunc func(numberPlate string) (*models.Vehicle, error)
	UpdateFunc            func(vehicle *models.Vehicle) error
	DeleteFunc            func(id uuid.UUID) error
}

// Create chama a função mock
func (m *VehicleRepositoryMock) Create(vehicle *models.Vehicle) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(vehicle)
	}
	return nil
}

// FindByID chama a função mock
func (m *VehicleRepositoryMock) FindByID(id uuid.UUID) (*models.Vehicle, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

// FindByCustomerID chama a função mock
func (m *VehicleRepositoryMock) FindByCustomerID(customerID uuid.UUID) ([]models.Vehicle, error) {
	if m.FindByCustomerIDFunc != nil {
		return m.FindByCustomerIDFunc(customerID)
	}
	return nil, nil
}

// FindByNumberPlate chama a função mock
func (m *VehicleRepositoryMock) FindByNumberPlate(numberPlate string) (*models.Vehicle, error) {
	if m.FindByNumberPlateFunc != nil {
		return m.FindByNumberPlateFunc(numberPlate)
	}
	return nil, nil
}

// Update chama a função mock
func (m *VehicleRepositoryMock) Update(vehicle *models.Vehicle) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(vehicle)
	}
	return nil
}

// Delete chama a função mock
func (m *VehicleRepositoryMock) Delete(id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
