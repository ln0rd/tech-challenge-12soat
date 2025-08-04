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

// FetchVehicleFromDB busca um vehicle específico do banco de dados
func (uc *UpdateByIdVehicle) FetchVehicleFromDB(id uuid.UUID) (*models.Vehicle, error) {
	var existingVehicle models.Vehicle
	if err := uc.DB.Where("id = ?", id).First(&existingVehicle).Error; err != nil {
		uc.Logger.Error("Database error finding vehicle to update", zap.Error(err), zap.String("id", id.String()))
		return nil, err
	}

	uc.Logger.Info("Found existing vehicle",
		zap.String("id", existingVehicle.ID.String()),
		zap.String("model", existingVehicle.Model),
		zap.String("brand", existingVehicle.Brand),
		zap.String("numberPlate", existingVehicle.NumberPlate))

	return &existingVehicle, nil
}

// ValidateNumberPlateUniqueness verifica se a placa do vehicle é única (para update)
func (uc *UpdateByIdVehicle) ValidateNumberPlateUniqueness(numberPlate string, vehicleID uuid.UUID) error {
	var vehicleWithSamePlate models.Vehicle
	if err := uc.DB.Where("number_plate = ? AND id != ?", numberPlate, vehicleID).First(&vehicleWithSamePlate).Error; err == nil {
		uc.Logger.Error("Number plate already exists", zap.String("numberPlate", numberPlate))
		return errors.New("number plate already exists")
	} else if err != gorm.ErrRecordNotFound {
		uc.Logger.Error("Error checking number plate uniqueness", zap.Error(err))
		return err
	}

	uc.Logger.Info("Number plate is unique", zap.String("numberPlate", numberPlate))
	return nil
}

// ValidateCustomerExists verifica se o customer existe no banco
func (uc *UpdateByIdVehicle) ValidateCustomerExists(customerID uuid.UUID) error {
	if customerID == uuid.Nil {
		uc.Logger.Error("Customer ID is required")
		return errors.New("customer ID is required")
	}

	var existingCustomer models.Customer
	if err := uc.DB.Where("id = ?", customerID).First(&existingCustomer).Error; err != nil {
		uc.Logger.Error("Customer not found", zap.String("customerID", customerID.String()))
		return errors.New("customer not found")
	}

	uc.Logger.Info("Customer found", zap.String("customerID", customerID.String()))
	return nil
}

// UpdateVehicleFields atualiza os campos do vehicle existente
func (uc *UpdateByIdVehicle) UpdateVehicleFields(existingVehicle *models.Vehicle, entity *domain.Vehicle) {
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
}

// SaveVehicleToDB salva as alterações do vehicle no banco de dados
func (uc *UpdateByIdVehicle) SaveVehicleToDB(vehicle *models.Vehicle) error {
	result := uc.DB.Save(vehicle)
	if result.Error != nil {
		uc.Logger.Error("Database error updating vehicle", zap.Error(result.Error))
		return result.Error
	}

	uc.Logger.Info("Vehicle updated successfully",
		zap.String("id", vehicle.ID.String()),
		zap.Int64("rowsAffected", result.RowsAffected))
	return nil
}

func (uc *UpdateByIdVehicle) Process(id uuid.UUID, entity *domain.Vehicle) error {
	uc.Logger.Info("Processing update vehicle by ID",
		zap.String("id", id.String()),
		zap.String("model", entity.Model),
		zap.String("brand", entity.Brand),
		zap.String("numberPlate", entity.NumberPlate))

	// Busca o vehicle existente
	existingVehicle, err := uc.FetchVehicleFromDB(id)
	if err != nil {
		return err
	}

	// Verifica se a nova placa já existe (se foi alterada)
	if entity.NumberPlate != existingVehicle.NumberPlate {
		if err := uc.ValidateNumberPlateUniqueness(entity.NumberPlate, id); err != nil {
			return err
		}
	}

	// Valida existência do customer
	if err := uc.ValidateCustomerExists(entity.CustomerID); err != nil {
		return err
	}

	// Atualiza os campos do vehicle
	uc.UpdateVehicleFields(existingVehicle, entity)

	// Salva as alterações
	err = uc.SaveVehicleToDB(existingVehicle)
	if err != nil {
		return err
	}

	return nil
}
