package models

type User struct {
	ID         int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Email      string `json:"email" gorm:"not null"`
	Password   string `json:"password" gorm:"not null"`
	Username   string `json:"username" gorm:"not null"`
	CustomerID string `json:"customer_id" gorm:"not null"`
	CreatedAt  string `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  string `json:"updated_at" gorm:"autoUpdateTime"`
}

func (u *User) TableName() string {
	return "users"
}
