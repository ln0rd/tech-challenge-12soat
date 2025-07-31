package user

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/user"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CreateUser struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (uc *CreateUser) Process(entity *domain.User) error {
	uc.Logger.Info("Processing user creation", zap.String("email", entity.Email))

	model := persistence.UserPersistence{}.ToModel(entity)

	uc.Logger.Info("Model created", zap.String("email", model.Email), zap.String("username", model.Username))

	result := uc.DB.Create(model)
	if result.Error != nil {
		uc.Logger.Error("Database error creating user", zap.Error(result.Error))
		return result.Error
	}

	uc.Logger.Info("User created in database", zap.String("email", model.Email), zap.Int64("rowsAffected", result.RowsAffected))

	return nil
}
