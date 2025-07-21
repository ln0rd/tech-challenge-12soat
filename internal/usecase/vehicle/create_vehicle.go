package vehicle

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/vehicle"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"gorm.io/gorm"
)

type CreateVehicle struct {
	DB *gorm.DB
}

func (uc *CreateVehicle) Process(entity *domain.Vehicle) error {
	model := persistence.VehiclePersistence{}.ToModel(entity)
	return uc.DB.Create(model).Error
}
