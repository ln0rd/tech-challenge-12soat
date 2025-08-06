package repository

import (
	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"gorm.io/gorm"
)

// OrderRepository define a interface para operações de order no banco
type OrderRepository interface {
	Create(order *models.Order) error
	FindByID(id uuid.UUID) (*models.Order, error)
	FindAll() ([]models.Order, error)
	Update(order *models.Order) error
	Delete(id uuid.UUID) error
}

// OrderRepositoryAdapter implementa OrderRepository usando GORM
type OrderRepositoryAdapter struct {
	db *gorm.DB
}

// NewOrderRepositoryAdapter cria uma nova instância do adaptador
func NewOrderRepositoryAdapter(db *gorm.DB) OrderRepository {
	return &OrderRepositoryAdapter{
		db: db,
	}
}

// Create implementa a criação de um order
func (o *OrderRepositoryAdapter) Create(order *models.Order) error {
	result := o.db.Create(order)
	return result.Error
}

// FindByID implementa a busca de order por ID
func (o *OrderRepositoryAdapter) FindByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	result := o.db.Where("id = ?", id).First(&order)
	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}

// FindAll implementa a busca de todos os orders
func (o *OrderRepositoryAdapter) FindAll() ([]models.Order, error) {
	var orders []models.Order
	result := o.db.Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}

// Update implementa a atualização de um order
func (o *OrderRepositoryAdapter) Update(order *models.Order) error {
	result := o.db.Model(order).Updates(order)
	return result.Error
}

// Delete implementa a exclusão de um order
func (o *OrderRepositoryAdapter) Delete(id uuid.UUID) error {
	result := o.db.Where("id = ?", id).Delete(&models.Order{})
	return result.Error
}
