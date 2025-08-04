package input

import (
	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DeleteByIdInput struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

// DeleteInputFromDB remove o input do banco de dados
func (uc *DeleteByIdInput) DeleteInputFromDB(id uuid.UUID) error {
	result := uc.DB.Where("id = ?", id).Delete(&models.Input{})
	if result.Error != nil {
		uc.Logger.Error("Database error deleting input", zap.Error(result.Error), zap.String("id", id.String()))
		return result.Error
	}

	if result.RowsAffected == 0 {
		uc.Logger.Warn("No input found to delete", zap.String("id", id.String()))
		return gorm.ErrRecordNotFound
	}

	uc.Logger.Info("Input deleted successfully", zap.String("id", id.String()), zap.Int64("rowsAffected", result.RowsAffected))
	return nil
}

func (uc *DeleteByIdInput) Process(id uuid.UUID) error {
	uc.Logger.Info("Processing delete input by ID", zap.String("id", id.String()))

	// Remove o input do banco
	err := uc.DeleteInputFromDB(id)
	if err != nil {
		return err
	}

	return nil
}
