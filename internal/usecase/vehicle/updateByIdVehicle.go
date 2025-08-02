package vehicle

import (
	"errors"

	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/vehicle"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UpdateByIdVehicle struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (uc *UpdateByIdVehicle) Process(id uuid.UUID, entity *domain.Vehicle) error {
	uc.Logger.Info("Processing update vehicle by ID",
		zap.String("id", id.String()),
		zap.String("model", entity.Model),
		zap.String("brand", entity.Brand),
		zap.String("numberPlate", entity.NumberPlate))

	// Primeiro verifica se o vehicle existe
	var existingVehicle models.Vehicle
	if err := uc.DB.Where("id = ?", id).First(&existingVehicle).Error; err != nil {
		uc.Logger.Error("Database error finding vehicle to update", zap.Error(err), zap.String("id", id.String()))
		return err
	}

	uc.Logger.Info("Found existing vehicle",
		zap.String("id", existingVehicle.ID.String()),
		zap.String("model", existingVehicle.Model),
		zap.String("brand", existingVehicle.Brand),
		zap.String("numberPlate", existingVehicle.NumberPlate))

	// Verifica se a nova placa já existe (se foi alterada)
	if entity.NumberPlate != existingVehicle.NumberPlate {
		var vehicleWithSamePlate models.Vehicle
		if err := uc.DB.Where("number_plate = ? AND id != ?", entity.NumberPlate, id).First(&vehicleWithSamePlate).Error; err == nil {
			uc.Logger.Error("Number plate already exists", zap.String("numberPlate", entity.NumberPlate))
			return errors.New("number plate already exists")
		} else if err != gorm.ErrRecordNotFound {
			uc.Logger.Error("Error checking number plate uniqueness", zap.Error(err))
			return err
		}
		uc.Logger.Info("Number plate is unique", zap.String("numberPlate", entity.NumberPlate))
	}

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

	// Atualiza os campos do vehicle existente
	existingVehicle.Model = entity.Model
	existingVehicle.Brand = entity.Brand
	existingVehicle.ReleaseYear = entity.ReleaseYear
	existingVehicle.VehicleIdentificationNumber = entity.VehicleIdentificationNumber
	existingVehicle.NumberPlate = entity.NumberPlate
	existingVehicle.Color = entity.Color
	existingVehicle.CustomerID = entity.CustomerID

	uc.Logger.Info("Updated vehicle fields",
		zap.String("model", existingVehicle.Model),
		zap.String("brand", existingVehicle.Brand),
		zap.String("numberPlate", existingVehicle.NumberPlate),
		zap.String("color", existingVehicle.Color))

	// Salva as alterações
	result := uc.DB.Save(&existingVehicle)
	if result.Error != nil {
		uc.Logger.Error("Database error updating vehicle", zap.Error(result.Error))
		return result.Error
	}

	uc.Logger.Info("Vehicle updated successfully",
		zap.String("id", existingVehicle.ID.String()),
		zap.Int64("rowsAffected", result.RowsAffected))

	return nil
}
