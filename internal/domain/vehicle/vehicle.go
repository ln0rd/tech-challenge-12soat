package vehicle

import (
	"time"

	"github.com/google/uuid"
)

type Vehicle struct {
	ID                          uuid.UUID `json:"id"`
	Model                       string    `json:"model"`
	Brand                       string    `json:"brand"`
	ReleaseYear                 int       `json:"release_year"`
	VehicleIdentificationNumber string    `json:"vehicle_identification_number"`
	NumberPlate                 string    `json:"number_plate"`
	Color                       string    `json:"color"`
	CustomerID                  uuid.UUID `json:"customer_id"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
}
