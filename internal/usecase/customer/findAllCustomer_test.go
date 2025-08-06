package customer

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/test/mocks"
	"go.uber.org/zap"
)

func TestFindAllCustomer_Process_Success(t *testing.T) {
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

	// Mock data
	mockCustomers := []models.Customer{
		{
			ID:             uuid.New(),
			Name:           "João Silva",
			DocumentNumber: "12345678901",
			CustomerType:   "individual",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New(),
			Name:           "Maria Santos",
			DocumentNumber: "98765432100",
			CustomerType:   "individual",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	customerRepoMock.FindAllFunc = func() ([]models.Customer, error) {
		return mockCustomers, nil
	}

	useCase := &FindAllCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	result, err := useCase.Process()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 customers, got %d", len(result))
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing find all customers",
		"Successfully fetched customers from database",
		"Successfully mapped customers to domain",
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

func TestFindAllCustomer_Process_DatabaseError(t *testing.T) {
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

	expectedError := errors.New("database connection failed")
	customerRepoMock.FindAllFunc = func() ([]models.Customer, error) {
		return nil, expectedError
	}

	useCase := &FindAllCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	result, err := useCase.Process()

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

	if loggedInfo[0] != "Processing find all customers" {
		t.Errorf("Expected log message 'Processing find all customers', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Database error fetching customers" {
		t.Errorf("Expected error log 'Database error fetching customers', got '%s'", loggedErrors[0])
	}
}

func TestFindAllCustomer_Process_EmptyResult(t *testing.T) {
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

	// Mock empty result
	customerRepoMock.FindAllFunc = func() ([]models.Customer, error) {
		return []models.Customer{}, nil
	}

	useCase := &FindAllCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	result, err := useCase.Process()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 customers, got %d", len(result))
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing find all customers",
		"Successfully fetched customers from database",
		"Successfully mapped customers to domain",
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

func TestFindAllCustomer_FetchCustomersFromDB_Success(t *testing.T) {
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

	mockCustomers := []models.Customer{
		{
			ID:             uuid.New(),
			Name:           "João Silva",
			DocumentNumber: "12345678901",
			CustomerType:   "individual",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	customerRepoMock.FindAllFunc = func() ([]models.Customer, error) {
		return mockCustomers, nil
	}

	useCase := &FindAllCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	result, err := useCase.FetchCustomersFromDB()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 customer, got %d", len(result))
	}

	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Successfully fetched customers from database" {
		t.Errorf("Expected log message 'Successfully fetched customers from database', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) > 0 {
		t.Errorf("Expected no error logs, got %d", len(loggedErrors))
	}
}

func TestFindAllCustomer_FetchCustomersFromDB_Error(t *testing.T) {
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

	expectedError := errors.New("database timeout")
	customerRepoMock.FindAllFunc = func() ([]models.Customer, error) {
		return nil, expectedError
	}

	useCase := &FindAllCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	result, err := useCase.FetchCustomersFromDB()

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

	if loggedErrors[0] != "Database error fetching customers" {
		t.Errorf("Expected error log 'Database error fetching customers', got '%s'", loggedErrors[0])
	}
}
