package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderSupplie struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OrderID    uuid.UUID `json:"order_id" gorm:"type:uuid"`
	SupplyID   uuid.UUID `json:"supply_id" gorm:"type:uuid"`
	Quantity   int       `json:"quantity" gorm:"not null;default:1"`
	TotalValue float64   `json:"total_value" gorm:"type:numeric(10,2);not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (o *OrderSupplie) TableName() string {
	return "order_supplies"
}
