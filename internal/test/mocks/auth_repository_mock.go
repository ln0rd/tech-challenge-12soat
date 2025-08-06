package mocks

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/auth"
)

// AuthRepositoryMock implementa AuthRepository para testes
type AuthRepositoryMock struct {
	FindUserByEmailFunc  func(email string) (*domain.UserInfo, error)
	ValidatePasswordFunc func(email, password string) error
}

// FindUserByEmail chama a função mock
func (m *AuthRepositoryMock) FindUserByEmail(email string) (*domain.UserInfo, error) {
	if m.FindUserByEmailFunc != nil {
		return m.FindUserByEmailFunc(email)
	}
	return nil, nil
}

// ValidatePassword chama a função mock
func (m *AuthRepositoryMock) ValidatePassword(email, password string) error {
	if m.ValidatePasswordFunc != nil {
		return m.ValidatePasswordFunc(email, password)
	}
	return nil
}
