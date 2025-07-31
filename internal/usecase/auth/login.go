package auth

import (
	"errors"
	"time"

	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/auth"
	"go.uber.org/zap"
)

type LoginUseCase struct {
	authRepository domain.AuthRepository
	tokenService   domain.TokenService
	logger         *zap.Logger
}

func NewLoginUseCase(authRepository domain.AuthRepository, tokenService domain.TokenService, logger *zap.Logger) *LoginUseCase {
	return &LoginUseCase{
		authRepository: authRepository,
		tokenService:   tokenService,
		logger:         logger,
	}
}

func (uc *LoginUseCase) Execute(request domain.LoginRequest) (*domain.LoginResponse, error) {
	uc.logger.Info("Processing login request", zap.String("email", request.Email))

	// Busca o usuário por email
	userInfo, err := uc.authRepository.FindUserByEmail(request.Email)
	if err != nil {
		uc.logger.Error("User not found", zap.Error(err), zap.String("email", request.Email))
		return nil, errors.New("invalid credentials")
	}

	uc.logger.Info("User found", zap.String("email", userInfo.Email), zap.String("username", userInfo.Username))

	// Valida a senha
	err = uc.authRepository.ValidatePassword(request.Email, request.Password)
	if err != nil {
		uc.logger.Error("Invalid password", zap.Error(err), zap.String("email", request.Email))
		return nil, errors.New("invalid credentials")
	}

	uc.logger.Info("Password validated successfully")

	// Gera o token JWT
	token, err := uc.tokenService.GenerateToken(*userInfo)
	if err != nil {
		uc.logger.Error("Error generating token", zap.Error(err))
		return nil, errors.New("error generating token")
	}

	// Gera o refresh token
	refreshToken, err := uc.tokenService.GenerateRefreshToken(userInfo.ID)
	if err != nil {
		uc.logger.Error("Error generating refresh token", zap.Error(err))
		return nil, errors.New("error generating refresh token")
	}

	// Calcula a data de expiração (24 horas)
	expiresAt := time.Now().Add(24 * time.Hour)

	uc.logger.Info("Login successful", zap.String("email", userInfo.Email))

	return &domain.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User:         *userInfo,
	}, nil
}
