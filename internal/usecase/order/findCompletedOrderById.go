package order

import (
	"errors"

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
	Order      *domain.Order       `json:"order"`
	Vehicle    VehicleDetails      `json:"vehicle"`
	Inputs     []OrderInputDetails `json:"inputs"`
	TotalPrice float64             `json:"total_price"`
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

	result := &OrderWithInputs{
		Order:      domainOrder,
		Vehicle:    vehicleDetails,
		Inputs:     inputs,
		TotalPrice: totalPrice,
	}

	uc.Logger.Info("Completed order with inputs retrieved successfully",
		zap.String("orderID", orderID.String()),
		zap.Int("inputsCount", len(inputs)))

	return result, nil
}
