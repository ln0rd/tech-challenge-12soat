package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/auth"
	"github.com/ln0rd/tech_challenge_12soat/internal/test/mocks"
	"go.uber.org/zap"
)

func TestLoginUseCase_Execute_Success(t *testing.T) {
	// Arrange
	authRepoMock := &mocks.AuthRepositoryMock{}
	tokenServiceMock := &mocks.TokenServiceMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	userID := uuid.New()
	mockUserInfo := &domain.UserInfo{
		ID:       userID,
		Email:    "joao@example.com",
		Username: "joao",
		UserType: "mechanic",
	}

	authRepoMock.FindUserByEmailFunc = func(email string) (*domain.UserInfo, error) {
		if email == "joao@example.com" {
			return mockUserInfo, nil
		}
		return nil, errors.New("user not found")
	}

	authRepoMock.ValidatePasswordFunc = func(email, password string) error {
		if email == "joao@example.com" && password == "password123" {
			return nil
		}
		return errors.New("invalid password")
	}

	tokenServiceMock.GenerateTokenFunc = func(userInfo domain.UserInfo) (string, error) {
		return "mock-jwt-token", nil
	}

	tokenServiceMock.GenerateRefreshTokenFunc = func(userID uuid.UUID) (string, error) {
		return "mock-refresh-token", nil
	}

	useCase := &LoginUseCase{
		authRepository: authRepoMock,
		tokenService:   tokenServiceMock,
		logger:         loggerMock,
	}

	request := domain.LoginRequest{
		Email:    "joao@example.com",
		Password: "password123",
	}

	// Act
	result, err := useCase.Execute(request)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected result, got nil")
	}

	if result.Token != "mock-jwt-token" {
		t.Errorf("Expected token 'mock-jwt-token', got '%s'", result.Token)
	}

	if result.RefreshToken != "mock-refresh-token" {
		t.Errorf("Expected refresh token 'mock-refresh-token', got '%s'", result.RefreshToken)
	}

	if result.User.ID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, result.User.ID)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing login request",
		"User found",
		"Password validated successfully",
		"Login successful",
	}

	for _, expectedLog := range expectedInfoLogs {
		found := false
		for _, actualLog := range loggedInfo {
			if actualLog == expectedLog {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected log message '%s' not found", expectedLog)
		}
	}

	if len(loggedErrors) > 0 {
		t.Errorf("Expected no error logs, got %d", len(loggedErrors))
	}
}

