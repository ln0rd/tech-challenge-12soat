package models

type Vehicle struct {
	ID                          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Model                       string `json:"model" gorm:"not null"`
	Brand                       string `json:"brand" gorm:"not null"`
	ReleaseYear                 int    `json:"release_year" gorm:"not null"`
	VehicleIdentificationNumber string `json:"vehicle_identification_number" gorm:"not null"`
	NumberPlate                 string `json:"number_plate" gorm:"not null"`
	Color                       string `json:"color" gorm:"not null"`
	CustomerID                  int    `json:"customer_id" gorm:"not null"`
	CreatedAt                   string `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt                   string `json:"updated_at" gorm:"autoUpdateTime"`
}

func (v *Vehicle) TableName() string {
	return "vehicles"
}
