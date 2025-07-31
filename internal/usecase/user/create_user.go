package user

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/user"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CreateUser struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (uc *CreateUser) Process(entity *domain.User) error {
	uc.Logger.Info("Processing user creation", zap.String("email", entity.Email))

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(entity.Password), bcrypt.DefaultCost)
	if err != nil {
		uc.Logger.Error("Error hashing password", zap.Error(err))
		return err
	}

	// Atualiza a senha com o hash
	entity.Password = string(hashedPassword)

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
