package order

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/test/mocks"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/order_status_history"
	"go.uber.org/zap"
)

func TestUpdateOrderStatus_Process_Success(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}

	orderID := uuid.New()
	newStatus := "In progress"

	mockOrder := &models.Order{
		ID:         orderID,
		CustomerID: uuid.New(),
		VehicleID:  uuid.New(),
		Status:     "Received",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock OrderRepository
	orderRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Order, error) {
		return mockOrder, nil
	}

	orderRepoMock.UpdateFunc = func(order *models.Order) error {
		return nil
	}

	// Mock OrderStatusHistoryRepository
	orderStatusHistoryRepoMock.FindCurrentByOrderIDFunc = func(orderID uuid.UUID) (*models.OrderStatusHistory, error) {
		return &models.OrderStatusHistory{
			ID:        uuid.New(),
			OrderID:   orderID,
			Status:    "Received",
			StartedAt: time.Now().Add(-time.Hour),
		}, nil
	}

	orderStatusHistoryRepoMock.UpdateFunc = func(statusHistory *models.OrderStatusHistory) error {
		return nil
	}

	orderStatusHistoryRepoMock.CreateFunc = func(statusHistory *models.OrderStatusHistory) error {
		return nil
	}

	// Instantiate the real ManageOrderStatusHistory usecase with its mocked dependency
	statusHistoryManager := &order_status_history.ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	useCase := &UpdateOrderStatus{
		OrderRepository:      orderRepoMock,
		Logger:               loggerMock,
		StatusHistoryManager: statusHistoryManager,
	}

	// Act
	err := useCase.Process(orderID, newStatus)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestUpdateOrderStatus_Process_InvalidStatus(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}

	orderID := uuid.New()
	invalidStatus := "Invalid Status"

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Instantiate the real ManageOrderStatusHistory usecase with its mocked dependency
	statusHistoryManager := &order_status_history.ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	useCase := &UpdateOrderStatus{
		OrderRepository:      orderRepoMock,
		Logger:               loggerMock,
		StatusHistoryManager: statusHistoryManager,
	}

	// Act
	err := useCase.Process(orderID, invalidStatus)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid status, got nil")
	}

	if err.Error() != "invalid order status" {
		t.Errorf("Expected error message 'invalid order status', got '%s'", err.Error())
	}
}

func TestUpdateOrderStatus_Process_OrderNotFound(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}

	orderID := uuid.New()
	newStatus := "In progress"

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock OrderRepository to return error
	orderRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Order, error) {
		return nil, errors.New("order not found")
	}

	// Instantiate the real ManageOrderStatusHistory usecase with its mocked dependency
	statusHistoryManager := &order_status_history.ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	useCase := &UpdateOrderStatus{
		OrderRepository:      orderRepoMock,
		Logger:               loggerMock,
		StatusHistoryManager: statusHistoryManager,
	}

	// Act
	err := useCase.Process(orderID, newStatus)

	// Assert
	if err == nil {
		t.Error("Expected error for order not found, got nil")
	}

	if err.Error() != "order not found" {
		t.Errorf("Expected error message 'order not found', got '%s'", err.Error())
	}
}

func TestUpdateOrderStatus_Process_UpdateOrderError(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}

	orderID := uuid.New()
	newStatus := "In progress"

	mockOrder := &models.Order{
		ID:         orderID,
		CustomerID: uuid.New(),
		VehicleID:  uuid.New(),
		Status:     "Received",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock OrderRepository
	orderRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Order, error) {
		return mockOrder, nil
	}

	orderRepoMock.UpdateFunc = func(order *models.Order) error {
		return errors.New("database error")
	}

	// Instantiate the real ManageOrderStatusHistory usecase with its mocked dependency
	statusHistoryManager := &order_status_history.ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	useCase := &UpdateOrderStatus{
		OrderRepository:      orderRepoMock,
		Logger:               loggerMock,
		StatusHistoryManager: statusHistoryManager,
	}

	// Act
	err := useCase.Process(orderID, newStatus)

	// Assert
	if err == nil {
		t.Error("Expected error for database update failure, got nil")
	}

	if err.Error() != "database error" {
		t.Errorf("Expected error message 'database error', got '%s'", err.Error())
	}
}

