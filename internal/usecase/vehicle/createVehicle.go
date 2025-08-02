package vehicle

import (
	"errors"

	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/vehicle"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CreateVehicle struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (uc *CreateVehicle) Process(entity *domain.Vehicle) error {
	uc.Logger.Info("Processing vehicle creation",
		zap.String("model", entity.Model),
		zap.String("brand", entity.Brand),
		zap.String("numberPlate", entity.NumberPlate))

	// Verifica se a placa já existe
	var existingVehicle models.Vehicle
	if err := uc.DB.Where("number_plate = ?", entity.NumberPlate).First(&existingVehicle).Error; err == nil {
		uc.Logger.Error("Number plate already exists", zap.String("numberPlate", entity.NumberPlate))
		return errors.New("number plate already exists")
	} else if err != gorm.ErrRecordNotFound {
		uc.Logger.Error("Error checking number plate uniqueness", zap.Error(err))
		return err
	}

	uc.Logger.Info("Number plate is unique", zap.String("numberPlate", entity.NumberPlate))

	// Verifica se o CustomerID existe
	if entity.CustomerID == uuid.Nil {
		uc.Logger.Error("Customer ID is required")
		return errors.New("customer ID is required")
	}

	var existingCustomer models.Customer
	if err := uc.DB.Where("id = ?", entity.CustomerID).First(&existingCustomer).Error; err != nil {
		uc.Logger.Error("Customer not found", zap.String("customerID", entity.CustomerID.String()))
		return errors.New("customer not found")
	}
	uc.Logger.Info("Customer found", zap.String("customerID", entity.CustomerID.String()))

	model := &models.Vehicle{
		ID:                          entity.ID,
		Model:                       entity.Model,
		Brand:                       entity.Brand,
		ReleaseYear:                 entity.ReleaseYear,
		VehicleIdentificationNumber: entity.VehicleIdentificationNumber,
		NumberPlate:                 entity.NumberPlate,
		Color:                       entity.Color,
		CustomerID:                  entity.CustomerID,
	}

	uc.Logger.Info("Model created",
		zap.String("model", model.Model),
		zap.String("brand", model.Brand),
		zap.String("numberPlate", model.NumberPlate))

	result := uc.DB.Create(model)
	if result.Error != nil {
		uc.Logger.Error("Database error creating vehicle", zap.Error(result.Error))
		return result.Error
	}

	uc.Logger.Info("Vehicle created in database",
		zap.String("id", model.ID.String()),
		zap.String("numberPlate", model.NumberPlate),
		zap.Int64("rowsAffected", result.RowsAffected))

	return nil
}
