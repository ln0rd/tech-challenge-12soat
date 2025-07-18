package models

type OrderSupplie struct {
	ID         int     `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID    int     `json:"order_id" gorm:"not null"`
	SupplyID   int     `json:"supply_id" gorm:"not null"`
	Quantity   int     `json:"quantity" gorm:"not null;default:1"`
	TotalValue float64 `json:"total_value" gorm:"type:numeric(10,2);not null"`
	CreatedAt  string  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  string  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (o *OrderSupplie) TableName() string {
	return "order_supplies"
}
