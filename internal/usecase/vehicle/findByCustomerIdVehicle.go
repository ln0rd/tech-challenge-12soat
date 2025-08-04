package vehicle

import (
	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/vehicle"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FindByCustomerIdVehicle struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

// FetchVehiclesFromDB busca vehicles por customer ID do banco de dados
func (uc *FindByCustomerIdVehicle) FetchVehiclesFromDB(customerID uuid.UUID) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	if err := uc.DB.Where("customer_id = ?", customerID).Find(&vehicles).Error; err != nil {
		uc.Logger.Error("Database error finding vehicles by customer ID", zap.Error(err), zap.String("customerID", customerID.String()))
		return []models.Vehicle{}, err
	}

	uc.Logger.Info("Found vehicles in database", zap.Int("count", len(vehicles)), zap.String("customerID", customerID.String()))
	return vehicles, nil
}

func (uc *FindByCustomerIdVehicle) Process(customerID uuid.UUID) ([]domain.Vehicle, error) {
	uc.Logger.Info("Processing find vehicles by customer ID", zap.String("customerID", customerID.String()))

	// Busca vehicles do banco
	vehicles, err := uc.FetchVehiclesFromDB(customerID)
	if err != nil {
		return []domain.Vehicle{}, err
	}

	// Mapeia para o domínio usando persistence
	domainVehicles := make([]domain.Vehicle, 0)
	for _, vehicle := range vehicles {
		domainVehicle := persistence.VehiclePersistence{}.ToEntity(&vehicle)
		domainVehicles = append(domainVehicles, *domainVehicle)
	}

	uc.Logger.Info("Successfully mapped vehicles to domain", zap.Int("count", len(domainVehicles)))

	// Sempre retorna uma lista (vazia se não encontrou veículos)
	return domainVehicles, nil
}