func TestUpdateOrderStatus_GetValidStatuses(t *testing.T) {
	// Arrange
	useCase := &UpdateOrderStatus{}

	// Act
	validStatuses := useCase.GetValidStatuses()

	// Assert
	expectedStatuses := []string{
		"Received",
		"Undergoing diagnosis",
		"Awaiting approval",
		"In progress",
		"Completed",
		"Delivered",
		"Canceled",
	}

	if len(validStatuses) != len(expectedStatuses) {
		t.Errorf("Expected %d statuses, got %d", len(expectedStatuses), len(validStatuses))
	}

	for i, status := range expectedStatuses {
		if validStatuses[i] != status {
			t.Errorf("Expected status '%s' at position %d, got '%s'", status, i, validStatuses[i])
		}
	}
}

func TestUpdateOrderStatus_ValidateOrderStatus_Valid(t *testing.T) {
	// Arrange
	loggerMock := &mocks.LoggerMock{}
	useCase := &UpdateOrderStatus{
		Logger: loggerMock,
	}

	validStatuses := []string{
		"Received",
		"Undergoing diagnosis",
		"Awaiting approval",
		"In progress",
		"Completed",
		"Delivered",
		"Canceled",
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	for _, status := range validStatuses {
		t.Run("Status_"+status, func(t *testing.T) {
			// Act
			err := useCase.ValidateOrderStatus(status)

			// Assert
			if err != nil {
				t.Errorf("Expected no error for valid status '%s', got %v", status, err)
			}
		})
	}
}

func TestUpdateOrderStatus_ValidateOrderStatus_Invalid(t *testing.T) {
	// Arrange
	loggerMock := &mocks.LoggerMock{}
	useCase := &UpdateOrderStatus{
		Logger: loggerMock,
	}

	invalidStatuses := []string{
		"Invalid Status",
		"",
		"pending",
		"IN PROGRESS",
		"completed",
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	for _, status := range invalidStatuses {
		t.Run("Status_"+status, func(t *testing.T) {
			// Act
			err := useCase.ValidateOrderStatus(status)

			// Assert
			if err == nil {
				t.Errorf("Expected error for invalid status '%s', got nil", status)
			}

			if err.Error() != "invalid order status" {
				t.Errorf("Expected error message 'invalid order status', got '%s'", err.Error())
			}
		})
	}
}

func TestUpdateOrderStatus_FetchOrderFromDB_Success(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()
	mockOrder := &models.Order{
		ID:         orderID,
		CustomerID: uuid.New(),
		VehicleID:  uuid.New(),
		Status:     "Received",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock OrderRepository
	orderRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Order, error) {
		return mockOrder, nil
	}

	useCase := &UpdateOrderStatus{
		OrderRepository: orderRepoMock,
		Logger:          loggerMock,
	}

	// Act
	result, err := useCase.FetchOrderFromDB(orderID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected order, got nil")
	}

	if result.ID != orderID {
		t.Errorf("Expected order ID %s, got %s", orderID, result.ID)
	}
}

func TestUpdateOrderStatus_FetchOrderFromDB_NotFound(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock OrderRepository to return error
	orderRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Order, error) {
		return nil, errors.New("order not found")
	}

	useCase := &UpdateOrderStatus{
		OrderRepository: orderRepoMock,
		Logger:          loggerMock,
	}

	// Act
	result, err := useCase.FetchOrderFromDB(orderID)

	// Assert
	if err == nil {
		t.Error("Expected error for order not found, got nil")
	}

	if result != nil {
		t.Error("Expected nil result, got order")
	}

	if err.Error() != "order not found" {
		t.Errorf("Expected error message 'order not found', got '%s'", err.Error())
	}
}

func TestUpdateOrderStatus_UpdateOrderStatusInDB_Success(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()
	newStatus := "In progress"
	mockOrder := &models.Order{
		ID:         orderID,
		CustomerID: uuid.New(),
		VehicleID:  uuid.New(),
		Status:     "Received",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock OrderRepository
	orderRepoMock.UpdateFunc = func(order *models.Order) error {
		return nil
	}

	useCase := &UpdateOrderStatus{
		OrderRepository: orderRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.UpdateOrderStatusInDB(mockOrder, newStatus)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mockOrder.Status != newStatus {
		t.Errorf("Expected status to be updated to '%s', got '%s'", newStatus, mockOrder.Status)
	}
}

func TestUpdateOrderStatus_UpdateOrderStatusInDB_Error(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()
	newStatus := "In progress"
	mockOrder := &models.Order{
		ID:         orderID,
		CustomerID: uuid.New(),
		VehicleID:  uuid.New(),
		Status:     "Received",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock OrderRepository to return error
	orderRepoMock.UpdateFunc = func(order *models.Order) error {
		return errors.New("database error")
	}

	useCase := &UpdateOrderStatus{
		OrderRepository: orderRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.UpdateOrderStatusInDB(mockOrder, newStatus)

	// Assert
	if err == nil {
		t.Error("Expected error for database update failure, got nil")
	}

	if err.Error() != "database error" {
		t.Errorf("Expected error message 'database error', got '%s'", err.Error())
	}
}

func TestUpdateOrderStatus_UpdateStatusHistory_Success(t *testing.T) {
	// Arrange
	loggerMock := &mocks.LoggerMock{}
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}

	orderID := uuid.New()
	newStatus := "In progress"

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock OrderStatusHistoryRepository
	orderStatusHistoryRepoMock.FindCurrentByOrderIDFunc = func(orderID uuid.UUID) (*models.OrderStatusHistory, error) {
		return &models.OrderStatusHistory{
			ID:        uuid.New(),
			OrderID:   orderID,
			Status:    "Received",
			StartedAt: time.Now().Add(-time.Hour),
		}, nil
	}

	orderStatusHistoryRepoMock.UpdateFunc = func(statusHistory *models.OrderStatusHistory) error {
		return nil
	}

	orderStatusHistoryRepoMock.CreateFunc = func(statusHistory *models.OrderStatusHistory) error {
		return nil
	}

	// Instantiate the real ManageOrderStatusHistory usecase with its mocked dependency
	statusHistoryManager := &order_status_history.ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	useCase := &UpdateOrderStatus{
		Logger:               loggerMock,
		StatusHistoryManager: statusHistoryManager,
	}

	// Act
	err := useCase.UpdateStatusHistory(orderID, newStatus)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestUpdateOrderStatus_UpdateStatusHistory_Error(t *testing.T) {
	// Arrange
	loggerMock := &mocks.LoggerMock{}
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}

	orderID := uuid.New()
	newStatus := "In progress"

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock OrderStatusHistoryRepository to return error
	orderStatusHistoryRepoMock.FindCurrentByOrderIDFunc = func(orderID uuid.UUID) (*models.OrderStatusHistory, error) {
		return nil, errors.New("database error")
	}

	// Instantiate the real ManageOrderStatusHistory usecase with its mocked dependency
	statusHistoryManager := &order_status_history.ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	useCase := &UpdateOrderStatus{
		Logger:               loggerMock,
		StatusHistoryManager: statusHistoryManager,
	}

	// Act
	err := useCase.UpdateStatusHistory(orderID, newStatus)

	// Assert
	// Should not return error even if status history update fails
	if err != nil {
		t.Errorf("Expected no error (status history errors are logged but not returned), got %v", err)
	}
}
