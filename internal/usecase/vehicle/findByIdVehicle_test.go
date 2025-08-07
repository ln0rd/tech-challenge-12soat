package vehicle

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/test/mocks"
	"go.uber.org/zap"
)

func TestFindByIdVehicle_Process_Success(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	vehicleID := uuid.New()
	mockVehicle := &models.Vehicle{
		ID:                          vehicleID,
		CustomerID:                  uuid.New(),
		Model:                       "Corolla",
		Brand:                       "Toyota",
		ReleaseYear:                 2020,
		VehicleIdentificationNumber: "VIN123456789",
		NumberPlate:                 "ABC1234",
		Color:                       "Prata",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository
	vehicleRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Vehicle, error) {
		return mockVehicle, nil
	}

	useCase := &FindByIdVehicle{
		VehicleRepository: vehicleRepoMock,
		Logger:            loggerMock,
	}

	// Act
	result, err := useCase.Process(vehicleID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected vehicle, got nil")
	}

	if result.ID != vehicleID {
		t.Errorf("Expected vehicle ID %s, got %s", vehicleID, result.ID)
	}

	if result.Model != "Corolla" {
		t.Errorf("Expected vehicle model 'Corolla', got '%s'", result.Model)
	}

	if result.Brand != "Toyota" {
		t.Errorf("Expected vehicle brand 'Toyota', got '%s'", result.Brand)
	}

	if result.NumberPlate != "ABC1234" {
		t.Errorf("Expected vehicle number plate 'ABC1234', got '%s'", result.NumberPlate)
	}
}

func TestFindByIdVehicle_Process_VehicleNotFound(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	vehicleID := uuid.New()

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository to return error
	vehicleRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Vehicle, error) {
		return nil, errors.New("vehicle not found")
	}

	useCase := &FindByIdVehicle{
		VehicleRepository: vehicleRepoMock,
		Logger:            loggerMock,
	}

	// Act
	result, err := useCase.Process(vehicleID)

	// Assert
	if err == nil {
		t.Error("Expected error for vehicle not found, got nil")
	}

	if result != nil {
		t.Error("Expected nil result for vehicle not found, got vehicle")
	}

	if err.Error() != "vehicle not found" {
		t.Errorf("Expected error message 'vehicle not found', got '%s'", err.Error())
	}
}

func TestFindByIdVehicle_FetchVehicleFromDB_Success(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	vehicleID := uuid.New()
	mockVehicle := &models.Vehicle{
		ID:                          vehicleID,
		CustomerID:                  uuid.New(),
		Model:                       "Corolla",
		Brand:                       "Toyota",
		ReleaseYear:                 2020,
		VehicleIdentificationNumber: "VIN123456789",
		NumberPlate:                 "ABC1234",
		Color:                       "Prata",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository
	vehicleRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Vehicle, error) {
		return mockVehicle, nil
	}

	useCase := &FindByIdVehicle{
		VehicleRepository: vehicleRepoMock,
		Logger:            loggerMock,
	}

	// Act
	result, err := useCase.FetchVehicleFromDB(vehicleID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected vehicle, got nil")
	}

	if result.ID != vehicleID {
		t.Errorf("Expected vehicle ID %s, got %s", vehicleID, result.ID)
	}

	if result.Model != "Corolla" {
		t.Errorf("Expected vehicle model 'Corolla', got '%s'", result.Model)
	}
}

func TestFindByIdVehicle_FetchVehicleFromDB_Error(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	vehicleID := uuid.New()

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository to return error
	vehicleRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Vehicle, error) {
		return nil, errors.New("database error")
	}

	useCase := &FindByIdVehicle{
		VehicleRepository: vehicleRepoMock,
		Logger:            loggerMock,
	}

	// Act
	result, err := useCase.FetchVehicleFromDB(vehicleID)

	// Assert
	if err == nil {
		t.Error("Expected error for database error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result for database error, got vehicle")
	}

	if err.Error() != "database error" {
		t.Errorf("Expected error message 'database error', got '%s'", err.Error())
	}
}
