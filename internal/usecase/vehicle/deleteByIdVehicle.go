package vehicle

import (
	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DeleteByIdVehicle struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (uc *DeleteByIdVehicle) Process(id uuid.UUID) error {
	uc.Logger.Info("Processing delete vehicle by ID", zap.String("id", id.String()))

	result := uc.DB.Where("id = ?", id).Delete(&models.Vehicle{})
	if result.Error != nil {
		uc.Logger.Error("Database error deleting vehicle", zap.Error(result.Error), zap.String("id", id.String()))
		return result.Error
	}

	if result.RowsAffected == 0 {
		uc.Logger.Warn("No vehicle found to delete", zap.String("id", id.String()))
		return gorm.ErrRecordNotFound
	}

	uc.Logger.Info("Vehicle deleted successfully", zap.String("id", id.String()), zap.Int64("rowsAffected", result.RowsAffected))

	return nil
}
