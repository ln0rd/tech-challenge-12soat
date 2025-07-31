package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/auth"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/auth"
	"go.uber.org/zap"
)

type AuthController struct {
	Logger       *zap.Logger
	LoginUseCase *auth.LoginUseCase
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (dto *LoginDTO) Validate() error {
	if dto.Email == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(dto.Email) {
		return errors.New("invalid email format")
	}
	if dto.Password == "" {
		return errors.New("password is required")
	}
	if len(dto.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var dto LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		ac.Logger.Error("Error decoding JSON", zap.Error(err))
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	ac.Logger.Info("Received login request", zap.String("email", dto.Email))

	if err := dto.Validate(); err != nil {
		ac.Logger.Error("Validation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ac.Logger.Info("Validation passed")

	request := domain.LoginRequest{
		Email:    dto.Email,
		Password: dto.Password,
	}

	ac.Logger.Info("Calling LoginUseCase.Execute...")
	response, err := ac.LoginUseCase.Execute(request)
	if err != nil {
		ac.Logger.Error("Login failed", zap.Error(err))
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	ac.Logger.Info("Login successful", zap.String("email", response.User.Email))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
