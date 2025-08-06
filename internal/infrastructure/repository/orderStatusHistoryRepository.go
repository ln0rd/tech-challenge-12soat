package repository

import (
	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"gorm.io/gorm"
)

// OrderStatusHistoryRepository define a interface para operações de order_status_history no banco
type OrderStatusHistoryRepository interface {
	Create(orderStatusHistory *models.OrderStatusHistory) error
	FindByID(id uuid.UUID) (*models.OrderStatusHistory, error)
	FindByOrderID(orderID uuid.UUID) ([]models.OrderStatusHistory, error)
	FindCurrentByOrderID(orderID uuid.UUID) (*models.OrderStatusHistory, error)
	Update(orderStatusHistory *models.OrderStatusHistory) error
	Delete(id uuid.UUID) error
}

// OrderStatusHistoryRepositoryAdapter implementa OrderStatusHistoryRepository usando GORM
type OrderStatusHistoryRepositoryAdapter struct {
	db *gorm.DB
}

// NewOrderStatusHistoryRepositoryAdapter cria uma nova instância do adaptador
func NewOrderStatusHistoryRepositoryAdapter(db *gorm.DB) OrderStatusHistoryRepository {
	return &OrderStatusHistoryRepositoryAdapter{
		db: db,
	}
}

// Create implementa a criação de um order_status_history
func (osh *OrderStatusHistoryRepositoryAdapter) Create(orderStatusHistory *models.OrderStatusHistory) error {
	result := osh.db.Create(orderStatusHistory)
	return result.Error
}

// FindByID implementa a busca de order_status_history por ID
func (osh *OrderStatusHistoryRepositoryAdapter) FindByID(id uuid.UUID) (*models.OrderStatusHistory, error) {
	var orderStatusHistory models.OrderStatusHistory
	result := osh.db.Where("id = ?", id).First(&orderStatusHistory)
	if result.Error != nil {
		return nil, result.Error
	}
	return &orderStatusHistory, nil
}

// FindByOrderID implementa a busca de order_status_history por order ID
func (osh *OrderStatusHistoryRepositoryAdapter) FindByOrderID(orderID uuid.UUID) ([]models.OrderStatusHistory, error) {
	var orderStatusHistories []models.OrderStatusHistory
	result := osh.db.Where("order_id = ?", orderID).Find(&orderStatusHistories)
	if result.Error != nil {
		return nil, result.Error
	}
	return orderStatusHistories, nil
}

// FindCurrentByOrderID implementa a busca do status atual de um order
func (osh *OrderStatusHistoryRepositoryAdapter) FindCurrentByOrderID(orderID uuid.UUID) (*models.OrderStatusHistory, error) {
	var orderStatusHistory models.OrderStatusHistory
	result := osh.db.Where("order_id = ? AND ended_at IS NULL", orderID).First(&orderStatusHistory)
	if result.Error != nil {
		return nil, result.Error
	}
	return &orderStatusHistory, nil
}

// Update implementa a atualização de um order_status_history
func (osh *OrderStatusHistoryRepositoryAdapter) Update(orderStatusHistory *models.OrderStatusHistory) error {
	result := osh.db.Model(orderStatusHistory).Updates(orderStatusHistory)
	return result.Error
}

// Delete implementa a exclusão de um order_status_history
func (osh *OrderStatusHistoryRepositoryAdapter) Delete(id uuid.UUID) error {
	result := osh.db.Where("id = ?", id).Delete(&models.OrderStatusHistory{})
	return result.Error
}
