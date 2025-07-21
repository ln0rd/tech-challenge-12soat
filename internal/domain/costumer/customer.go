package domain

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	UserID         string    `json:"user_id"`
	DocumentNumber string    `json:"document_number"`
	CustomerType   string    `json:"customer_type"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
