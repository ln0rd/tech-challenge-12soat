package repository

import (
	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"gorm.io/gorm"
)

// VehicleRepository define a interface para operações de vehicle no banco
type VehicleRepository interface {
	Create(vehicle *models.Vehicle) error
	FindByID(id uuid.UUID) (*models.Vehicle, error)
	FindByCustomerID(customerID uuid.UUID) ([]models.Vehicle, error)
	FindByNumberPlate(numberPlate string) (*models.Vehicle, error)
	Update(vehicle *models.Vehicle) error
	Delete(id uuid.UUID) error
}

// VehicleRepositoryAdapter implementa VehicleRepository usando GORM
type VehicleRepositoryAdapter struct {
	db *gorm.DB
}

// NewVehicleRepositoryAdapter cria uma nova instância do adaptador
func NewVehicleRepositoryAdapter(db *gorm.DB) VehicleRepository {
	return &VehicleRepositoryAdapter{
		db: db,
	}
}

// Create implementa a criação de um vehicle
func (v *VehicleRepositoryAdapter) Create(vehicle *models.Vehicle) error {
	result := v.db.Create(vehicle)
	return result.Error
}

// FindByID implementa a busca de vehicle por ID
func (v *VehicleRepositoryAdapter) FindByID(id uuid.UUID) (*models.Vehicle, error) {
	var vehicle models.Vehicle
	result := v.db.Where("id = ?", id).First(&vehicle)
	if result.Error != nil {
		return nil, result.Error
	}
	return &vehicle, nil
}

// FindByCustomerID implementa a busca de vehicles por customer ID
func (v *VehicleRepositoryAdapter) FindByCustomerID(customerID uuid.UUID) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	result := v.db.Where("customer_id = ?", customerID).Find(&vehicles)
	if result.Error != nil {
		return nil, result.Error
	}
	return vehicles, nil
}

// FindByNumberPlate implementa a busca de vehicle por placa
func (v *VehicleRepositoryAdapter) FindByNumberPlate(numberPlate string) (*models.Vehicle, error) {
	var vehicle models.Vehicle
	result := v.db.Where("number_plate = ?", numberPlate).First(&vehicle)
	if result.Error != nil {
		return nil, result.Error
	}
	return &vehicle, nil
}

// Update implementa a atualização de um vehicle
func (v *VehicleRepositoryAdapter) Update(vehicle *models.Vehicle) error {
	result := v.db.Model(vehicle).Updates(vehicle)
	return result.Error
}

// Delete implementa a exclusão de um vehicle
func (v *VehicleRepositoryAdapter) Delete(id uuid.UUID) error {
	result := v.db.Where("id = ?", id).Delete(&models.Vehicle{})
	return result.Error
}
