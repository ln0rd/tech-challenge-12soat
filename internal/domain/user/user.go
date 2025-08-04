package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID  `json:"id"`
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	Username   string     `json:"username"`
	UserType   string     `json:"user_type"` // admin, mechanic, vehicle_owner
	CustomerID *uuid.UUID `json:"customer_id,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// Constantes para os tipos de usuário
const (
	UserTypeAdmin        = "admin"
	UserTypeMechanic     = "mechanic"
	UserTypeVehicleOwner = "vehicle_owner"
)
