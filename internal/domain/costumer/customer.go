package domain

type Customer struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	UserID         int    `json:"user_id"`
	DocumentNumber string `json:"document_number"`
	CustomerType   string `json:"customer_type"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
