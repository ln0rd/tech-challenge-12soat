package repository

import (
	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"gorm.io/gorm"
)

// OrderInputRepository define a interface para operações de order_input no banco
type OrderInputRepository interface {
	Create(orderInput *models.OrderInput) error
	FindByID(id uuid.UUID) (*models.OrderInput, error)
	FindByOrderID(orderID uuid.UUID) ([]models.OrderInput, error)
	FindByOrderIDAndInputID(orderID uuid.UUID, inputID uuid.UUID) (*models.OrderInput, error)
	Update(orderInput *models.OrderInput) error
	Delete(id uuid.UUID) error
	DeleteByOrderIDAndInputID(orderID uuid.UUID, inputID uuid.UUID) error
}

// OrderInputRepositoryAdapter implementa OrderInputRepository usando GORM
type OrderInputRepositoryAdapter struct {
	db *gorm.DB
}

// NewOrderInputRepositoryAdapter cria uma nova instância do adaptador
func NewOrderInputRepositoryAdapter(db *gorm.DB) OrderInputRepository {
	return &OrderInputRepositoryAdapter{
		db: db,
	}
}

// Create implementa a criação de um order_input
func (oi *OrderInputRepositoryAdapter) Create(orderInput *models.OrderInput) error {
	result := oi.db.Create(orderInput)
	return result.Error
}

// FindByID implementa a busca de order_input por ID
func (oi *OrderInputRepositoryAdapter) FindByID(id uuid.UUID) (*models.OrderInput, error) {
	var orderInput models.OrderInput
	result := oi.db.Where("id = ?", id).First(&orderInput)
	if result.Error != nil {
		return nil, result.Error
	}
	return &orderInput, nil
}

// FindByOrderID implementa a busca de order_inputs por order ID
func (oi *OrderInputRepositoryAdapter) FindByOrderID(orderID uuid.UUID) ([]models.OrderInput, error) {
	var orderInputs []models.OrderInput
	result := oi.db.Where("order_id = ?", orderID).Find(&orderInputs)
	if result.Error != nil {
		return nil, result.Error
	}
	return orderInputs, nil
}

// FindByOrderIDAndInputID implementa a busca de order_input por order ID e input ID
func (oi *OrderInputRepositoryAdapter) FindByOrderIDAndInputID(orderID uuid.UUID, inputID uuid.UUID) (*models.OrderInput, error) {
	var orderInput models.OrderInput
	result := oi.db.Where("order_id = ? AND input_id = ?", orderID, inputID).First(&orderInput)
	if result.Error != nil {
		return nil, result.Error
	}
	return &orderInput, nil
}

// Update implementa a atualização de um order_input
func (oi *OrderInputRepositoryAdapter) Update(orderInput *models.OrderInput) error {
	result := oi.db.Model(orderInput).Updates(orderInput)
	return result.Error
}

// Delete implementa a exclusão de um order_input
func (oi *OrderInputRepositoryAdapter) Delete(id uuid.UUID) error {
	result := oi.db.Where("id = ?", id).Delete(&models.OrderInput{})
	return result.Error
}

// DeleteByOrderIDAndInputID implementa a exclusão de order_input por order ID e input ID
func (oi *OrderInputRepositoryAdapter) DeleteByOrderIDAndInputID(orderID uuid.UUID, inputID uuid.UUID) error {
	result := oi.db.Where("order_id = ? AND input_id = ?", orderID, inputID).Delete(&models.OrderInput{})
	return result.Error
}
