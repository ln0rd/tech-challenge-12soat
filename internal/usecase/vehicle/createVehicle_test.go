package vehicle

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	vehicleDomain "github.com/ln0rd/tech_challenge_12soat/internal/domain/vehicle"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/test/mocks"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestCreateVehicle_Process_Success(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	vehicleEntity := &vehicleDomain.Vehicle{
		ID:                          uuid.New(),
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

	mockCustomer := &models.Customer{
		ID:             vehicleEntity.CustomerID,
		Name:           "John Doe",
		DocumentNumber: "12345678901",
		CustomerType:   "Individual",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository
	vehicleRepoMock.FindByNumberPlateFunc = func(numberPlate string) (*models.Vehicle, error) {
		return nil, gorm.ErrRecordNotFound
	}

	vehicleRepoMock.CreateFunc = func(vehicle *models.Vehicle) error {
		return nil
	}

	// Mock CustomerRepository
	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		return mockCustomer, nil
	}

	useCase := &CreateVehicle{
		VehicleRepository:  vehicleRepoMock,
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.Process(vehicleEntity)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCreateVehicle_Process_NumberPlateAlreadyExists(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	vehicleEntity := &vehicleDomain.Vehicle{
		ID:                          uuid.New(),
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

	existingVehicle := &models.Vehicle{
		ID:                          uuid.New(),
		CustomerID:                  uuid.New(),
		Model:                       "Civic",
		Brand:                       "Honda",
		ReleaseYear:                 2019,
		VehicleIdentificationNumber: "VIN987654321",
		NumberPlate:                 "ABC1234",
		Color:                       "Preto",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository to return existing vehicle
	vehicleRepoMock.FindByNumberPlateFunc = func(numberPlate string) (*models.Vehicle, error) {
		return existingVehicle, nil
	}

	useCase := &CreateVehicle{
		VehicleRepository:  vehicleRepoMock,
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.Process(vehicleEntity)

	// Assert
	if err == nil {
		t.Error("Expected error for existing number plate, got nil")
	}

	if err.Error() != "number plate already exists" {
		t.Errorf("Expected error message 'number plate already exists', got '%s'", err.Error())
	}
}

func TestCreateVehicle_Process_CustomerNotFound(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	vehicleEntity := &vehicleDomain.Vehicle{
		ID:                          uuid.New(),
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
	vehicleRepoMock.FindByNumberPlateFunc = func(numberPlate string) (*models.Vehicle, error) {
		return nil, gorm.ErrRecordNotFound
	}

	// Mock CustomerRepository to return error
	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		return nil, errors.New("customer not found")
	}

	useCase := &CreateVehicle{
		VehicleRepository:  vehicleRepoMock,
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.Process(vehicleEntity)

	// Assert
	if err == nil {
		t.Error("Expected error for customer not found, got nil")
	}

	if err.Error() != "customer not found" {
		t.Errorf("Expected error message 'customer not found', got '%s'", err.Error())
	}
}

func TestCreateVehicle_Process_DatabaseError(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	vehicleEntity := &vehicleDomain.Vehicle{
		ID:                          uuid.New(),
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

	mockCustomer := &models.Customer{
		ID:             vehicleEntity.CustomerID,
		Name:           "John Doe",
		DocumentNumber: "12345678901",
		CustomerType:   "Individual",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository
	vehicleRepoMock.FindByNumberPlateFunc = func(numberPlate string) (*models.Vehicle, error) {
		return nil, gorm.ErrRecordNotFound
	}

	vehicleRepoMock.CreateFunc = func(vehicle *models.Vehicle) error {
		return errors.New("database error")
	}

	// Mock CustomerRepository
	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		return mockCustomer, nil
	}

	useCase := &CreateVehicle{
		VehicleRepository:  vehicleRepoMock,
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.Process(vehicleEntity)

	// Assert
	if err == nil {
		t.Error("Expected error for database error, got nil")
	}

	if err.Error() != "database error" {
		t.Errorf("Expected error message 'database error', got '%s'", err.Error())
	}
}

func TestCreateVehicle_ValidateNumberPlateUniqueness_Success(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	numberPlate := "ABC1234"

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository to return not found
	vehicleRepoMock.FindByNumberPlateFunc = func(numberPlate string) (*models.Vehicle, error) {
		return nil, gorm.ErrRecordNotFound
	}

	useCase := &CreateVehicle{
		VehicleRepository: vehicleRepoMock,
		Logger:            loggerMock,
	}

	// Act
	err := useCase.ValidateNumberPlateUniqueness(numberPlate)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCreateVehicle_ValidateNumberPlateUniqueness_AlreadyExists(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	numberPlate := "ABC1234"
	existingVehicle := &models.Vehicle{
		ID:                          uuid.New(),
		CustomerID:                  uuid.New(),
		Model:                       "Civic",
		Brand:                       "Honda",
		ReleaseYear:                 2019,
		VehicleIdentificationNumber: "VIN987654321",
		NumberPlate:                 numberPlate,
		Color:                       "Preto",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository to return existing vehicle
	vehicleRepoMock.FindByNumberPlateFunc = func(numberPlate string) (*models.Vehicle, error) {
		return existingVehicle, nil
	}

	useCase := &CreateVehicle{
		VehicleRepository: vehicleRepoMock,
		Logger:            loggerMock,
	}

	// Act
	err := useCase.ValidateNumberPlateUniqueness(numberPlate)

	// Assert
	if err == nil {
		t.Error("Expected error for existing number plate, got nil")
	}

	if err.Error() != "number plate already exists" {
		t.Errorf("Expected error message 'number plate already exists', got '%s'", err.Error())
	}
}

func TestCreateVehicle_ValidateNumberPlateUniqueness_DatabaseError(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	numberPlate := "ABC1234"

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository to return error
	vehicleRepoMock.FindByNumberPlateFunc = func(numberPlate string) (*models.Vehicle, error) {
		return nil, errors.New("database error")
	}

	useCase := &CreateVehicle{
		VehicleRepository: vehicleRepoMock,
		Logger:            loggerMock,
	}

	// Act
	err := useCase.ValidateNumberPlateUniqueness(numberPlate)

	// Assert
	if err == nil {
		t.Error("Expected error for database error, got nil")
	}

	if err.Error() != "database error" {
		t.Errorf("Expected error message 'database error', got '%s'", err.Error())
	}
}

func TestCreateVehicle_ValidateCustomerExists_Success(t *testing.T) {
	// Arrange
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	customerID := uuid.New()
	mockCustomer := &models.Customer{
		ID:             customerID,
		Name:           "John Doe",
		DocumentNumber: "12345678901",
		CustomerType:   "Individual",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock CustomerRepository
	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		return mockCustomer, nil
	}

	useCase := &CreateVehicle{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.ValidateCustomerExists(customerID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCreateVehicle_ValidateCustomerExists_NilCustomerID(t *testing.T) {
	// Arrange
	loggerMock := &mocks.LoggerMock{}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	useCase := &CreateVehicle{
		Logger: loggerMock,
	}

	// Act
	err := useCase.ValidateCustomerExists(uuid.Nil)

	// Assert
	if err == nil {
		t.Error("Expected error for nil customer ID, got nil")
	}

	if err.Error() != "customer ID is required" {
		t.Errorf("Expected error message 'customer ID is required', got '%s'", err.Error())
	}
}

func TestCreateVehicle_ValidateCustomerExists_CustomerNotFound(t *testing.T) {
	// Arrange
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	customerID := uuid.New()

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock CustomerRepository to return error
	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		return nil, errors.New("customer not found")
	}

	useCase := &CreateVehicle{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.ValidateCustomerExists(customerID)

	// Assert
	if err == nil {
		t.Error("Expected error for customer not found, got nil")
	}

	if err.Error() != "customer not found" {
		t.Errorf("Expected error message 'customer not found', got '%s'", err.Error())
	}
}

func TestCreateVehicle_SaveVehicleToDB_Success(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	vehicleModel := &models.Vehicle{
		ID:                          uuid.New(),
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
	vehicleRepoMock.CreateFunc = func(vehicle *models.Vehicle) error {
		return nil
	}

	useCase := &CreateVehicle{
		VehicleRepository: vehicleRepoMock,
		Logger:            loggerMock,
	}

	// Act
	err := useCase.SaveVehicleToDB(vehicleModel)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCreateVehicle_SaveVehicleToDB_Error(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	vehicleModel := &models.Vehicle{
		ID:                          uuid.New(),
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

	// Mock VehicleRepository to return error
	vehicleRepoMock.CreateFunc = func(vehicle *models.Vehicle) error {
		return errors.New("database error")
	}

	useCase := &CreateVehicle{
		VehicleRepository: vehicleRepoMock,
		Logger:            loggerMock,
	}

	// Act
	err := useCase.SaveVehicleToDB(vehicleModel)

	// Assert
	if err == nil {
		t.Error("Expected error for database error, got nil")
	}

	if err.Error() != "database error" {
		t.Errorf("Expected error message 'database error', got '%s'", err.Error())
	}
}
