package order_status_history

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/test/mocks"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestManageOrderStatusHistory_IsFinalStatus(t *testing.T) {
	// Arrange
	useCase := &ManageOrderStatusHistory{}

	testCases := []struct {
		name     string
		status   string
		expected bool
	}{
		{"Delivered status", "Delivered", true},
		{"Canceled status", "Canceled", true},
		{"Received status", "Received", false},
		{"In progress status", "In progress", false},
		{"Completed status", "Completed", false},
		{"Empty status", "", false},
		{"Random status", "Random", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := useCase.IsFinalStatus(tc.status)

			// Assert
			if result != tc.expected {
				t.Errorf("Expected %v for status '%s', got %v", tc.expected, tc.status, result)
			}
		})
	}
}

func TestManageOrderStatusHistory_FetchCurrentStatusFromDB_Success(t *testing.T) {
	// Arrange
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()
	mockCurrentStatus := &models.OrderStatusHistory{
		ID:        uuid.New(),
		OrderID:   orderID,
		Status:    "In progress",
		StartedAt: time.Now().Add(-time.Hour),
		EndedAt:   nil,
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock repository
	orderStatusHistoryRepoMock.FindCurrentByOrderIDFunc = func(orderID uuid.UUID) (*models.OrderStatusHistory, error) {
		return mockCurrentStatus, nil
	}

	useCase := &ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	result, err := useCase.FetchCurrentStatusFromDB(orderID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected current status, got nil")
	}

	if result.ID != mockCurrentStatus.ID {
		t.Errorf("Expected status ID %s, got %s", mockCurrentStatus.ID, result.ID)
	}

	if result.Status != "In progress" {
		t.Errorf("Expected status 'In progress', got '%s'", result.Status)
	}
}

func TestManageOrderStatusHistory_FetchCurrentStatusFromDB_NotFound(t *testing.T) {
	// Arrange
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock repository to return gorm.ErrRecordNotFound
	orderStatusHistoryRepoMock.FindCurrentByOrderIDFunc = func(orderID uuid.UUID) (*models.OrderStatusHistory, error) {
		return nil, gorm.ErrRecordNotFound
	}

	useCase := &ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	result, err := useCase.FetchCurrentStatusFromDB(orderID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for gorm.ErrRecordNotFound, got %v", err)
	}

	if result != nil {
		t.Error("Expected nil result for not found, got status")
	}
}

func TestManageOrderStatusHistory_FetchCurrentStatusFromDB_Error(t *testing.T) {
	// Arrange
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock repository to return error
	orderStatusHistoryRepoMock.FindCurrentByOrderIDFunc = func(orderID uuid.UUID) (*models.OrderStatusHistory, error) {
		return nil, errors.New("database error")
	}

	useCase := &ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	result, err := useCase.FetchCurrentStatusFromDB(orderID)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result for error, got status")
	}

	if err.Error() != "database error" {
		t.Errorf("Expected error message 'database error', got '%s'", err.Error())
	}
}

