package vehicle

import (
	"errors"

	"github.com/google/uuid"
	interfaces "github.com/ln0rd/tech_challenge_12soat/internal/domain/interfaces"
	vehicleDomain "github.com/ln0rd/tech_challenge_12soat/internal/domain/vehicle"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/repository"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CreateVehicle struct {
	VehicleRepository  repository.VehicleRepository
	CustomerRepository repository.CustomerRepository
	Logger             interfaces.Logger
}

// ValidateNumberPlateUniqueness verifica se a placa do vehicle é única
func (uc *CreateVehicle) ValidateNumberPlateUniqueness(numberPlate string) error {
	_, err := uc.VehicleRepository.FindByNumberPlate(numberPlate)
	if err == nil {
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
func (uc *CreateVehicle) ValidateCustomerExists(customerID uuid.UUID) error {
	if customerID == uuid.Nil {
		uc.Logger.Error("Customer ID is required")
		return errors.New("customer ID is required")
	}

	_, err := uc.CustomerRepository.FindByID(customerID)
	if err != nil {
		uc.Logger.Error("Customer not found", zap.String("customerID", customerID.String()))
		return errors.New("customer not found")
	}

	uc.Logger.Info("Customer found", zap.String("customerID", customerID.String()))
	return nil
}

// SaveVehicleToDB salva o vehicle no banco de dados
func (uc *CreateVehicle) SaveVehicleToDB(model *models.Vehicle) error {
	err := uc.VehicleRepository.Create(model)
	if err != nil {
		uc.Logger.Error("Database error creating vehicle", zap.Error(err))
		return err
	}

	uc.Logger.Info("Vehicle created in database",
		zap.String("id", model.ID.String()),
		zap.String("numberPlate", model.NumberPlate))
	return nil
}

func (uc *CreateVehicle) Process(entity *vehicleDomain.Vehicle) error {
	uc.Logger.Info("Processing vehicle creation",
		zap.String("model", entity.Model),
		zap.String("brand", entity.Brand),
		zap.String("numberPlate", entity.NumberPlate))

	// Valida unicidade da placa
	if err := uc.ValidateNumberPlateUniqueness(entity.NumberPlate); err != nil {
		return err
	}

	// Valida existência do customer
	if err := uc.ValidateCustomerExists(entity.CustomerID); err != nil {
		return err
	}

	// Mapeia entidade para modelo usando persistence
	model := persistence.VehiclePersistence{}.ToModel(entity)
	uc.Logger.Info("Model created",
		zap.String("model", model.Model),
		zap.String("brand", model.Brand),
		zap.String("numberPlate", model.NumberPlate))

	// Salva no banco
	err := uc.SaveVehicleToDB(model)
	if err != nil {
		return err
	}

	return nil
}
