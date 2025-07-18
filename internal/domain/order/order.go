package order

type Order struct {
	ID         int    `json:"id"`
	CustomerID int    `json:"customer_id"`
	VehicleID  int    `json:"vehicle_id"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
