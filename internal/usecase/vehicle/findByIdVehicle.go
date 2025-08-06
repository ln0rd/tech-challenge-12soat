package vehicle

import (
	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/vehicle"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/logger"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/repository"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"go.uber.org/zap"
)

type FindByIdVehicle struct {
	VehicleRepository repository.VehicleRepository
	Logger            logger.Logger
}

// FetchVehicleFromDB busca um vehicle específico do banco de dados
func (uc *FindByIdVehicle) FetchVehicleFromDB(id uuid.UUID) (*models.Vehicle, error) {
	vehicle, err := uc.VehicleRepository.FindByID(id)
	if err != nil {
		uc.Logger.Error("Database error finding vehicle by ID", zap.Error(err), zap.String("id", id.String()))
		return nil, err
	}

	uc.Logger.Info("Found vehicle in database",
		zap.String("id", vehicle.ID.String()),
		zap.String("model", vehicle.Model),
		zap.String("brand", vehicle.Brand),
		zap.String("numberPlate", vehicle.NumberPlate))

	return vehicle, nil
}

func (uc *FindByIdVehicle) Process(id uuid.UUID) (*domain.Vehicle, error) {
	uc.Logger.Info("Processing find vehicle by ID", zap.String("id", id.String()))

	// Busca vehicle do banco
	vehicle, err := uc.FetchVehicleFromDB(id)
	if err != nil {
		return nil, err
	}

	// Mapeia para o domínio usando persistence
	domainVehicle := persistence.VehiclePersistence{}.ToEntity(vehicle)
	uc.Logger.Info("Successfully mapped vehicle to domain", zap.String("id", domainVehicle.ID.String()))

	return domainVehicle, nil
}
