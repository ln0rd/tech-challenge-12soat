package vehicle

import (
	"github.com/google/uuid"
	interfaces "github.com/ln0rd/tech_challenge_12soat/internal/domain/interfaces"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/repository"
	"go.uber.org/zap"
)

type DeleteByIdVehicle struct {
	VehicleRepository repository.VehicleRepository
	Logger            interfaces.Logger
}

// DeleteVehicleFromDB remove o vehicle do banco de dados
func (uc *DeleteByIdVehicle) DeleteVehicleFromDB(id uuid.UUID) error {
	err := uc.VehicleRepository.Delete(id)
	if err != nil {
		uc.Logger.Error("Database error deleting vehicle", zap.Error(err), zap.String("id", id.String()))
		return err
	}

	uc.Logger.Info("Vehicle deleted successfully", zap.String("id", id.String()))
	return nil
}

func (uc *DeleteByIdVehicle) Process(id uuid.UUID) error {
	uc.Logger.Info("Processing delete vehicle by ID", zap.String("id", id.String()))

	// Remove o vehicle do banco
	err := uc.DeleteVehicleFromDB(id)
	if err != nil {
		return err
	}

	return nil
}
