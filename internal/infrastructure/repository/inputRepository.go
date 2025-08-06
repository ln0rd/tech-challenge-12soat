package repository

import (
	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"gorm.io/gorm"
)

// InputRepository define a interface para operações de input no banco
type InputRepository interface {
	Create(input *models.Input) error
	FindByID(id uuid.UUID) (*models.Input, error)
	FindAll() ([]models.Input, error)
	FindByName(name string) (*models.Input, error)
	Update(input *models.Input) error
	Delete(id uuid.UUID) error
}

// InputRepositoryAdapter implementa InputRepository usando GORM
type InputRepositoryAdapter struct {
	db *gorm.DB
}

// NewInputRepositoryAdapter cria uma nova instância do adaptador
func NewInputRepositoryAdapter(db *gorm.DB) InputRepository {
	return &InputRepositoryAdapter{
		db: db,
	}
}

// Create implementa a criação de um input
func (i *InputRepositoryAdapter) Create(input *models.Input) error {
	result := i.db.Create(input)
	return result.Error
}

// FindByID implementa a busca de input por ID
func (i *InputRepositoryAdapter) FindByID(id uuid.UUID) (*models.Input, error) {
	var input models.Input
	result := i.db.Where("id = ?", id).First(&input)
	if result.Error != nil {
		return nil, result.Error
	}
	return &input, nil
}

// FindAll implementa a busca de todos os inputs
func (i *InputRepositoryAdapter) FindAll() ([]models.Input, error) {
	var inputs []models.Input
	result := i.db.Find(&inputs)
	if result.Error != nil {
		return nil, result.Error
	}
	return inputs, nil
}

// FindByName implementa a busca de input por nome
func (i *InputRepositoryAdapter) FindByName(name string) (*models.Input, error) {
	var input models.Input
	result := i.db.Where("name = ?", name).First(&input)
	if result.Error != nil {
		return nil, result.Error
	}
	return &input, nil
}

// Update implementa a atualização de um input
func (i *InputRepositoryAdapter) Update(input *models.Input) error {
	result := i.db.Model(input).Updates(input)
	return result.Error
}

// Delete implementa a exclusão de um input
func (i *InputRepositoryAdapter) Delete(id uuid.UUID) error {
	result := i.db.Where("id = ?", id).Delete(&models.Input{})
	return result.Error
}
