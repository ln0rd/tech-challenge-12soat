package vehicle

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/test/mocks"
	"go.uber.org/zap"
)

func TestDeleteByIdVehicle_Process_Success(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	vehicleID := uuid.New()

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository
	vehicleRepoMock.DeleteFunc = func(id uuid.UUID) error {
		return nil
	}

	useCase := &DeleteByIdVehicle{
		VehicleRepository: vehicleRepoMock,
		Logger:            loggerMock,
	}

	// Act
	err := useCase.Process(vehicleID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDeleteByIdVehicle_Process_DatabaseError(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	vehicleID := uuid.New()

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository to return error
	vehicleRepoMock.DeleteFunc = func(id uuid.UUID) error {
		return errors.New("database error")
	}

	useCase := &DeleteByIdVehicle{
		VehicleRepository: vehicleRepoMock,
		Logger:            loggerMock,
	}

	// Act
	err := useCase.Process(vehicleID)

	// Assert
	if err == nil {
		t.Error("Expected error for database error, got nil")
	}

	if err.Error() != "database error" {
		t.Errorf("Expected error message 'database error', got '%s'", err.Error())
	}
}

func TestDeleteByIdVehicle_DeleteVehicleFromDB_Success(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	vehicleID := uuid.New()

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository
	vehicleRepoMock.DeleteFunc = func(id uuid.UUID) error {
		return nil
	}

	useCase := &DeleteByIdVehicle{
		VehicleRepository: vehicleRepoMock,
		Logger:            loggerMock,
	}

	// Act
	err := useCase.DeleteVehicleFromDB(vehicleID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDeleteByIdVehicle_DeleteVehicleFromDB_Error(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	vehicleID := uuid.New()

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository to return error
	vehicleRepoMock.DeleteFunc = func(id uuid.UUID) error {
		return errors.New("database error")
	}

	useCase := &DeleteByIdVehicle{
		VehicleRepository: vehicleRepoMock,
		Logger:            loggerMock,
	}

	// Act
	err := useCase.DeleteVehicleFromDB(vehicleID)

	// Assert
	if err == nil {
		t.Error("Expected error for database error, got nil")
	}

	if err.Error() != "database error" {
		t.Errorf("Expected error message 'database error', got '%s'", err.Error())
	}
}
