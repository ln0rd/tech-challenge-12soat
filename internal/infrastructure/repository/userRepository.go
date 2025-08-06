package repository

import (
	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"gorm.io/gorm"
)

// UserRepository define a interface para operações de user no banco
type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uuid.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
}

// UserRepositoryAdapter implementa UserRepository usando GORM
type UserRepositoryAdapter struct {
	db *gorm.DB
}

// NewUserRepositoryAdapter cria uma nova instância do adaptador
func NewUserRepositoryAdapter(db *gorm.DB) UserRepository {
	return &UserRepositoryAdapter{
		db: db,
	}
}

// Create implementa a criação de um user
func (u *UserRepositoryAdapter) Create(user *models.User) error {
	result := u.db.Create(user)
	return result.Error
}

// FindByID implementa a busca de user por ID
func (u *UserRepositoryAdapter) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	result := u.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindByEmail implementa a busca de user por email
func (u *UserRepositoryAdapter) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := u.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Update implementa a atualização de um user
func (u *UserRepositoryAdapter) Update(user *models.User) error {
	result := u.db.Model(user).Updates(user)
	return result.Error
}

// Delete implementa a exclusão de um user
func (u *UserRepositoryAdapter) Delete(id uuid.UUID) error {
	result := u.db.Where("id = ?", id).Delete(&models.User{})
	return result.Error
}
