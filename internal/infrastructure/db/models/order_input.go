package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderInput struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OrderID    uuid.UUID `json:"order_id" gorm:"type:uuid;not null"`
	InputID    uuid.UUID `json:"input_id" gorm:"type:uuid;not null"`
	Quantity   int       `json:"quantity" gorm:"not null"`
	UnitPrice  float64   `json:"unit_price" gorm:"not null"`
	TotalPrice float64   `json:"total_price" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (oi *OrderInput) TableName() string {
	return "order_inputs"
}
