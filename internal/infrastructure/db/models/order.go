package models

type Order struct {
	ID         int    `json:"id" gorm:"primaryKey;autoIncrement"`
	CustomerID int    `json:"customer_id" gorm:"not null"`
	VehicleID  int    `json:"vehicle_id" gorm:"not null"`
	CreatedAt  string `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  string `json:"updated_at" gorm:"autoUpdateTime"`
}

func (o *Order) TableName() string {
	return "orders"
}