func TestManageOrderStatusHistory_FinalizeCurrentStatus_Success(t *testing.T) {
	// Arrange
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()
	mockCurrentStatus := &models.OrderStatusHistory{
		ID:        uuid.New(),
		OrderID:   orderID,
		Status:    "In progress",
		StartedAt: time.Now().Add(-time.Hour),
		EndedAt:   nil,
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock repository
	orderStatusHistoryRepoMock.FindCurrentByOrderIDFunc = func(orderID uuid.UUID) (*models.OrderStatusHistory, error) {
		return mockCurrentStatus, nil
	}

	orderStatusHistoryRepoMock.UpdateFunc = func(statusHistory *models.OrderStatusHistory) error {
		return nil
	}

	useCase := &ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	err := useCase.FinalizeCurrentStatus(orderID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestManageOrderStatusHistory_FinalizeCurrentStatus_NoCurrentStatus(t *testing.T) {
	// Arrange
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock repository to return nil (no current status)
	orderStatusHistoryRepoMock.FindCurrentByOrderIDFunc = func(orderID uuid.UUID) (*models.OrderStatusHistory, error) {
		return nil, nil
	}

	useCase := &ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	err := useCase.FinalizeCurrentStatus(orderID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for no current status, got %v", err)
	}
}

func TestManageOrderStatusHistory_FinalizeCurrentStatus_UpdateError(t *testing.T) {
	// Arrange
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()
	mockCurrentStatus := &models.OrderStatusHistory{
		ID:        uuid.New(),
		OrderID:   orderID,
		Status:    "In progress",
		StartedAt: time.Now().Add(-time.Hour),
		EndedAt:   nil,
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock repository
	orderStatusHistoryRepoMock.FindCurrentByOrderIDFunc = func(orderID uuid.UUID) (*models.OrderStatusHistory, error) {
		return mockCurrentStatus, nil
	}

	orderStatusHistoryRepoMock.UpdateFunc = func(statusHistory *models.OrderStatusHistory) error {
		return errors.New("update error")
	}

	useCase := &ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	err := useCase.FinalizeCurrentStatus(orderID)

	// Assert
	if err == nil {
		t.Error("Expected error for update failure, got nil")
	}

	if err.Error() != "update error" {
		t.Errorf("Expected error message 'update error', got '%s'", err.Error())
	}
}

func TestManageOrderStatusHistory_CreateNewStatus_Success(t *testing.T) {
	// Arrange
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()
	status := "In progress"

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock repository
	orderStatusHistoryRepoMock.CreateFunc = func(statusHistory *models.OrderStatusHistory) error {
		return nil
	}

	useCase := &ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	err := useCase.CreateNewStatus(orderID, status)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestManageOrderStatusHistory_CreateNewStatus_Error(t *testing.T) {
	// Arrange
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()
	status := "In progress"

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock repository to return error
	orderStatusHistoryRepoMock.CreateFunc = func(statusHistory *models.OrderStatusHistory) error {
		return errors.New("create error")
	}

	useCase := &ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	err := useCase.CreateNewStatus(orderID, status)

	// Assert
	if err == nil {
		t.Error("Expected error for create failure, got nil")
	}

	if err.Error() != "create error" {
		t.Errorf("Expected error message 'create error', got '%s'", err.Error())
	}
}

func TestManageOrderStatusHistory_UpdateCurrentStatusToFinal_Success(t *testing.T) {
	// Arrange
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()
	finalStatus := "Delivered"
	mockCurrentStatus := &models.OrderStatusHistory{
		ID:        uuid.New(),
		OrderID:   orderID,
		Status:    "In progress",
		StartedAt: time.Now().Add(-time.Hour),
		EndedAt:   nil,
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock repository
	orderStatusHistoryRepoMock.FindCurrentByOrderIDFunc = func(orderID uuid.UUID) (*models.OrderStatusHistory, error) {
		return mockCurrentStatus, nil
	}

	orderStatusHistoryRepoMock.UpdateFunc = func(statusHistory *models.OrderStatusHistory) error {
		return nil
	}

	useCase := &ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	err := useCase.UpdateCurrentStatusToFinal(orderID, finalStatus)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestManageOrderStatusHistory_UpdateCurrentStatusToFinal_NoCurrentStatus(t *testing.T) {
	// Arrange
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()
	finalStatus := "Delivered"

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock repository to return nil (no current status)
	orderStatusHistoryRepoMock.FindCurrentByOrderIDFunc = func(orderID uuid.UUID) (*models.OrderStatusHistory, error) {
		return nil, nil
	}

	useCase := &ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	err := useCase.UpdateCurrentStatusToFinal(orderID, finalStatus)

	// Assert
	if err == nil {
		t.Error("Expected error for no current status, got nil")
	}

	if err.Error() != "no current status found" {
		t.Errorf("Expected error message 'no current status found', got '%s'", err.Error())
	}
}

func TestManageOrderStatusHistory_UpdateCurrentStatusToFinal_UpdateError(t *testing.T) {
	// Arrange
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()
	finalStatus := "Delivered"
	mockCurrentStatus := &models.OrderStatusHistory{
		ID:        uuid.New(),
		OrderID:   orderID,
		Status:    "In progress",
		StartedAt: time.Now().Add(-time.Hour),
		EndedAt:   nil,
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock repository
	orderStatusHistoryRepoMock.FindCurrentByOrderIDFunc = func(orderID uuid.UUID) (*models.OrderStatusHistory, error) {
		return mockCurrentStatus, nil
	}

	orderStatusHistoryRepoMock.UpdateFunc = func(statusHistory *models.OrderStatusHistory) error {
		return errors.New("update error")
	}

	useCase := &ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	err := useCase.UpdateCurrentStatusToFinal(orderID, finalStatus)

	// Assert
	if err == nil {
		t.Error("Expected error for update failure, got nil")
	}

	if err.Error() != "update error" {
		t.Errorf("Expected error message 'update error', got '%s'", err.Error())
	}
}

func TestManageOrderStatusHistory_FetchOrderHistoryFromDB_Success(t *testing.T) {
	// Arrange
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()
	mockHistory := []models.OrderStatusHistory{
		{
			ID:        uuid.New(),
			OrderID:   orderID,
			Status:    "Received",
			StartedAt: time.Now().Add(-2 * time.Hour),
			EndedAt:   &[]time.Time{time.Now().Add(-time.Hour)}[0],
		},
		{
			ID:        uuid.New(),
			OrderID:   orderID,
			Status:    "In progress",
			StartedAt: time.Now().Add(-time.Hour),
			EndedAt:   nil,
		},
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock repository
	orderStatusHistoryRepoMock.FindByOrderIDFunc = func(orderID uuid.UUID) ([]models.OrderStatusHistory, error) {
		return mockHistory, nil
	}

	useCase := &ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	result, err := useCase.FetchOrderHistoryFromDB(orderID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 history items, got %d", len(result))
	}

	if result[0].Status != "Received" {
		t.Errorf("Expected first status 'Received', got '%s'", result[0].Status)
	}

	if result[1].Status != "In progress" {
		t.Errorf("Expected second status 'In progress', got '%s'", result[1].Status)
	}
}

func TestManageOrderStatusHistory_FetchOrderHistoryFromDB_Error(t *testing.T) {
	// Arrange
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock repository to return error
	orderStatusHistoryRepoMock.FindByOrderIDFunc = func(orderID uuid.UUID) ([]models.OrderStatusHistory, error) {
		return nil, errors.New("database error")
	}

	useCase := &ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	result, err := useCase.FetchOrderHistoryFromDB(orderID)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result for error, got history")
	}

	if err.Error() != "database error" {
		t.Errorf("Expected error message 'database error', got '%s'", err.Error())
	}
}

func TestManageOrderStatusHistory_StartNewStatus(t *testing.T) {
	// Arrange
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()
	status := "In progress"

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock repository
	orderStatusHistoryRepoMock.CreateFunc = func(statusHistory *models.OrderStatusHistory) error {
		return nil
	}

	useCase := &ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	err := useCase.StartNewStatus(orderID, status)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestManageOrderStatusHistory_UpdateStatus_NonFinalStatus_Success(t *testing.T) {
	// Arrange
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()
	newStatus := "In progress"
	mockCurrentStatus := &models.OrderStatusHistory{
		ID:        uuid.New(),
		OrderID:   orderID,
		Status:    "Received",
		StartedAt: time.Now().Add(-time.Hour),
		EndedAt:   nil,
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock repository
	orderStatusHistoryRepoMock.FindCurrentByOrderIDFunc = func(orderID uuid.UUID) (*models.OrderStatusHistory, error) {
		return mockCurrentStatus, nil
	}

	orderStatusHistoryRepoMock.UpdateFunc = func(statusHistory *models.OrderStatusHistory) error {
		return nil
	}

	orderStatusHistoryRepoMock.CreateFunc = func(statusHistory *models.OrderStatusHistory) error {
		return nil
	}

	useCase := &ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	err := useCase.UpdateStatus(orderID, newStatus)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestManageOrderStatusHistory_UpdateStatus_FinalStatus_Success(t *testing.T) {
	// Arrange
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()
	newStatus := "Delivered"
	mockCurrentStatus := &models.OrderStatusHistory{
		ID:        uuid.New(),
		OrderID:   orderID,
		Status:    "In progress",
		StartedAt: time.Now().Add(-time.Hour),
		EndedAt:   nil,
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock repository
	orderStatusHistoryRepoMock.FindCurrentByOrderIDFunc = func(orderID uuid.UUID) (*models.OrderStatusHistory, error) {
		return mockCurrentStatus, nil
	}

	orderStatusHistoryRepoMock.UpdateFunc = func(statusHistory *models.OrderStatusHistory) error {
		return nil
	}

	useCase := &ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	err := useCase.UpdateStatus(orderID, newStatus)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestManageOrderStatusHistory_UpdateStatus_FinalizeError(t *testing.T) {
	// Arrange
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()
	newStatus := "In progress"

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock repository to return error on finalize
	orderStatusHistoryRepoMock.FindCurrentByOrderIDFunc = func(orderID uuid.UUID) (*models.OrderStatusHistory, error) {
		return &models.OrderStatusHistory{
			ID:        uuid.New(),
			OrderID:   orderID,
			Status:    "Received",
			StartedAt: time.Now().Add(-time.Hour),
			EndedAt:   nil,
		}, nil
	}

	orderStatusHistoryRepoMock.UpdateFunc = func(statusHistory *models.OrderStatusHistory) error {
		return errors.New("finalize error")
	}

	useCase := &ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	err := useCase.UpdateStatus(orderID, newStatus)

	// Assert
	if err == nil {
		t.Error("Expected error for finalize failure, got nil")
	}

	if err.Error() != "finalize error" {
		t.Errorf("Expected error message 'finalize error', got '%s'", err.Error())
	}
}

func TestManageOrderStatusHistory_UpdateStatus_CreateNewError(t *testing.T) {
	// Arrange
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()
	newStatus := "In progress"

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock repository
	orderStatusHistoryRepoMock.FindCurrentByOrderIDFunc = func(orderID uuid.UUID) (*models.OrderStatusHistory, error) {
		return &models.OrderStatusHistory{
			ID:        uuid.New(),
			OrderID:   orderID,
			Status:    "Received",
			StartedAt: time.Now().Add(-time.Hour),
			EndedAt:   nil,
		}, nil
	}

	orderStatusHistoryRepoMock.UpdateFunc = func(statusHistory *models.OrderStatusHistory) error {
		return nil
	}

	orderStatusHistoryRepoMock.CreateFunc = func(statusHistory *models.OrderStatusHistory) error {
		return errors.New("create error")
	}

	useCase := &ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	err := useCase.UpdateStatus(orderID, newStatus)

	// Assert
	if err == nil {
		t.Error("Expected error for create failure, got nil")
	}

	if err.Error() != "create error" {
		t.Errorf("Expected error message 'create error', got '%s'", err.Error())
	}
}

func TestManageOrderStatusHistory_GetOrderHistory(t *testing.T) {
	// Arrange
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	orderID := uuid.New()
	mockHistory := []models.OrderStatusHistory{
		{
			ID:        uuid.New(),
			OrderID:   orderID,
			Status:    "Received",
			StartedAt: time.Now().Add(-2 * time.Hour),
			EndedAt:   &[]time.Time{time.Now().Add(-time.Hour)}[0],
		},
	}

	// Mock logger
	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {}
	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {}

	// Mock repository
	orderStatusHistoryRepoMock.FindByOrderIDFunc = func(orderID uuid.UUID) ([]models.OrderStatusHistory, error) {
		return mockHistory, nil
	}

	useCase := &ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	result, err := useCase.GetOrderHistory(orderID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 history item, got %d", len(result))
	}

	if result[0].Status != "Received" {
		t.Errorf("Expected status 'Received', got '%s'", result[0].Status)
	}
}
