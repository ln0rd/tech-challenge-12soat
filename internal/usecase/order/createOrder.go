package order

import (
	"errors"

	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/order"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CreateOrder struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (uc *CreateOrder) Process(entity *domain.Order) error {
	uc.Logger.Info("Processing order creation",
		zap.String("customerID", entity.CustomerID.String()),
		zap.String("vehicleID", entity.VehicleID.String()),
		zap.String("status", entity.Status))

	// Verifica se o customer existe
	var existingCustomer models.Customer
	if err := uc.DB.Where("id = ?", entity.CustomerID).First(&existingCustomer).Error; err != nil {
		uc.Logger.Error("Customer not found", zap.String("customerID", entity.CustomerID.String()))
		return errors.New("customer not found")
	}
	uc.Logger.Info("Customer found", zap.String("customerID", entity.CustomerID.String()))

	// Verifica se o vehicle existe
	var existingVehicle models.Vehicle
	if err := uc.DB.Where("id = ?", entity.VehicleID).First(&existingVehicle).Error; err != nil {
		uc.Logger.Error("Vehicle not found", zap.String("vehicleID", entity.VehicleID.String()))
		return errors.New("vehicle not found")
	}
	uc.Logger.Info("Vehicle found", zap.String("vehicleID", entity.VehicleID.String()))

	// Verifica se o vehicle pertence ao customer
	if existingVehicle.CustomerID != entity.CustomerID {
		uc.Logger.Error("Vehicle does not belong to customer",
			zap.String("vehicleID", entity.VehicleID.String()),
			zap.String("customerID", entity.CustomerID.String()))
		return errors.New("vehicle does not belong to customer")
	}
	uc.Logger.Info("Vehicle belongs to customer")

	model := &models.Order{
		ID:         entity.ID,
		CustomerID: entity.CustomerID,
		VehicleID:  entity.VehicleID,
		Status:     entity.Status,
	}

	uc.Logger.Info("Model created",
		zap.String("id", model.ID.String()),
		zap.String("customerID", model.CustomerID.String()),
		zap.String("vehicleID", model.VehicleID.String()),
		zap.String("status", model.Status))

	result := uc.DB.Create(model)
	if result.Error != nil {
		uc.Logger.Error("Database error creating order", zap.Error(result.Error))
		return result.Error
	}

	uc.Logger.Info("Order created in database",
		zap.String("id", model.ID.String()),
		zap.Int64("rowsAffected", result.RowsAffected))

	return nil
}
