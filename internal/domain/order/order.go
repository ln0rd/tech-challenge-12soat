package order

import (
	"time"

	"github.com/google/uuid"
)

const (
	StatusReceived            = "Received"
	StatusUndergoingDiagnosis = "Undergoing diagnosis"
	StatusAwaitingApproval    = "Awaiting approval"
	StatusInProgress          = "In progress"
	StatusCompleted           = "Completed"
	StatusDelivered           = "Delivered"
)

type Order struct {
	ID         uuid.UUID `json:"id"`
	CustomerID uuid.UUID `json:"customer_id"`
	VehicleID  uuid.UUID `json:"vehicle_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
