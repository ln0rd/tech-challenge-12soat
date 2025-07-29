package customer

import (
	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"gorm.io/gorm"
)

type DeleteByIdCustomer struct {
	DB *gorm.DB
}

func (uc *DeleteByIdCustomer) Process(id uuid.UUID) error {
	result := uc.DB.Where("id = ?", id).Delete(&models.Customer{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
