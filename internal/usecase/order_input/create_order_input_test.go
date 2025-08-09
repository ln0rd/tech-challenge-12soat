package order_input

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/order_input"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/test/mocks"
	"go.uber.org/zap"
)

func TestCreateOrderInput_Process_Success(t *testing.T) {
	// Arrange
	orderInputRepoMock := &mocks.OrderInputRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	orderInputRepoMock.CreateFunc = func(orderInput *models.OrderInput) error {
		// Simula sucesso na criação
		return nil
	}

	useCase := &CreateOrderInput{
		OrderInputRepository: orderInputRepoMock,
		Logger:               loggerMock,
	}

	orderInput := &domain.OrderInput{
		ID:        uuid.New(),
		OrderID:   uuid.New(),
		InputID:   uuid.New(),
		Quantity:  5,
		UnitPrice: 10.50,
		TotalPrice: 52.50,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	err := useCase.Process(orderInput)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verifica se os logs corretos foram chamados
	if len(loggedInfo) < 3 {
		t.Errorf("Expected at least 3 info logs, got %d", len(loggedInfo))
	}

	expectedInfoLogs := []string{
		"Processing order input creation",
		"Model created",
		"Order input created in database",
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

func TestCreateOrderInput_Process_DatabaseError(t *testing.T) {
	// Arrange
	orderInputRepoMock := &mocks.OrderInputRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	expectedError := errors.New("database connection failed")
	orderInputRepoMock.CreateFunc = func(orderInput *models.OrderInput) error {
		return expectedError
	}

	useCase := &CreateOrderInput{
		OrderInputRepository: orderInputRepoMock,
		Logger:               loggerMock,
	}

	orderInput := &domain.OrderInput{
		ID:        uuid.New(),
		OrderID:   uuid.New(),
		InputID:   uuid.New(),
		Quantity:  3,
		UnitPrice: 15.00,
		TotalPrice: 45.00,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	err := useCase.Process(orderInput)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	// Verifica se o log de erro foi chamado
	if len(loggedErrors) == 0 {
		t.Error("Expected error log, got none")
	}

	expectedErrorLog := "Database error creating order input"
	found := false
	for _, actualLog := range loggedErrors {
		if actualLog == expectedErrorLog {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected error log message '%s' not found", expectedErrorLog)
	}
}

func TestCreateOrderInput_SaveOrderInputToDB_Success(t *testing.T) {
	// Arrange
	orderInputRepoMock := &mocks.OrderInputRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	orderInputRepoMock.CreateFunc = func(orderInput *models.OrderInput) error {
		// Simula sucesso na criação
		return nil
	}

	useCase := &CreateOrderInput{
		OrderInputRepository: orderInputRepoMock,
		Logger:               loggerMock,
	}

	orderInput := &models.OrderInput{
		ID:         uuid.New(),
		OrderID:    uuid.New(),
		InputID:    uuid.New(),
		Quantity:   2,
		UnitPrice:  25.00,
		TotalPrice: 50.00,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Act
	err := useCase.SaveOrderInputToDB(orderInput)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verifica se o log de sucesso foi chamado
	if len(loggedInfo) == 0 {
		t.Error("Expected info log, got none")
	}

	expectedInfoLog := "Order input created in database"
	found := false
	for _, actualLog := range loggedInfo {
		if actualLog == expectedInfoLog {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected info log message '%s' not found", expectedInfoLog)
	}

	if len(loggedErrors) > 0 {
		t.Errorf("Expected no error logs, got %d", len(loggedErrors))
	}
}

func TestCreateOrderInput_SaveOrderInputToDB_Error(t *testing.T) {
	// Arrange
	orderInputRepoMock := &mocks.OrderInputRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	expectedError := errors.New("database constraint violation")
	orderInputRepoMock.CreateFunc = func(orderInput *models.OrderInput) error {
		return expectedError
	}

	useCase := &CreateOrderInput{
		OrderInputRepository: orderInputRepoMock,
		Logger:               loggerMock,
	}

	orderInput := &models.OrderInput{
		ID:         uuid.New(),
		OrderID:    uuid.New(),
		InputID:    uuid.New(),
		Quantity:   1,
		UnitPrice:  30.00,
		TotalPrice: 30.00,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Act
	err := useCase.SaveOrderInputToDB(orderInput)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	// Verifica se o log de erro foi chamado
	if len(loggedErrors) == 0 {
		t.Error("Expected error log, got none")
	}

	expectedErrorLog := "Database error creating order input"
	found := false
	for _, actualLog := range loggedErrors {
		if actualLog == expectedErrorLog {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected error log message '%s' not found", expectedErrorLog)
	}
}
