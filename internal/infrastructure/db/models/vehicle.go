package models

import (
	"time"

	"github.com/google/uuid"
)

type Vehicle struct {
	ID                          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Model                       string    `json:"model" gorm:"not null"`
	Brand                       string    `json:"brand" gorm:"not null"`
	ReleaseYear                 int       `json:"release_year" gorm:"not null"`
	VehicleIdentificationNumber string    `json:"vehicle_identification_number" gorm:"not null"`
	NumberPlate                 string    `json:"number_plate" gorm:"not null;unique"`
	Color                       string    `json:"color" gorm:"not null"`
	CustomerID                  uuid.UUID `json:"customer_id" gorm:"type:uuid;not null"`
	CreatedAt                   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt                   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (v *Vehicle) TableName() string {
	return "vehicles"
}
