package user

import (
	"errors"

	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/user"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/logger"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/repository"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CreateUser struct {
	UserRepository repository.UserRepository
	Logger         logger.Logger
}

// ValidateEmailUniqueness verifica se o email é único
func (uc *CreateUser) ValidateEmailUniqueness(email string) error {
	_, err := uc.UserRepository.FindByEmail(email)
	if err == nil {
		uc.Logger.Error("Email already exists", zap.String("email", email))
		return errors.New("email already exists")
	} else if err != gorm.ErrRecordNotFound {
		uc.Logger.Error("Error checking email uniqueness", zap.Error(err))
		return err
	}

	uc.Logger.Info("Email is unique", zap.String("email", email))
	return nil
}

// HashPassword faz o hash da senha
func (uc *CreateUser) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		uc.Logger.Error("Error hashing password", zap.Error(err))
		return "", err
	}

	uc.Logger.Info("Password hashed successfully")
	return string(hashedPassword), nil
}

// SaveUserToDB salva o user no banco de dados
func (uc *CreateUser) SaveUserToDB(model *models.User) error {
	err := uc.UserRepository.Create(model)
	if err != nil {
		uc.Logger.Error("Database error creating user", zap.Error(err))
		return err
	}

	uc.Logger.Info("User created in database",
		zap.String("id", model.ID.String()),
		zap.String("email", model.Email))
	return nil
}

func (uc *CreateUser) Process(entity *domain.User) error {
	uc.Logger.Info("Processing user creation", zap.String("email", entity.Email))

	// Valida unicidade do email
	if err := uc.ValidateEmailUniqueness(entity.Email); err != nil {
		return err
	}

	// Hash da senha
	hashedPassword, err := uc.HashPassword(entity.Password)
	if err != nil {
		return err
	}

	// Atualiza a senha com o hash
	entity.Password = hashedPassword

	// Mapeia entidade para modelo usando persistence
	model := persistence.UserPersistence{}.ToModel(entity)
	uc.Logger.Info("Model created",
		zap.String("email", model.Email),
		zap.String("username", model.Username))

	// Salva no banco
	err = uc.SaveUserToDB(model)
	if err != nil {
		return err
	}

	return nil
}
