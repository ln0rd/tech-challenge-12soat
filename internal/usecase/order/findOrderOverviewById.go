package order

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/order"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/logger"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/repository"
	"go.uber.org/zap"
)

type FindOrderOverviewById struct {
	OrderRepository              repository.OrderRepository
	VehicleRepository            repository.VehicleRepository
	OrderInputRepository         repository.OrderInputRepository
	OrderStatusHistoryRepository repository.OrderStatusHistoryRepository
	InputRepository              repository.InputRepository
	Logger                       logger.Logger
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

// FormatDurationFromSeconds converte segundos para formato HH:MM:SS
func FormatDurationFromSeconds(seconds int) string {
	if seconds <= 0 {
		return "00:00:00"
	}

	hours := seconds / 3600
	remainingSeconds := seconds % 3600
	minutes := remainingSeconds / 60
	secs := remainingSeconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}

// FetchOrderFromDB busca um order específico do banco de dados
func (uc *FindOrderOverviewById) FetchOrderFromDB(orderID uuid.UUID) (*models.Order, error) {
	order, err := uc.OrderRepository.FindByID(orderID)
	if err != nil {
		uc.Logger.Error("Order not found", zap.String("orderID", orderID.String()))
		return nil, errors.New("order not found")
	}

	uc.Logger.Info("Order found",
		zap.String("orderID", order.ID.String()),
		zap.String("customerID", order.CustomerID.String()),
		zap.String("vehicleID", order.VehicleID.String()),
		zap.String("status", order.Status))

	return order, nil
}

// FetchVehicleFromDB busca um vehicle específico do banco de dados
func (uc *FindOrderOverviewById) FetchVehicleFromDB(vehicleID uuid.UUID) (*models.Vehicle, error) {
	vehicle, err := uc.VehicleRepository.FindByID(vehicleID)
	if err != nil {
		uc.Logger.Error("Vehicle not found",
			zap.String("vehicleID", vehicleID.String()))
		return nil, errors.New("vehicle not found")
	}

	uc.Logger.Info("Vehicle found",
		zap.String("vehicleID", vehicle.ID.String()),
		zap.String("model", vehicle.Model),
		zap.String("brand", vehicle.Brand),
		zap.String("numberPlate", vehicle.NumberPlate))

	return vehicle, nil
}

// FetchOrderInputsFromDB busca os inputs da order do banco de dados
func (uc *FindOrderOverviewById) FetchOrderInputsFromDB(orderID uuid.UUID) ([]models.OrderInput, error) {
	orderInputs, err := uc.OrderInputRepository.FindByOrderID(orderID)
	if err != nil {
		uc.Logger.Error("Database error finding order inputs", zap.Error(err))
		return nil, err
	}

	uc.Logger.Info("Order inputs found", zap.Int("count", len(orderInputs)))
	return orderInputs, nil
}

// MapVehicleToDetails mapeia o vehicle para VehicleDetails
func (uc *FindOrderOverviewById) MapVehicleToDetails(vehicle *models.Vehicle) VehicleDetails {
	return VehicleDetails{
		ID:                          vehicle.ID.String(),
		Model:                       vehicle.Model,
		Brand:                       vehicle.Brand,
		ReleaseYear:                 vehicle.ReleaseYear,
		VehicleIdentificationNumber: vehicle.VehicleIdentificationNumber,
		NumberPlate:                 vehicle.NumberPlate,
		Color:                       vehicle.Color,
	}
}

// MapOrderInputToDetails mapeia um OrderInput para OrderInputDetails
func (uc *FindOrderOverviewById) MapOrderInputToDetails(orderInput models.OrderInput, inputName string) OrderInputDetails {
	return OrderInputDetails{
		ID:         orderInput.ID.String(),
		InputID:    orderInput.InputID.String(),
		InputName:  inputName,
		Quantity:   orderInput.Quantity,
		UnitPrice:  orderInput.UnitPrice,
		TotalPrice: orderInput.TotalPrice,
	}
}

// ProcessOrderInputs processa os inputs da order e calcula o total
func (uc *FindOrderOverviewById) ProcessOrderInputs(orderInputs []models.OrderInput) ([]OrderInputDetails, float64) {
	var inputs []OrderInputDetails
	var totalPrice float64 = 0

	for _, orderInput := range orderInputs {
		// Busca o nome do input
		input, err := uc.InputRepository.FindByID(orderInput.InputID)
		if err != nil {
			uc.Logger.Error("Input not found for order input",
				zap.String("inputID", orderInput.InputID.String()),
				zap.String("orderInputID", orderInput.ID.String()))
			continue // Pula este input se não encontrar
		}

		inputDetail := uc.MapOrderInputToDetails(orderInput, input.Name)
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
	return inputs, totalPrice
}

// MapOrderToDomain mapeia a order para o domínio
func (uc *FindOrderOverviewById) MapOrderToDomain(order *models.Order) *domain.Order {
	return &domain.Order{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		VehicleID:  order.VehicleID,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}

// CalculateTimeline calcula timeline e tempo médio baseado no histórico de status
func (uc *FindOrderOverviewById) CalculateTimeline(orderID uuid.UUID) (map[string]string, string) {
	history, err := uc.OrderStatusHistoryRepository.FindByOrderID(orderID)
	if err != nil {
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

			timeline[status.Status] = FormatDurationFromSeconds(durationSeconds)
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
		averageTime = FormatDurationFromSeconds(averageSeconds)
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

func (uc *FindOrderOverviewById) Process(orderID uuid.UUID) (*OrderWithInputs, error) {
	uc.Logger.Info("Processing find completed order by ID", zap.String("orderID", orderID.String()))

	// Busca a order
	order, err := uc.FetchOrderFromDB(orderID)
	if err != nil {
		return nil, err
	}

	// Busca as informações do vehicle
	vehicle, err := uc.FetchVehicleFromDB(order.VehicleID)
	if err != nil {
		return nil, err
	}

	// Busca os inputs relacionados à order
	orderInputs, err := uc.FetchOrderInputsFromDB(orderID)
	if err != nil {
		return nil, err
	}

	// Processa os inputs e calcula o total
	inputs, totalPrice := uc.ProcessOrderInputs(orderInputs)

	// Mapeia para o domínio
	domainOrder := uc.MapOrderToDomain(order)

	// Mapeia os detalhes do vehicle
	vehicleDetails := uc.MapVehicleToDetails(vehicle)

	// Calcula timeline e tempo médio
	timeline, averageTime := uc.CalculateTimeline(orderID)

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
