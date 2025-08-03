package order

import (
	"time"

	"github.com/google/uuid"
)

const (
	StatusReceived = "Received"
	// ------------------------------ Mecanico abre OS
	StatusUndergoingDiagnosis = "Undergoing diagnosis"
	// ------------------------------ Cliente aprova OS
	StatusAwaitingApproval = "Awaiting approval"
	// ------------------------------ Mecanico inicia serviço
	StatusInProgress = "In progress"
	// ------------------------------ Mecanico finaliza serviço
	StatusCompleted = "Completed"
	// ------------------------------ Cliente retira veículo
	StatusDelivered = "Delivered"
	// ------------------------------ Order cancelada
	StatusCanceled = "Canceled"
)

type Order struct {
	ID         uuid.UUID `json:"id"`
	CustomerID uuid.UUID `json:"customer_id"`
	VehicleID  uuid.UUID `json:"vehicle_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
