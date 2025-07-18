package order_supplies

type OrderSupplie struct {
	ID         int     `json:"id"`
	OrderID    int     `json:"order_id"`
	SupplyID   int     `json:"supply_id"`
	Quantity   int     `json:"quantity"`
	TotalValue float64 `json:"total_value"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}
