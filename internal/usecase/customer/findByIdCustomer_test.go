package customer

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

func TestFindByIdCustomer_Process_Success(t *testing.T) {
	// Arrange
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

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
		return nil, gorm.ErrRecordNotFound
	}

	useCase := &FindByIdCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	result, err := useCase.Process(customerID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected customer result, got nil")
	}

	if result.ID != customerID {
		t.Errorf("Expected customer ID %s, got %s", customerID, result.ID)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing find customer by ID",
		"Successfully fetched customer from database",
		"Successfully mapped customer to domain",
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

func TestFindByIdCustomer_Process_CustomerNotFound(t *testing.T) {
	// Arrange
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	customerID := uuid.New()

	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		return nil, gorm.ErrRecordNotFound
	}

	useCase := &FindByIdCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	result, err := useCase.Process(customerID)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err != gorm.ErrRecordNotFound {
		t.Errorf("Expected error %v, got %v", gorm.ErrRecordNotFound, err)
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	// Verifica se os logs corretos foram chamados
	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Processing find customer by ID" {
		t.Errorf("Expected log message 'Processing find customer by ID', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Customer not found" {
		t.Errorf("Expected error log 'Customer not found', got '%s'", loggedErrors[0])
	}
}

func TestFindByIdCustomer_Process_DatabaseError(t *testing.T) {
	// Arrange
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	customerID := uuid.New()
	expectedError := errors.New("database connection failed")

	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		return nil, expectedError
	}

	useCase := &FindByIdCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	result, err := useCase.Process(customerID)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	// Verifica se os logs corretos foram chamados
	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Processing find customer by ID" {
		t.Errorf("Expected log message 'Processing find customer by ID', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Database error fetching customer" {
		t.Errorf("Expected error log 'Database error fetching customer', got '%s'", loggedErrors[0])
	}
}

func TestFindByIdCustomer_FetchCustomerFromDB_Success(t *testing.T) {
	// Arrange
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

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
		return nil, gorm.ErrRecordNotFound
	}

	useCase := &FindByIdCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	result, err := useCase.FetchCustomerFromDB(customerID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected customer result, got nil")
	}

	if result.ID != customerID {
		t.Errorf("Expected customer ID %s, got %s", customerID, result.ID)
	}

	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Successfully fetched customer from database" {
		t.Errorf("Expected log message 'Successfully fetched customer from database', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) > 0 {
		t.Errorf("Expected no error logs, got %d", len(loggedErrors))
	}
}

func TestFindByIdCustomer_FetchCustomerFromDB_NotFound(t *testing.T) {
	// Arrange
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	customerID := uuid.New()

	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		return nil, gorm.ErrRecordNotFound
	}

	useCase := &FindByIdCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	result, err := useCase.FetchCustomerFromDB(customerID)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err != gorm.ErrRecordNotFound {
		t.Errorf("Expected error %v, got %v", gorm.ErrRecordNotFound, err)
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

func TestFindByIdCustomer_FetchCustomerFromDB_DatabaseError(t *testing.T) {
	// Arrange
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

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

	useCase := &FindByIdCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	result, err := useCase.FetchCustomerFromDB(customerID)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
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

	if loggedErrors[0] != "Database error fetching customer" {
		t.Errorf("Expected error log 'Database error fetching customer', got '%s'", loggedErrors[0])
	}
}
