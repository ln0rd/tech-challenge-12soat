package input

import (
	"time"

	"github.com/google/uuid"
)

type Input struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	InputType   string    `json:"input_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
