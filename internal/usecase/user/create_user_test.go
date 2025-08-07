package user

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/user"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/test/mocks"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestCreateUser_Process_Success(t *testing.T) {
	// Arrange
	userRepoMock := &mocks.UserRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	userEntity := &domain.User{
		ID:        uuid.New(),
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock UserRepository
	userRepoMock.FindByEmailFunc = func(email string) (*models.User, error) {
		return nil, gorm.ErrRecordNotFound
	}

	userRepoMock.CreateFunc = func(user *models.User) error {
		return nil
	}

	useCase := &CreateUser{
		UserRepository: userRepoMock,
		Logger:         loggerMock,
	}

	// Act
	err := useCase.Process(userEntity)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verifica se a senha foi hasheada
	if userEntity.Password == "password123" {
		t.Error("Expected password to be hashed, got original password")
	}
}

func TestCreateUser_Process_EmailAlreadyExists(t *testing.T) {
	// Arrange
	userRepoMock := &mocks.UserRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	userEntity := &domain.User{
		ID:        uuid.New(),
		Username:  "testuser",
		Email:     "existing@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	existingUser := &models.User{
		ID:        uuid.New(),
		Username:  "existinguser",
		Email:     "existing@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock UserRepository to return existing user
	userRepoMock.FindByEmailFunc = func(email string) (*models.User, error) {
		return existingUser, nil
	}

	useCase := &CreateUser{
		UserRepository: userRepoMock,
		Logger:         loggerMock,
	}

	// Act
	err := useCase.Process(userEntity)

	// Assert
	if err == nil {
		t.Error("Expected error for existing email, got nil")
	}

	if err.Error() != "email already exists" {
		t.Errorf("Expected error message 'email already exists', got '%s'", err.Error())
	}
}

func TestCreateUser_Process_EmailCheckError(t *testing.T) {
	// Arrange
	userRepoMock := &mocks.UserRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	userEntity := &domain.User{
		ID:        uuid.New(),
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock UserRepository to return error
	userRepoMock.FindByEmailFunc = func(email string) (*models.User, error) {
		return nil, errors.New("database error")
	}

	useCase := &CreateUser{
		UserRepository: userRepoMock,
		Logger:         loggerMock,
	}

	// Act
	err := useCase.Process(userEntity)

	// Assert
	if err == nil {
		t.Error("Expected error for database error, got nil")
	}

	if err.Error() != "database error" {
		t.Errorf("Expected error message 'database error', got '%s'", err.Error())
	}
}

func TestCreateUser_Process_DatabaseError(t *testing.T) {
	// Arrange
	userRepoMock := &mocks.UserRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	userEntity := &domain.User{
		ID:        uuid.New(),
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock UserRepository
	userRepoMock.FindByEmailFunc = func(email string) (*models.User, error) {
		return nil, gorm.ErrRecordNotFound
	}

	userRepoMock.CreateFunc = func(user *models.User) error {
		return errors.New("database error")
	}

	useCase := &CreateUser{
		UserRepository: userRepoMock,
		Logger:         loggerMock,
	}

	// Act
	err := useCase.Process(userEntity)

	// Assert
	if err == nil {
		t.Error("Expected error for database error, got nil")
	}

	if err.Error() != "database error" {
		t.Errorf("Expected error message 'database error', got '%s'", err.Error())
	}
}

func TestCreateUser_ValidateEmailUniqueness_Success(t *testing.T) {
	// Arrange
	userRepoMock := &mocks.UserRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	email := "test@example.com"

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock UserRepository to return not found
	userRepoMock.FindByEmailFunc = func(email string) (*models.User, error) {
		return nil, gorm.ErrRecordNotFound
	}

	useCase := &CreateUser{
		UserRepository: userRepoMock,
		Logger:         loggerMock,
	}

	// Act
	err := useCase.ValidateEmailUniqueness(email)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCreateUser_ValidateEmailUniqueness_EmailExists(t *testing.T) {
	// Arrange
	userRepoMock := &mocks.UserRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	email := "existing@example.com"
	existingUser := &models.User{
		ID:        uuid.New(),
		Username:  "existinguser",
		Email:     email,
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock UserRepository to return existing user
	userRepoMock.FindByEmailFunc = func(email string) (*models.User, error) {
		return existingUser, nil
	}

	useCase := &CreateUser{
		UserRepository: userRepoMock,
		Logger:         loggerMock,
	}

	// Act
	err := useCase.ValidateEmailUniqueness(email)

	// Assert
	if err == nil {
		t.Error("Expected error for existing email, got nil")
	}

	if err.Error() != "email already exists" {
		t.Errorf("Expected error message 'email already exists', got '%s'", err.Error())
	}
}

func TestCreateUser_ValidateEmailUniqueness_DatabaseError(t *testing.T) {
	// Arrange
	userRepoMock := &mocks.UserRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	email := "test@example.com"

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock UserRepository to return error
	userRepoMock.FindByEmailFunc = func(email string) (*models.User, error) {
		return nil, errors.New("database error")
	}

	useCase := &CreateUser{
		UserRepository: userRepoMock,
		Logger:         loggerMock,
	}

	// Act
	err := useCase.ValidateEmailUniqueness(email)

	// Assert
	if err == nil {
		t.Error("Expected error for database error, got nil")
	}

	if err.Error() != "database error" {
		t.Errorf("Expected error message 'database error', got '%s'", err.Error())
	}
}

func TestCreateUser_HashPassword_Success(t *testing.T) {
	// Arrange
	loggerMock := &mocks.LoggerMock{}
	useCase := &CreateUser{
		Logger: loggerMock,
	}

	password := "password123"

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Act
	hashedPassword, err := useCase.HashPassword(password)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if hashedPassword == password {
		t.Error("Expected password to be hashed, got original password")
	}

	if len(hashedPassword) == 0 {
		t.Error("Expected hashed password to not be empty")
	}
}

func TestCreateUser_HashPassword_EmptyPassword(t *testing.T) {
	// Arrange
	loggerMock := &mocks.LoggerMock{}
	useCase := &CreateUser{
		Logger: loggerMock,
	}

	password := ""

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Act
	hashedPassword, err := useCase.HashPassword(password)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for empty password, got %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Error("Expected hashed password to not be empty")
	}
}

func TestCreateUser_SaveUserToDB_Success(t *testing.T) {
	// Arrange
	userRepoMock := &mocks.UserRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	userModel := &models.User{
		ID:        uuid.New(),
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock UserRepository
	userRepoMock.CreateFunc = func(user *models.User) error {
		return nil
	}

	useCase := &CreateUser{
		UserRepository: userRepoMock,
		Logger:         loggerMock,
	}

	// Act
	err := useCase.SaveUserToDB(userModel)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCreateUser_SaveUserToDB_Error(t *testing.T) {
	// Arrange
	userRepoMock := &mocks.UserRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	userModel := &models.User{
		ID:        uuid.New(),
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock UserRepository to return error
	userRepoMock.CreateFunc = func(user *models.User) error {
		return errors.New("database error")
	}

	useCase := &CreateUser{
		UserRepository: userRepoMock,
		Logger:         loggerMock,
	}

	// Act
	err := useCase.SaveUserToDB(userModel)

	// Assert
	if err == nil {
		t.Error("Expected error for database error, got nil")
	}

	if err.Error() != "database error" {
		t.Errorf("Expected error message 'database error', got '%s'", err.Error())
	}
}

func TestCreateUser_Process_PasswordHashing(t *testing.T) {
	// Arrange
	userRepoMock := &mocks.UserRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	originalPassword := "password123"
	userEntity := &domain.User{
		ID:        uuid.New(),
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  originalPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock UserRepository
	userRepoMock.FindByEmailFunc = func(email string) (*models.User, error) {
		return nil, gorm.ErrRecordNotFound
	}

	userRepoMock.CreateFunc = func(user *models.User) error {
		return nil
	}

	useCase := &CreateUser{
		UserRepository: userRepoMock,
		Logger:         loggerMock,
	}

	// Act
	err := useCase.Process(userEntity)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verifica se a senha foi hasheada
	if userEntity.Password == originalPassword {
		t.Error("Expected password to be hashed, got original password")
	}

	// Verifica se o hash é diferente do original
	if len(userEntity.Password) <= len(originalPassword) {
		t.Error("Expected hashed password to be longer than original")
	}

	// Verifica se o hash começa com $2a$ (bcrypt)
	if len(userEntity.Password) < 4 || userEntity.Password[:4] != "$2a$" {
		t.Error("Expected password to be hashed with bcrypt")
	}
}
