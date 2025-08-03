package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatusHistory struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OrderID         uuid.UUID  `json:"order_id" gorm:"type:uuid;not null"`
	Status          string     `json:"status" gorm:"not null"`
	StartedAt       time.Time  `json:"started_at" gorm:"not null"`
	EndedAt         *time.Time `json:"ended_at"`
	DurationMinutes *int       `json:"duration_minutes"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

func (osh *OrderStatusHistory) TableName() string {
	return "order_status_history"
}
