package models

import (
	"time"

	"github.com/google/uuid"
)

type Supplie struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name              string    `json:"name" gorm:"not null"`
	Description       string    `json:"description"`
	Price             float64   `json:"price" gorm:"not null"`
	QuantityAvailable int       `json:"quantity_available" gorm:"not null;default:0"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
