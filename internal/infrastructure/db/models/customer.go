package models

type Customer struct {
	ID             int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name           string `json:"name" gorm:"not null"`
	Email          string `json:"email" gorm:"not null"`
	UserID         int    `json:"user_id"`
	DocumentNumber string `json:"document_number" gorm:"not null"`
	CustomerType   string `json:"customer_type" gorm:"not null"`
	CreatedAt      string `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      string `json:"updated_at" gorm:"autoUpdateTime"`
}

func (c *Customer) TableName() string {
	return "customers"
}
