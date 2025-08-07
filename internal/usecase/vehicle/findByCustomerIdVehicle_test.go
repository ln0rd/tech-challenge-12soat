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

func TestFindByCustomerIdVehicle_Process_Success(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	customerID := uuid.New()
	mockVehicles := []models.Vehicle{
		{
			ID:                          uuid.New(),
			CustomerID:                  customerID,
			Model:                       "Corolla",
			Brand:                       "Toyota",
			ReleaseYear:                 2020,
			VehicleIdentificationNumber: "VIN123456789",
			NumberPlate:                 "ABC1234",
			Color:                       "Prata",
			CreatedAt:                   time.Now(),
			UpdatedAt:                   time.Now(),
		},
		{
			ID:                          uuid.New(),
			CustomerID:                  customerID,
			Model:                       "Civic",
			Brand:                       "Honda",
			ReleaseYear:                 2019,
			VehicleIdentificationNumber: "VIN987654321",
			NumberPlate:                 "XYZ5678",
			Color:                       "Preto",
			CreatedAt:                   time.Now(),
			UpdatedAt:                   time.Now(),
		},
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository
	vehicleRepoMock.FindByCustomerIDFunc = func(customerID uuid.UUID) ([]models.Vehicle, error) {
		return mockVehicles, nil
	}

	useCase := &FindByCustomerIdVehicle{
		VehicleRepository: vehicleRepoMock,
		Logger:            loggerMock,
	}

	// Act
	result, err := useCase.Process(customerID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 vehicles, got %d", len(result))
	}

	if result[0].Model != "Corolla" {
		t.Errorf("Expected first vehicle model 'Corolla', got '%s'", result[0].Model)
	}

	if result[1].Model != "Civic" {
		t.Errorf("Expected second vehicle model 'Civic', got '%s'", result[1].Model)
	}
}

func TestFindByCustomerIdVehicle_Process_EmptyResult(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	customerID := uuid.New()

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository to return empty list
	vehicleRepoMock.FindByCustomerIDFunc = func(customerID uuid.UUID) ([]models.Vehicle, error) {
		return []models.Vehicle{}, nil
	}

	useCase := &FindByCustomerIdVehicle{
		VehicleRepository: vehicleRepoMock,
		Logger:            loggerMock,
	}

	// Act
	result, err := useCase.Process(customerID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 vehicles, got %d", len(result))
	}
}

func TestFindByCustomerIdVehicle_Process_DatabaseError(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	customerID := uuid.New()

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository to return error
	vehicleRepoMock.FindByCustomerIDFunc = func(customerID uuid.UUID) ([]models.Vehicle, error) {
		return []models.Vehicle{}, errors.New("database error")
	}

	useCase := &FindByCustomerIdVehicle{
		VehicleRepository: vehicleRepoMock,
		Logger:            loggerMock,
	}

	// Act
	result, err := useCase.Process(customerID)

	// Assert
	if err == nil {
		t.Error("Expected error for database error, got nil")
	}

	if len(result) != 0 {
		t.Errorf("Expected empty result for error, got %d vehicles", len(result))
	}

	if err.Error() != "database error" {
		t.Errorf("Expected error message 'database error', got '%s'", err.Error())
	}
}

func TestFindByCustomerIdVehicle_FetchVehiclesFromDB_Success(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	customerID := uuid.New()
	mockVehicles := []models.Vehicle{
		{
			ID:                          uuid.New(),
			CustomerID:                  customerID,
			Model:                       "Corolla",
			Brand:                       "Toyota",
			ReleaseYear:                 2020,
			VehicleIdentificationNumber: "VIN123456789",
			NumberPlate:                 "ABC1234",
			Color:                       "Prata",
			CreatedAt:                   time.Now(),
			UpdatedAt:                   time.Now(),
		},
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository
	vehicleRepoMock.FindByCustomerIDFunc = func(customerID uuid.UUID) ([]models.Vehicle, error) {
		return mockVehicles, nil
	}

	useCase := &FindByCustomerIdVehicle{
		VehicleRepository: vehicleRepoMock,
		Logger:            loggerMock,
	}

	// Act
	result, err := useCase.FetchVehiclesFromDB(customerID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 vehicle, got %d", len(result))
	}

	if result[0].Model != "Corolla" {
		t.Errorf("Expected vehicle model 'Corolla', got '%s'", result[0].Model)
	}
}

func TestFindByCustomerIdVehicle_FetchVehiclesFromDB_Error(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	customerID := uuid.New()

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository to return error
	vehicleRepoMock.FindByCustomerIDFunc = func(customerID uuid.UUID) ([]models.Vehicle, error) {
		return []models.Vehicle{}, errors.New("database error")
	}

	useCase := &FindByCustomerIdVehicle{
		VehicleRepository: vehicleRepoMock,
		Logger:            loggerMock,
	}

	// Act
	result, err := useCase.FetchVehiclesFromDB(customerID)

	// Assert
	if err == nil {
		t.Error("Expected error for database error, got nil")
	}

	if len(result) != 0 {
		t.Errorf("Expected empty result for error, got %d vehicles", len(result))
	}

	if err.Error() != "database error" {
		t.Errorf("Expected error message 'database error', got '%s'", err.Error())
	}
}
