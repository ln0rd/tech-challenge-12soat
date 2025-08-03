package order_status_history

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatusHistory struct {
	ID              uuid.UUID  `json:"id"`
	OrderID         uuid.UUID  `json:"order_id"`
	Status          string     `json:"status"`
	StartedAt       time.Time  `json:"started_at"`
	EndedAt         *time.Time `json:"ended_at,omitempty"`
	DurationMinutes *int       `json:"duration_minutes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}
