package customer

import (
	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/costumer"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"gorm.io/gorm"
)

type UpdateByIdCustomer struct {
	DB *gorm.DB
}

func (uc *UpdateByIdCustomer) Process(id uuid.UUID, entity *domain.Customer) error {
	// Primeiro verifica se o customer existe
	var existingCustomer models.Customer
	if err := uc.DB.Where("id = ?", id).First(&existingCustomer).Error; err != nil {
		return err
	}

	// Atualiza os campos do customer existente
	existingCustomer.Name = entity.Name
	existingCustomer.Email = entity.Email
	existingCustomer.UserID = entity.UserID
	existingCustomer.DocumentNumber = entity.DocumentNumber
	existingCustomer.CustomerType = entity.CustomerType

	// Salva as alterações
	return uc.DB.Save(&existingCustomer).Error
}
