package order_supplies

import (
	"time"

	"github.com/google/uuid"
)

type OrderSupplie struct {
	ID         uuid.UUID `json:"id"`
	OrderID    uuid.UUID `json:"order_id"`
	SupplyID   uuid.UUID `json:"supply_id"`
	Quantity   int       `json:"quantity"`
	TotalValue float64   `json:"total_value"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