func TestLoginUseCase_Execute_UserNotFound(t *testing.T) {
	// Arrange
	authRepoMock := &mocks.AuthRepositoryMock{}
	tokenServiceMock := &mocks.TokenServiceMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	authRepoMock.FindUserByEmailFunc = func(email string) (*domain.UserInfo, error) {
		return nil, errors.New("user not found")
	}

	useCase := &LoginUseCase{
		authRepository: authRepoMock,
		tokenService:   tokenServiceMock,
		logger:         loggerMock,
	}

	request := domain.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	// Act
	result, err := useCase.Execute(request)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "invalid credentials" {
		t.Errorf("Expected error 'invalid credentials', got '%s'", err.Error())
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	// Verifica se os logs corretos foram chamados
	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Processing login request" {
		t.Errorf("Expected log message 'Processing login request', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "User not found" {
		t.Errorf("Expected error log 'User not found', got '%s'", loggedErrors[0])
	}
}

func TestLoginUseCase_Execute_InvalidPassword(t *testing.T) {
	// Arrange
	authRepoMock := &mocks.AuthRepositoryMock{}
	tokenServiceMock := &mocks.TokenServiceMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	userID := uuid.New()
	mockUserInfo := &domain.UserInfo{
		ID:       userID,
		Email:    "joao@example.com",
		Username: "joao",
		UserType: "mechanic",
	}

	authRepoMock.FindUserByEmailFunc = func(email string) (*domain.UserInfo, error) {
		if email == "joao@example.com" {
			return mockUserInfo, nil
		}
		return nil, errors.New("user not found")
	}

	authRepoMock.ValidatePasswordFunc = func(email, password string) error {
		return errors.New("invalid password")
	}

	useCase := &LoginUseCase{
		authRepository: authRepoMock,
		tokenService:   tokenServiceMock,
		logger:         loggerMock,
	}

	request := domain.LoginRequest{
		Email:    "joao@example.com",
		Password: "wrongpassword",
	}

	// Act
	result, err := useCase.Execute(request)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "invalid credentials" {
		t.Errorf("Expected error 'invalid credentials', got '%s'", err.Error())
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing login request",
		"User found",
	}

	for _, expectedLog := range expectedInfoLogs {
		found := false
		for _, actualLog := range loggedInfo {
			if actualLog == expectedLog {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected log message '%s' not found", expectedLog)
		}
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Invalid password" {
		t.Errorf("Expected error log 'Invalid password', got '%s'", loggedErrors[0])
	}
}

func TestLoginUseCase_Execute_TokenGenerationError(t *testing.T) {
	// Arrange
	authRepoMock := &mocks.AuthRepositoryMock{}
	tokenServiceMock := &mocks.TokenServiceMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	userID := uuid.New()
	mockUserInfo := &domain.UserInfo{
		ID:       userID,
		Email:    "joao@example.com",
		Username: "joao",
		UserType: "mechanic",
	}

	authRepoMock.FindUserByEmailFunc = func(email string) (*domain.UserInfo, error) {
		if email == "joao@example.com" {
			return mockUserInfo, nil
		}
		return nil, errors.New("user not found")
	}

	authRepoMock.ValidatePasswordFunc = func(email, password string) error {
		if email == "joao@example.com" && password == "password123" {
			return nil
		}
		return errors.New("invalid password")
	}

	tokenServiceMock.GenerateTokenFunc = func(userInfo domain.UserInfo) (string, error) {
		return "", errors.New("token generation failed")
	}

	useCase := &LoginUseCase{
		authRepository: authRepoMock,
		tokenService:   tokenServiceMock,
		logger:         loggerMock,
	}

	request := domain.LoginRequest{
		Email:    "joao@example.com",
		Password: "password123",
	}

	// Act
	result, err := useCase.Execute(request)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "error generating token" {
		t.Errorf("Expected error 'error generating token', got '%s'", err.Error())
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing login request",
		"User found",
		"Password validated successfully",
	}

	for _, expectedLog := range expectedInfoLogs {
		found := false
		for _, actualLog := range loggedInfo {
			if actualLog == expectedLog {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected log message '%s' not found", expectedLog)
		}
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Error generating token" {
		t.Errorf("Expected error log 'Error generating token', got '%s'", loggedErrors[0])
	}
}

func TestLoginUseCase_Execute_RefreshTokenGenerationError(t *testing.T) {
	// Arrange
	authRepoMock := &mocks.AuthRepositoryMock{}
	tokenServiceMock := &mocks.TokenServiceMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	userID := uuid.New()
	mockUserInfo := &domain.UserInfo{
		ID:       userID,
		Email:    "joao@example.com",
		Username: "joao",
		UserType: "mechanic",
	}

	authRepoMock.FindUserByEmailFunc = func(email string) (*domain.UserInfo, error) {
		if email == "joao@example.com" {
			return mockUserInfo, nil
		}
		return nil, errors.New("user not found")
	}

	authRepoMock.ValidatePasswordFunc = func(email, password string) error {
		if email == "joao@example.com" && password == "password123" {
			return nil
		}
		return errors.New("invalid password")
	}

	tokenServiceMock.GenerateTokenFunc = func(userInfo domain.UserInfo) (string, error) {
		return "mock-jwt-token", nil
	}

	tokenServiceMock.GenerateRefreshTokenFunc = func(userID uuid.UUID) (string, error) {
		return "", errors.New("refresh token generation failed")
	}

	useCase := &LoginUseCase{
		authRepository: authRepoMock,
		tokenService:   tokenServiceMock,
		logger:         loggerMock,
	}

	request := domain.LoginRequest{
		Email:    "joao@example.com",
		Password: "password123",
	}

	// Act
	result, err := useCase.Execute(request)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "error generating refresh token" {
		t.Errorf("Expected error 'error generating refresh token', got '%s'", err.Error())
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing login request",
		"User found",
		"Password validated successfully",
	}

	for _, expectedLog := range expectedInfoLogs {
		found := false
		for _, actualLog := range loggedInfo {
			if actualLog == expectedLog {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected log message '%s' not found", expectedLog)
		}
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Error generating refresh token" {
		t.Errorf("Expected error log 'Error generating refresh token', got '%s'", loggedErrors[0])
	}
}

func TestLoginUseCase_Execute_ExpiresAtCalculation(t *testing.T) {
	// Arrange
	authRepoMock := &mocks.AuthRepositoryMock{}
	tokenServiceMock := &mocks.TokenServiceMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	userID := uuid.New()
	mockUserInfo := &domain.UserInfo{
		ID:       userID,
		Email:    "joao@example.com",
		Username: "joao",
		UserType: "mechanic",
	}

	authRepoMock.FindUserByEmailFunc = func(email string) (*domain.UserInfo, error) {
		if email == "joao@example.com" {
			return mockUserInfo, nil
		}
		return nil, errors.New("user not found")
	}

	authRepoMock.ValidatePasswordFunc = func(email, password string) error {
		if email == "joao@example.com" && password == "password123" {
			return nil
		}
		return errors.New("invalid password")
	}

	tokenServiceMock.GenerateTokenFunc = func(userInfo domain.UserInfo) (string, error) {
		return "mock-jwt-token", nil
	}

	tokenServiceMock.GenerateRefreshTokenFunc = func(userID uuid.UUID) (string, error) {
		return "mock-refresh-token", nil
	}

	useCase := &LoginUseCase{
		authRepository: authRepoMock,
		tokenService:   tokenServiceMock,
		logger:         loggerMock,
	}

	request := domain.LoginRequest{
		Email:    "joao@example.com",
		Password: "password123",
	}

	// Act
	result, err := useCase.Execute(request)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected result, got nil")
	}

	// Verifica se a data de expiração está correta (24 horas no futuro)
	expectedExpiresAt := time.Now().Add(24 * time.Hour)
	timeDiff := result.ExpiresAt.Sub(expectedExpiresAt)
	if timeDiff < -time.Second || timeDiff > time.Second {
		t.Errorf("Expected expires_at to be approximately 24 hours in the future, got %v", result.ExpiresAt)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing login request",
		"User found",
		"Password validated successfully",
		"Login successful",
	}

	for _, expectedLog := range expectedInfoLogs {
		found := false
		for _, actualLog := range loggedInfo {
			if actualLog == expectedLog {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected log message '%s' not found", expectedLog)
		}
	}
}
