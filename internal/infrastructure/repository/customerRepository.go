package repository

import (
	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"gorm.io/gorm"
)

// CustomerRepository define a interface para operações de customer no banco
type CustomerRepository interface {
	Create(customer *models.Customer) error
	FindByID(id uuid.UUID) (*models.Customer, error)
	FindAll() ([]models.Customer, error)
	Update(customer *models.Customer) error
	Delete(id uuid.UUID) error
}

// CustomerRepositoryAdapter implementa CustomerRepository usando GORM
type CustomerRepositoryAdapter struct {
	db *gorm.DB
}

// NewCustomerRepositoryAdapter cria uma nova instância do adaptador
func NewCustomerRepositoryAdapter(db *gorm.DB) CustomerRepository {
	return &CustomerRepositoryAdapter{
		db: db,
	}
}

// Create implementa a criação de um customer
func (c *CustomerRepositoryAdapter) Create(customer *models.Customer) error {
	result := c.db.Create(customer)
	return result.Error
}

// FindByID implementa a busca de customer por ID
func (c *CustomerRepositoryAdapter) FindByID(id uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	result := c.db.Where("id = ?", id).First(&customer)
	if result.Error != nil {
		return nil, result.Error
	}
	return &customer, nil
}

// FindAll implementa a busca de todos os customers
func (c *CustomerRepositoryAdapter) FindAll() ([]models.Customer, error) {
	var customers []models.Customer
	result := c.db.Find(&customers)
	if result.Error != nil {
		return nil, result.Error
	}
	return customers, nil
}

// Update implementa a atualização de um customer
func (c *CustomerRepositoryAdapter) Update(customer *models.Customer) error {
	result := c.db.Model(customer).Updates(customer)
	return result.Error
}

// Delete implementa a exclusão de um customer
func (c *CustomerRepositoryAdapter) Delete(id uuid.UUID) error {
	result := c.db.Where("id = ?", id).Delete(&models.Customer{})
	return result.Error
}
