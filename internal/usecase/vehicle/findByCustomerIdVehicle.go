package vehicle

import (
	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/vehicle"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FindByCustomerIdVehicle struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (uc *FindByCustomerIdVehicle) Process(customerID uuid.UUID) ([]domain.Vehicle, error) {
	uc.Logger.Info("Processing find vehicles by customer ID", zap.String("customerID", customerID.String()))

	var vehicles []models.Vehicle
	if err := uc.DB.Where("customer_id = ?", customerID).Find(&vehicles).Error; err != nil {
		uc.Logger.Error("Database error finding vehicles by customer ID", zap.Error(err), zap.String("customerID", customerID.String()))
		return []domain.Vehicle{}, err
	}

	uc.Logger.Info("Found vehicles in database", zap.Int("count", len(vehicles)), zap.String("customerID", customerID.String()))

	// Inicializa uma lista vazia
	domainVehicles := make([]domain.Vehicle, 0)

	for _, vehicle := range vehicles {
		domainVehicles = append(domainVehicles, domain.Vehicle{
			ID:                          vehicle.ID,
			Model:                       vehicle.Model,
			Brand:                       vehicle.Brand,
			ReleaseYear:                 vehicle.ReleaseYear,
			VehicleIdentificationNumber: vehicle.VehicleIdentificationNumber,
			NumberPlate:                 vehicle.NumberPlate,
			Color:                       vehicle.Color,
			CustomerID:                  vehicle.CustomerID,
			CreatedAt:                   vehicle.CreatedAt,
			UpdatedAt:                   vehicle.UpdatedAt,
		})
	}

	uc.Logger.Info("Successfully mapped vehicles to domain", zap.Int("count", len(domainVehicles)))

	// Sempre retorna uma lista (vazia se não encontrou veículos)
	return domainVehicles, nil
}
