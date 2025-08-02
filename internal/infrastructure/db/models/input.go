package models

import (
	"time"

	"github.com/google/uuid"
)

type Input struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `json:"name" gorm:"not null;unique"`
	Description string    `json:"description"`
	Price       float64   `json:"price" gorm:"not null"`
	Quantity    int       `json:"quantity" gorm:"not null"`
	InputType   string    `json:"input_type" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (i *Input) TableName() string {
	return "inputs"
}
