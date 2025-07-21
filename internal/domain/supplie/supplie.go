package supplie

import (
	"time"

	"github.com/google/uuid"
)

type Supplie struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	Price             float64   `json:"price"`
	QuantityAvailable int       `json:"quantity_available"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
