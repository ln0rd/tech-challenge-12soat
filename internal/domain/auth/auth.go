package auth

import (
	"time"

	"github.com/google/uuid"
)

// LoginRequest representa a requisição de login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse representa a resposta de login
type LoginResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserInfo  `json:"user"`
}

// UserInfo representa as informações do usuário no token
type UserInfo struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	UserType string    `json:"user_type"`
}

// Claims representa as claims do JWT
type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	UserType string    `json:"user_type"`
	Exp      int64     `json:"exp"`
	Iat      int64     `json:"iat"`
}

// TokenService define a interface para serviços de token
type TokenService interface {
	GenerateToken(userInfo UserInfo) (string, error)
	ValidateToken(token string) (*Claims, error)
	GenerateRefreshToken(userID uuid.UUID) (string, error)
	ValidateRefreshToken(refreshToken string) (uuid.UUID, error)
}

// AuthRepository define a interface para repositório de autenticação
type AuthRepository interface {
	FindUserByEmail(email string) (*UserInfo, error)
	ValidatePassword(email, password string) error
}
