package user

type User struct {
	ID         int    `json:"id"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Username   string `json:"username"`
	CustomerID string `json:"customer_id"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
