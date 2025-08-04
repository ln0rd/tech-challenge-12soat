package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email      string     `json:"email" gorm:"not null;unique"`
	Password   string     `json:"password" gorm:"not null"`
	Username   string     `json:"username" gorm:"not null"`
	UserType   string     `json:"user_type" gorm:"not null;check:user_type IN ('admin', 'mechanic', 'vehicle_owner')"`
	CustomerID *uuid.UUID `json:"customer_id" gorm:"type:uuid"`
	CreatedAt  time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (u *User) TableName() string {
	return "users"
}
