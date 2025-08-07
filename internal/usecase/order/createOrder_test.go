package order

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/order"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/test/mocks"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/order_status_history"
	"go.uber.org/zap"
)

func TestCreateOrder_Process_Success(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	customerID := uuid.New()
	vehicleID := uuid.New()

	mockCustomer := &models.Customer{
		ID:             customerID,
		Name:           "João Silva",
		DocumentNumber: "12345678901",
		CustomerType:   "individual",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	mockVehicle := &models.Vehicle{
		ID:                          vehicleID,
		CustomerID:                  customerID,
		NumberPlate:                 "ABC1234",
		Brand:                       "Toyota",
		Model:                       "Corolla",
		ReleaseYear:                 2020,
		VehicleIdentificationNumber: "VIN123456789",
		Color:                       "Prata",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		if id == customerID {
			return mockCustomer, nil
		}
		return nil, errors.New("customer not found")
	}

	vehicleRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Vehicle, error) {
		if id == vehicleID {
			return mockVehicle, nil
		}
		return nil, errors.New("vehicle not found")
	}

	orderRepoMock.CreateFunc = func(order *models.Order) error {
		return nil
	}

	orderStatusHistoryRepoMock.CreateFunc = func(statusHistory *models.OrderStatusHistory) error {
		return nil
	}

	statusHistoryManager := &order_status_history.ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	useCase := &CreateOrder{
		OrderRepository:      orderRepoMock,
		CustomerRepository:   customerRepoMock,
		VehicleRepository:    vehicleRepoMock,
		Logger:               loggerMock,
		StatusHistoryManager: statusHistoryManager,
	}

	order := &domain.Order{
		ID:         uuid.New(),
		CustomerID: customerID,
		VehicleID:  vehicleID,
		Status:     "Received",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Act
	err := useCase.Process(order)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing order creation",
		"Customer found",
		"Vehicle found",
		"Vehicle belongs to customer",
		"Model created",
		"Order created in database",
		"New status started",
		"Status history started successfully",
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

func TestCreateOrder_Process_CustomerNotFound(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	customerID := uuid.New()
	vehicleID := uuid.New()

	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		return nil, errors.New("customer not found")
	}

	statusHistoryManager := &order_status_history.ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	useCase := &CreateOrder{
		OrderRepository:      orderRepoMock,
		CustomerRepository:   customerRepoMock,
		VehicleRepository:    vehicleRepoMock,
		Logger:               loggerMock,
		StatusHistoryManager: statusHistoryManager,
	}

	order := &domain.Order{
		ID:         uuid.New(),
		CustomerID: customerID,
		VehicleID:  vehicleID,
		Status:     "Received",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Act
	err := useCase.Process(order)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "customer not found" {
		t.Errorf("Expected error 'customer not found', got '%s'", err.Error())
	}

	// Verifica se os logs corretos foram chamados
	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Processing order creation" {
		t.Errorf("Expected log message 'Processing order creation', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Customer not found" {
		t.Errorf("Expected error log 'Customer not found', got '%s'", loggedErrors[0])
	}
}

func TestCreateOrder_Process_VehicleNotFound(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	customerID := uuid.New()
	vehicleID := uuid.New()

	mockCustomer := &models.Customer{
		ID:             customerID,
		Name:           "João Silva",
		DocumentNumber: "12345678901",
		CustomerType:   "individual",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		if id == customerID {
			return mockCustomer, nil
		}
		return nil, errors.New("customer not found")
	}

	vehicleRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Vehicle, error) {
		return nil, errors.New("vehicle not found")
	}

	statusHistoryManager := &order_status_history.ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	useCase := &CreateOrder{
		OrderRepository:      orderRepoMock,
		CustomerRepository:   customerRepoMock,
		VehicleRepository:    vehicleRepoMock,
		Logger:               loggerMock,
		StatusHistoryManager: statusHistoryManager,
	}

	order := &domain.Order{
		ID:         uuid.New(),
		CustomerID: customerID,
		VehicleID:  vehicleID,
		Status:     "Received",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Act
	err := useCase.Process(order)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "vehicle not found" {
		t.Errorf("Expected error 'vehicle not found', got '%s'", err.Error())
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing order creation",
		"Customer found",
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

	if loggedErrors[0] != "Vehicle not found" {
		t.Errorf("Expected error log 'Vehicle not found', got '%s'", loggedErrors[0])
	}
}

func TestCreateOrder_Process_VehicleDoesNotBelongToCustomer(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	customerID := uuid.New()
	vehicleID := uuid.New()
	differentCustomerID := uuid.New()

	mockCustomer := &models.Customer{
		ID:             customerID,
		Name:           "João Silva",
		DocumentNumber: "12345678901",
		CustomerType:   "individual",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	mockVehicle := &models.Vehicle{
		ID:                          vehicleID,
		CustomerID:                  differentCustomerID, // Vehicle pertence a outro customer
		NumberPlate:                 "ABC1234",
		Brand:                       "Toyota",
		Model:                       "Corolla",
		ReleaseYear:                 2020,
		VehicleIdentificationNumber: "VIN123456789",
		Color:                       "Prata",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		if id == customerID {
			return mockCustomer, nil
		}
		return nil, errors.New("customer not found")
	}

	vehicleRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Vehicle, error) {
		if id == vehicleID {
			return mockVehicle, nil
		}
		return nil, errors.New("vehicle not found")
	}

	statusHistoryManager := &order_status_history.ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	useCase := &CreateOrder{
		OrderRepository:      orderRepoMock,
		CustomerRepository:   customerRepoMock,
		VehicleRepository:    vehicleRepoMock,
		Logger:               loggerMock,
		StatusHistoryManager: statusHistoryManager,
	}

	order := &domain.Order{
		ID:         uuid.New(),
		CustomerID: customerID,
		VehicleID:  vehicleID,
		Status:     "Received",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Act
	err := useCase.Process(order)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "vehicle does not belong to customer" {
		t.Errorf("Expected error 'vehicle does not belong to customer', got '%s'", err.Error())
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing order creation",
		"Customer found",
		"Vehicle found",
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

	if loggedErrors[0] != "Vehicle does not belong to customer" {
		t.Errorf("Expected error log 'Vehicle does not belong to customer', got '%s'", loggedErrors[0])
	}
}

func TestCreateOrder_Process_DatabaseError(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	customerID := uuid.New()
	vehicleID := uuid.New()

	mockCustomer := &models.Customer{
		ID:             customerID,
		Name:           "João Silva",
		DocumentNumber: "12345678901",
		CustomerType:   "individual",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	mockVehicle := &models.Vehicle{
		ID:                          vehicleID,
		CustomerID:                  customerID,
		NumberPlate:                 "ABC1234",
		Brand:                       "Toyota",
		Model:                       "Corolla",
		ReleaseYear:                 2020,
		VehicleIdentificationNumber: "VIN123456789",
		Color:                       "Prata",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		if id == customerID {
			return mockCustomer, nil
		}
		return nil, errors.New("customer not found")
	}

	vehicleRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Vehicle, error) {
		if id == vehicleID {
			return mockVehicle, nil
		}
		return nil, errors.New("vehicle not found")
	}

	expectedError := errors.New("database connection failed")
	orderRepoMock.CreateFunc = func(order *models.Order) error {
		return expectedError
	}

	statusHistoryManager := &order_status_history.ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	useCase := &CreateOrder{
		OrderRepository:      orderRepoMock,
		CustomerRepository:   customerRepoMock,
		VehicleRepository:    vehicleRepoMock,
		Logger:               loggerMock,
		StatusHistoryManager: statusHistoryManager,
	}

	order := &domain.Order{
		ID:         uuid.New(),
		CustomerID: customerID,
		VehicleID:  vehicleID,
		Status:     "Received",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Act
	err := useCase.Process(order)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing order creation",
		"Customer found",
		"Vehicle found",
		"Vehicle belongs to customer",
		"Model created",
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

	if loggedErrors[0] != "Database error creating order" {
		t.Errorf("Expected error log 'Database error creating order', got '%s'", loggedErrors[0])
	}
}

func TestCreateOrder_FetchCustomerFromDB_Success(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	customerID := uuid.New()
	mockCustomer := &models.Customer{
		ID:             customerID,
		Name:           "João Silva",
		DocumentNumber: "12345678901",
		CustomerType:   "individual",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		if id == customerID {
			return mockCustomer, nil
		}
		return nil, errors.New("customer not found")
	}

	statusHistoryManager := &order_status_history.ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	useCase := &CreateOrder{
		OrderRepository:      orderRepoMock,
		CustomerRepository:   customerRepoMock,
		VehicleRepository:    vehicleRepoMock,
		Logger:               loggerMock,
		StatusHistoryManager: statusHistoryManager,
	}

	// Act
	result, err := useCase.FetchCustomerFromDB(customerID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected result, got nil")
	}

	if result.ID != customerID {
		t.Errorf("Expected customer ID %s, got %s", customerID, result.ID)
	}

	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Customer found" {
		t.Errorf("Expected log message 'Customer found', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) > 0 {
		t.Errorf("Expected no error logs, got %d", len(loggedErrors))
	}
}

func TestCreateOrder_FetchCustomerFromDB_Error(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	customerID := uuid.New()
	expectedError := errors.New("database timeout")

	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		return nil, expectedError
	}

	statusHistoryManager := &order_status_history.ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	useCase := &CreateOrder{
		OrderRepository:      orderRepoMock,
		CustomerRepository:   customerRepoMock,
		VehicleRepository:    vehicleRepoMock,
		Logger:               loggerMock,
		StatusHistoryManager: statusHistoryManager,
	}

	// Act
	result, err := useCase.FetchCustomerFromDB(customerID)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "customer not found" {
		t.Errorf("Expected error 'customer not found', got '%s'", err.Error())
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	if len(loggedInfo) > 0 {
		t.Errorf("Expected no info logs, got %d", len(loggedInfo))
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Customer not found" {
		t.Errorf("Expected error log 'Customer not found', got '%s'", loggedErrors[0])
	}
}

func TestCreateOrder_ValidateVehicleOwnership_Success(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	customerID := uuid.New()
	vehicleID := uuid.New()

	mockVehicle := &models.Vehicle{
		ID:                          vehicleID,
		CustomerID:                  customerID,
		NumberPlate:                 "ABC1234",
		Brand:                       "Toyota",
		Model:                       "Corolla",
		ReleaseYear:                 2020,
		VehicleIdentificationNumber: "VIN123456789",
		Color:                       "Prata",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	statusHistoryManager := &order_status_history.ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	useCase := &CreateOrder{
		OrderRepository:      orderRepoMock,
		CustomerRepository:   customerRepoMock,
		VehicleRepository:    vehicleRepoMock,
		Logger:               loggerMock,
		StatusHistoryManager: statusHistoryManager,
	}

	// Act
	err := useCase.ValidateVehicleOwnership(mockVehicle, customerID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Vehicle belongs to customer" {
		t.Errorf("Expected log message 'Vehicle belongs to customer', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) > 0 {
		t.Errorf("Expected no error logs, got %d", len(loggedErrors))
	}
}

func TestCreateOrder_ValidateVehicleOwnership_Error(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	customerID := uuid.New()
	differentCustomerID := uuid.New()
	vehicleID := uuid.New()

	mockVehicle := &models.Vehicle{
		ID:                          vehicleID,
		CustomerID:                  differentCustomerID, // Vehicle pertence a outro customer
		NumberPlate:                 "ABC1234",
		Brand:                       "Toyota",
		Model:                       "Corolla",
		ReleaseYear:                 2020,
		VehicleIdentificationNumber: "VIN123456789",
		Color:                       "Prata",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	statusHistoryManager := &order_status_history.ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	useCase := &CreateOrder{
		OrderRepository:      orderRepoMock,
		CustomerRepository:   customerRepoMock,
		VehicleRepository:    vehicleRepoMock,
		Logger:               loggerMock,
		StatusHistoryManager: statusHistoryManager,
	}

	// Act
	err := useCase.ValidateVehicleOwnership(mockVehicle, customerID)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "vehicle does not belong to customer" {
		t.Errorf("Expected error 'vehicle does not belong to customer', got '%s'", err.Error())
	}

	if len(loggedInfo) > 0 {
		t.Errorf("Expected no info logs, got %d", len(loggedInfo))
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Vehicle does not belong to customer" {
		t.Errorf("Expected error log 'Vehicle does not belong to customer', got '%s'", loggedErrors[0])
	}
}
