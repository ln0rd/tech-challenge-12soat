package auth

import (
	"errors"

	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/auth"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewAuthRepository(db *gorm.DB, logger *zap.Logger) *AuthRepository {
	return &AuthRepository{
		db:     db,
		logger: logger,
	}
}

func (r *AuthRepository) FindUserByEmail(email string) (*domain.UserInfo, error) {
	r.logger.Info("Finding user by email", zap.String("email", email))

	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		r.logger.Error("User not found", zap.Error(err), zap.String("email", email))
		return nil, err
	}

	r.logger.Info("User found", zap.String("email", user.Email), zap.String("username", user.Username))

	return &domain.UserInfo{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		UserType: user.UserType,
	}, nil
}

func (r *AuthRepository) ValidatePassword(email, password string) error {
	r.logger.Info("Validating password", zap.String("email", email))

	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		r.logger.Error("User not found for password validation", zap.Error(err), zap.String("email", email))
		return err
	}

	// Compara a senha fornecida com o hash armazenado
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		r.logger.Error("Invalid password", zap.Error(err), zap.String("email", email))
		return errors.New("invalid password")
	}

	r.logger.Info("Password validated successfully", zap.String("email", email))
	return nil
}
