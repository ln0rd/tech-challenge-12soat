package order

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/order"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FindCompletedOrderById struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

type OrderWithInputs struct {
	Order       *domain.Order       `json:"order"`
	Vehicle     VehicleDetails      `json:"vehicle"`
	Inputs      []OrderInputDetails `json:"inputs"`
	TotalPrice  float64             `json:"total_price"`
	Timeline    map[string]string   `json:"timeline"`
	AverageTime string              `json:"average_time"`
}

type VehicleDetails struct {
	ID                          string `json:"id"`
	Model                       string `json:"model"`
	Brand                       string `json:"brand"`
	ReleaseYear                 int    `json:"release_year"`
	VehicleIdentificationNumber string `json:"vehicle_identification_number"`
	NumberPlate                 string `json:"number_plate"`
	Color                       string `json:"color"`
}

type OrderInputDetails struct {
	ID         string  `json:"id"`
	InputID    string  `json:"input_id"`
	InputName  string  `json:"input_name"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
}

// Formata duração em minutos para HH:MM:SS
func formatDuration(minutes int) string {
	if minutes <= 0 {
		return "00:00:00"
	}

	hours := minutes / 60
	remainingMinutes := minutes % 60
	seconds := 0 // Como não temos segundos, sempre será 0

	return fmt.Sprintf("%02d:%02d:%02d", hours, remainingMinutes, seconds)
}

// Converte segundos para formato HH:MM:SS
func formatDurationFromSeconds(seconds int) string {
	if seconds <= 0 {
		return "00:00:00"
	}

	hours := seconds / 3600
	remainingSeconds := seconds % 3600
	minutes := remainingSeconds / 60
	secs := remainingSeconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}

// Calcula timeline e tempo médio baseado no histórico de status
func (uc *FindCompletedOrderById) calculateTimeline(orderID uuid.UUID) (map[string]string, string) {
	var history []models.OrderStatusHistory

	if err := uc.DB.Where("order_id = ? ORDER BY started_at ASC", orderID).Find(&history).Error; err != nil {
		uc.Logger.Error("Error fetching order status history", zap.Error(err))
		return make(map[string]string), "00:00:00"
	}

	uc.Logger.Info("Found order status history",
		zap.String("orderID", orderID.String()),
		zap.Int("historyCount", len(history)))

	timeline := make(map[string]string)
	var totalSeconds int
	var completedStatuses int

	for _, status := range history {
		if status.EndedAt != nil {
			// Status finalizado - calcula duração baseada em started_at e ended_at
			duration := status.EndedAt.Sub(status.StartedAt)
			durationSeconds := int(duration.Seconds())

			timeline[status.Status] = formatDurationFromSeconds(durationSeconds)
			totalSeconds += durationSeconds
			completedStatuses++

			uc.Logger.Info("Status duration calculated",
				zap.String("status", status.Status),
				zap.Time("startedAt", status.StartedAt),
				zap.Time("endedAt", *status.EndedAt),
				zap.Int("durationSeconds", durationSeconds))
		} else {
			// Status atual (não finalizado)
			timeline[status.Status] = "00:00:00"
			uc.Logger.Info("Status not completed yet",
				zap.String("status", status.Status),
				zap.Time("startedAt", status.StartedAt))
		}
	}

	// Calcula tempo médio
	var averageTime string
	if completedStatuses > 0 {
		averageSeconds := totalSeconds / completedStatuses
		averageTime = formatDurationFromSeconds(averageSeconds)
	} else {
		averageTime = "00:00:00"
	}

	uc.Logger.Info("Timeline calculated",
		zap.String("orderID", orderID.String()),
		zap.Int("totalSeconds", totalSeconds),
		zap.Int("completedStatuses", completedStatuses),
		zap.String("averageTime", averageTime))

	return timeline, averageTime
}

func (uc *FindCompletedOrderById) Process(orderID uuid.UUID) (*OrderWithInputs, error) {
	uc.Logger.Info("Processing find completed order by ID", zap.String("orderID", orderID.String()))

	// Busca a order
	var order models.Order
	if err := uc.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		uc.Logger.Error("Order not found", zap.String("orderID", orderID.String()))
		return nil, errors.New("order not found")
	}

	uc.Logger.Info("Order found",
		zap.String("orderID", order.ID.String()),
		zap.String("customerID", order.CustomerID.String()),
		zap.String("vehicleID", order.VehicleID.String()),
		zap.String("status", order.Status))

	// Busca as informações do vehicle
	var vehicle models.Vehicle
	if err := uc.DB.Where("id = ?", order.VehicleID).First(&vehicle).Error; err != nil {
		uc.Logger.Error("Vehicle not found",
			zap.String("vehicleID", order.VehicleID.String()),
			zap.String("orderID", orderID.String()))
		return nil, errors.New("vehicle not found")
	}

	uc.Logger.Info("Vehicle found",
		zap.String("vehicleID", vehicle.ID.String()),
		zap.String("model", vehicle.Model),
		zap.String("brand", vehicle.Brand),
		zap.String("numberPlate", vehicle.NumberPlate))

	// Busca os inputs relacionados à order
	var orderInputs []models.OrderInput
	if err := uc.DB.Where("order_id = ?", orderID).Find(&orderInputs).Error; err != nil {
		uc.Logger.Error("Database error finding order inputs", zap.Error(err))
		return nil, err
	}

	uc.Logger.Info("Found order inputs", zap.Int("count", len(orderInputs)))

	// Busca os detalhes dos inputs e calcula o total
	var inputs []OrderInputDetails
	var totalPrice float64 = 0

	for _, orderInput := range orderInputs {
		// Busca o nome do input
		var input models.Input
		if err := uc.DB.Where("id = ?", orderInput.InputID).First(&input).Error; err != nil {
			uc.Logger.Error("Input not found for order input",
				zap.String("inputID", orderInput.InputID.String()),
				zap.String("orderInputID", orderInput.ID.String()))
			continue // Pula este input se não encontrar
		}

		inputDetail := OrderInputDetails{
			ID:         orderInput.ID.String(),
			InputID:    orderInput.InputID.String(),
			InputName:  input.Name,
			Quantity:   orderInput.Quantity,
			UnitPrice:  orderInput.UnitPrice,
			TotalPrice: orderInput.TotalPrice,
		}

		inputs = append(inputs, inputDetail)
		totalPrice += orderInput.TotalPrice

		uc.Logger.Info("Added input detail",
			zap.String("inputID", orderInput.InputID.String()),
			zap.String("inputName", input.Name),
			zap.Int("quantity", orderInput.Quantity),
			zap.Float64("unitPrice", orderInput.UnitPrice),
			zap.Float64("totalPrice", orderInput.TotalPrice))
	}

	uc.Logger.Info("Calculated total price", zap.Float64("totalPrice", totalPrice))

	// Mapeia para o domínio
	domainOrder := &domain.Order{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		VehicleID:  order.VehicleID,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}

	// Mapeia os detalhes do vehicle
	vehicleDetails := VehicleDetails{
		ID:                          vehicle.ID.String(),
		Model:                       vehicle.Model,
		Brand:                       vehicle.Brand,
		ReleaseYear:                 vehicle.ReleaseYear,
		VehicleIdentificationNumber: vehicle.VehicleIdentificationNumber,
		NumberPlate:                 vehicle.NumberPlate,
		Color:                       vehicle.Color,
	}

	// Calcula timeline e tempo médio
	timeline, averageTime := uc.calculateTimeline(orderID)

	result := &OrderWithInputs{
		Order:       domainOrder,
		Vehicle:     vehicleDetails,
		Inputs:      inputs,
		TotalPrice:  totalPrice,
		Timeline:    timeline,
		AverageTime: averageTime,
	}

	uc.Logger.Info("Completed order with inputs and timeline retrieved successfully",
		zap.String("orderID", orderID.String()),
		zap.Int("inputsCount", len(inputs)),
		zap.Int("timelineEntries", len(timeline)),
		zap.String("averageTime", averageTime))

	return result, nil
}
