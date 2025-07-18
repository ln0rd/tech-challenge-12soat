package vehicle

type Vehicle struct {
	ID                          int    `json:"id"`
	Model                       string `json:"model"`
	Brand                       string `json:"brand"`
	ReleaseYear                 int    `json:"release_year"`
	VehicleIdentificationNumber string `json:"vehicle_identification_number"`
	NumberPlate                 string `json:"number_plate"`
	Color                       string `json:"color"`
	CustomerID                  int    `json:"customer_id"`
	CreatedAt                   string `json:"created_at"`
	UpdatedAt                   string `json:"updated_at"`
}
