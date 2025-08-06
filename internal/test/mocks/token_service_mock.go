package mocks

import (
	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/auth"
)

// TokenServiceMock implementa TokenService para testes
type TokenServiceMock struct {
	GenerateTokenFunc        func(userInfo domain.UserInfo) (string, error)
	ValidateTokenFunc        func(token string) (*domain.Claims, error)
	GenerateRefreshTokenFunc func(userID uuid.UUID) (string, error)
	ValidateRefreshTokenFunc func(refreshToken string) (uuid.UUID, error)
}

// GenerateToken chama a função mock
func (m *TokenServiceMock) GenerateToken(userInfo domain.UserInfo) (string, error) {
	if m.GenerateTokenFunc != nil {
		return m.GenerateTokenFunc(userInfo)
	}
	return "", nil
}

// ValidateToken chama a função mock
func (m *TokenServiceMock) ValidateToken(token string) (*domain.Claims, error) {
	if m.ValidateTokenFunc != nil {
		return m.ValidateTokenFunc(token)
	}
	return nil, nil
}

// GenerateRefreshToken chama a função mock
func (m *TokenServiceMock) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	if m.GenerateRefreshTokenFunc != nil {
		return m.GenerateRefreshTokenFunc(userID)
	}
	return "", nil
}

// ValidateRefreshToken chama a função mock
func (m *TokenServiceMock) ValidateRefreshToken(refreshToken string) (uuid.UUID, error) {
	if m.ValidateRefreshTokenFunc != nil {
		return m.ValidateRefreshTokenFunc(refreshToken)
	}
	return uuid.Nil, nil
}
