package models

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name           string    `json:"name" gorm:"not null"`
	Email          string    `json:"email" gorm:"not null"`
	UserID         string    `json:"user_id"`
	DocumentNumber string    `json:"document_number" gorm:"not null"`
	CustomerType   string    `json:"customer_type" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (c *Customer) TableName() string {
	return "customers"
}
