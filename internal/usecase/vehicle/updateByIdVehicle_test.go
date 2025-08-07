package vehicle

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/vehicle"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/test/mocks"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestUpdateByIdVehicle_Process_Success(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	vehicleID := uuid.New()
	customerID := uuid.New()

	existingVehicle := &models.Vehicle{
		ID:                          vehicleID,
		CustomerID:                  customerID,
		Model:                       "Corolla",
		Brand:                       "Toyota",
		ReleaseYear:                 2020,
		VehicleIdentificationNumber: "VIN123456789",
		NumberPlate:                 "ABC1234",
		Color:                       "Prata",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	vehicleEntity := &domain.Vehicle{
		ID:                          vehicleID,
		CustomerID:                  customerID,
		Model:                       "Corolla XSE",
		Brand:                       "Toyota",
		ReleaseYear:                 2021,
		VehicleIdentificationNumber: "VIN123456789",
		NumberPlate:                 "ABC1234",
		Color:                       "Branco",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

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

	// Mock VehicleRepository
	vehicleRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Vehicle, error) {
		return existingVehicle, nil
	}

	vehicleRepoMock.FindByNumberPlateFunc = func(numberPlate string) (*models.Vehicle, error) {
		return nil, gorm.ErrRecordNotFound
	}

	vehicleRepoMock.UpdateFunc = func(vehicle *models.Vehicle) error {
		return nil
	}

	// Mock CustomerRepository
	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		return mockCustomer, nil
	}

	useCase := &UpdateByIdVehicle{
		VehicleRepository:  vehicleRepoMock,
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.Process(vehicleID, vehicleEntity)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestUpdateByIdVehicle_Process_VehicleNotFound(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	vehicleID := uuid.New()
	vehicleEntity := &domain.Vehicle{
		ID:                          vehicleID,
		CustomerID:                  uuid.New(),
		Model:                       "Corolla XSE",
		Brand:                       "Toyota",
		ReleaseYear:                 2021,
		VehicleIdentificationNumber: "VIN123456789",
		NumberPlate:                 "ABC1234",
		Color:                       "Branco",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository to return error
	vehicleRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Vehicle, error) {
		return nil, errors.New("vehicle not found")
	}

	useCase := &UpdateByIdVehicle{
		VehicleRepository:  vehicleRepoMock,
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.Process(vehicleID, vehicleEntity)

	// Assert
	if err == nil {
		t.Error("Expected error for vehicle not found, got nil")
	}

	if err.Error() != "vehicle not found" {
		t.Errorf("Expected error message 'vehicle not found', got '%s'", err.Error())
	}
}

func TestUpdateByIdVehicle_Process_NumberPlateAlreadyExists(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	vehicleID := uuid.New()
	customerID := uuid.New()

	existingVehicle := &models.Vehicle{
		ID:                          vehicleID,
		CustomerID:                  customerID,
		Model:                       "Corolla",
		Brand:                       "Toyota",
		ReleaseYear:                 2020,
		VehicleIdentificationNumber: "VIN123456789",
		NumberPlate:                 "ABC1234",
		Color:                       "Prata",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	vehicleEntity := &domain.Vehicle{
		ID:                          vehicleID,
		CustomerID:                  customerID,
		Model:                       "Corolla XSE",
		Brand:                       "Toyota",
		ReleaseYear:                 2021,
		VehicleIdentificationNumber: "VIN123456789",
		NumberPlate:                 "XYZ5678", // Different number plate
		Color:                       "Branco",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	existingVehicleWithSamePlate := &models.Vehicle{
		ID:                          uuid.New(),
		CustomerID:                  uuid.New(),
		Model:                       "Civic",
		Brand:                       "Honda",
		ReleaseYear:                 2019,
		VehicleIdentificationNumber: "VIN987654321",
		NumberPlate:                 "XYZ5678",
		Color:                       "Preto",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository
	vehicleRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Vehicle, error) {
		return existingVehicle, nil
	}

	vehicleRepoMock.FindByNumberPlateFunc = func(numberPlate string) (*models.Vehicle, error) {
		return existingVehicleWithSamePlate, nil
	}

	useCase := &UpdateByIdVehicle{
		VehicleRepository:  vehicleRepoMock,
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.Process(vehicleID, vehicleEntity)

	// Assert
	if err == nil {
		t.Error("Expected error for existing number plate, got nil")
	}

	if err.Error() != "number plate already exists" {
		t.Errorf("Expected error message 'number plate already exists', got '%s'", err.Error())
	}
}

func TestUpdateByIdVehicle_Process_CustomerNotFound(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	vehicleID := uuid.New()
	customerID := uuid.New()

	existingVehicle := &models.Vehicle{
		ID:                          vehicleID,
		CustomerID:                  customerID,
		Model:                       "Corolla",
		Brand:                       "Toyota",
		ReleaseYear:                 2020,
		VehicleIdentificationNumber: "VIN123456789",
		NumberPlate:                 "ABC1234",
		Color:                       "Prata",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	vehicleEntity := &domain.Vehicle{
		ID:                          vehicleID,
		CustomerID:                  customerID,
		Model:                       "Corolla XSE",
		Brand:                       "Toyota",
		ReleaseYear:                 2021,
		VehicleIdentificationNumber: "VIN123456789",
		NumberPlate:                 "ABC1234",
		Color:                       "Branco",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository
	vehicleRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Vehicle, error) {
		return existingVehicle, nil
	}

	vehicleRepoMock.FindByNumberPlateFunc = func(numberPlate string) (*models.Vehicle, error) {
		return nil, gorm.ErrRecordNotFound
	}

	// Mock CustomerRepository to return error
	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		return nil, errors.New("customer not found")
	}

	useCase := &UpdateByIdVehicle{
		VehicleRepository:  vehicleRepoMock,
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.Process(vehicleID, vehicleEntity)

	// Assert
	if err == nil {
		t.Error("Expected error for customer not found, got nil")
	}

	if err.Error() != "customer not found" {
		t.Errorf("Expected error message 'customer not found', got '%s'", err.Error())
	}
}

func TestUpdateByIdVehicle_Process_DatabaseError(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	vehicleID := uuid.New()
	customerID := uuid.New()

	existingVehicle := &models.Vehicle{
		ID:                          vehicleID,
		CustomerID:                  customerID,
		Model:                       "Corolla",
		Brand:                       "Toyota",
		ReleaseYear:                 2020,
		VehicleIdentificationNumber: "VIN123456789",
		NumberPlate:                 "ABC1234",
		Color:                       "Prata",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	vehicleEntity := &domain.Vehicle{
		ID:                          vehicleID,
		CustomerID:                  customerID,
		Model:                       "Corolla XSE",
		Brand:                       "Toyota",
		ReleaseYear:                 2021,
		VehicleIdentificationNumber: "VIN123456789",
		NumberPlate:                 "ABC1234",
		Color:                       "Branco",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

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

	// Mock VehicleRepository
	vehicleRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Vehicle, error) {
		return existingVehicle, nil
	}

	vehicleRepoMock.FindByNumberPlateFunc = func(numberPlate string) (*models.Vehicle, error) {
		return nil, gorm.ErrRecordNotFound
	}

	vehicleRepoMock.UpdateFunc = func(vehicle *models.Vehicle) error {
		return errors.New("database error")
	}

	// Mock CustomerRepository
	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		return mockCustomer, nil
	}

	useCase := &UpdateByIdVehicle{
		VehicleRepository:  vehicleRepoMock,
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.Process(vehicleID, vehicleEntity)

	// Assert
	if err == nil {
		t.Error("Expected error for database error, got nil")
	}

	if err.Error() != "database error" {
		t.Errorf("Expected error message 'database error', got '%s'", err.Error())
	}
}

func TestUpdateByIdVehicle_FetchVehicleFromDB_Success(t *testing.T) {
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

	useCase := &UpdateByIdVehicle{
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

func TestUpdateByIdVehicle_FetchVehicleFromDB_Error(t *testing.T) {
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

	useCase := &UpdateByIdVehicle{
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

func TestUpdateByIdVehicle_ValidateNumberPlateUniqueness_Success(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	numberPlate := "ABC1234"
	vehicleID := uuid.New()

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository to return not found
	vehicleRepoMock.FindByNumberPlateFunc = func(numberPlate string) (*models.Vehicle, error) {
		return nil, gorm.ErrRecordNotFound
	}

	useCase := &UpdateByIdVehicle{
		VehicleRepository: vehicleRepoMock,
		Logger:            loggerMock,
	}

	// Act
	err := useCase.ValidateNumberPlateUniqueness(numberPlate, vehicleID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestUpdateByIdVehicle_ValidateNumberPlateUniqueness_AlreadyExists(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	numberPlate := "ABC1234"
	vehicleID := uuid.New()
	otherVehicleID := uuid.New()

	existingVehicle := &models.Vehicle{
		ID:                          otherVehicleID,
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

	// Mock VehicleRepository to return existing vehicle with different ID
	vehicleRepoMock.FindByNumberPlateFunc = func(numberPlate string) (*models.Vehicle, error) {
		return existingVehicle, nil
	}

	useCase := &UpdateByIdVehicle{
		VehicleRepository: vehicleRepoMock,
		Logger:            loggerMock,
	}

	// Act
	err := useCase.ValidateNumberPlateUniqueness(numberPlate, vehicleID)

	// Assert
	if err == nil {
		t.Error("Expected error for existing number plate, got nil")
	}

	if err.Error() != "number plate already exists" {
		t.Errorf("Expected error message 'number plate already exists', got '%s'", err.Error())
	}
}

func TestUpdateByIdVehicle_ValidateNumberPlateUniqueness_SameVehicle(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	numberPlate := "ABC1234"
	vehicleID := uuid.New()

	existingVehicle := &models.Vehicle{
		ID:                          vehicleID, // Same ID
		CustomerID:                  uuid.New(),
		Model:                       "Corolla",
		Brand:                       "Toyota",
		ReleaseYear:                 2020,
		VehicleIdentificationNumber: "VIN123456789",
		NumberPlate:                 numberPlate,
		Color:                       "Prata",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository to return existing vehicle with same ID
	vehicleRepoMock.FindByNumberPlateFunc = func(numberPlate string) (*models.Vehicle, error) {
		return existingVehicle, nil
	}

	useCase := &UpdateByIdVehicle{
		VehicleRepository: vehicleRepoMock,
		Logger:            loggerMock,
	}

	// Act
	err := useCase.ValidateNumberPlateUniqueness(numberPlate, vehicleID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for same vehicle, got %v", err)
	}
}

func TestUpdateByIdVehicle_ValidateCustomerExists_Success(t *testing.T) {
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

	useCase := &UpdateByIdVehicle{
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

func TestUpdateByIdVehicle_ValidateCustomerExists_NilCustomerID(t *testing.T) {
	// Arrange
	loggerMock := &mocks.LoggerMock{}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	useCase := &UpdateByIdVehicle{
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

func TestUpdateByIdVehicle_ValidateCustomerExists_CustomerNotFound(t *testing.T) {
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

	useCase := &UpdateByIdVehicle{
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

func TestUpdateByIdVehicle_UpdateVehicleFields(t *testing.T) {
	// Arrange
	loggerMock := &mocks.LoggerMock{}

	existingVehicle := &models.Vehicle{
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

	vehicleEntity := &domain.Vehicle{
		ID:                          existingVehicle.ID,
		CustomerID:                  uuid.New(),
		Model:                       "Corolla XSE",
		Brand:                       "Toyota",
		ReleaseYear:                 2021,
		VehicleIdentificationNumber: "VIN987654321",
		NumberPlate:                 "XYZ5678",
		Color:                       "Branco",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	useCase := &UpdateByIdVehicle{
		Logger: loggerMock,
	}

	// Act
	useCase.UpdateVehicleFields(existingVehicle, vehicleEntity)

	// Assert
	if existingVehicle.Model != "Corolla XSE" {
		t.Errorf("Expected model 'Corolla XSE', got '%s'", existingVehicle.Model)
	}

	if existingVehicle.ReleaseYear != 2021 {
		t.Errorf("Expected release year 2021, got %d", existingVehicle.ReleaseYear)
	}

	if existingVehicle.VehicleIdentificationNumber != "VIN987654321" {
		t.Errorf("Expected VIN 'VIN987654321', got '%s'", existingVehicle.VehicleIdentificationNumber)
	}

	if existingVehicle.NumberPlate != "XYZ5678" {
		t.Errorf("Expected number plate 'XYZ5678', got '%s'", existingVehicle.NumberPlate)
	}

	if existingVehicle.Color != "Branco" {
		t.Errorf("Expected color 'Branco', got '%s'", existingVehicle.Color)
	}
}

func TestUpdateByIdVehicle_SaveVehicleToDB_Success(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	vehicleModel := &models.Vehicle{
		ID:                          uuid.New(),
		CustomerID:                  uuid.New(),
		Model:                       "Corolla XSE",
		Brand:                       "Toyota",
		ReleaseYear:                 2021,
		VehicleIdentificationNumber: "VIN987654321",
		NumberPlate:                 "XYZ5678",
		Color:                       "Branco",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository
	vehicleRepoMock.UpdateFunc = func(vehicle *models.Vehicle) error {
		return nil
	}

	useCase := &UpdateByIdVehicle{
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

func TestUpdateByIdVehicle_SaveVehicleToDB_Error(t *testing.T) {
	// Arrange
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	vehicleModel := &models.Vehicle{
		ID:                          uuid.New(),
		CustomerID:                  uuid.New(),
		Model:                       "Corolla XSE",
		Brand:                       "Toyota",
		ReleaseYear:                 2021,
		VehicleIdentificationNumber: "VIN987654321",
		NumberPlate:                 "XYZ5678",
		Color:                       "Branco",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock VehicleRepository to return error
	vehicleRepoMock.UpdateFunc = func(vehicle *models.Vehicle) error {
		return errors.New("database error")
	}

	useCase := &UpdateByIdVehicle{
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
