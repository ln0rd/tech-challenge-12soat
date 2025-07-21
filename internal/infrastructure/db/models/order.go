package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CustomerID uuid.UUID `json:"customer_id" gorm:"type:uuid"`
	VehicleID  uuid.UUID `json:"vehicle_id" gorm:"type:uuid"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (o *Order) TableName() string {
	return "orders"
}
