package customer

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/costumer"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/test/mocks"
	"go.uber.org/zap"
)

func TestCreateCustomer_Process_Success(t *testing.T) {
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

	customerRepoMock.CreateFunc = func(customer *models.Customer) error {
		// Simula sucesso na criação
		return nil
	}

	useCase := &CreateCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	customer := &domain.Customer{
		ID:             uuid.New(),
		Name:           "João Silva",
		DocumentNumber: "12345678901",
		CustomerType:   "individual",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Act
	err := useCase.Process(customer)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verifica se os logs corretos foram chamados
	if len(loggedInfo) < 2 {
		t.Errorf("Expected at least 2 info logs, got %d", len(loggedInfo))
	}

	expectedInfoLogs := []string{
		"Processing customer creation",
		"Model created",
		"Customer created in database",
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

func TestCreateCustomer_Process_DatabaseError(t *testing.T) {
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
	customerRepoMock.CreateFunc = func(customer *models.Customer) error {
		return expectedError
	}

	useCase := &CreateCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	customer := &domain.Customer{
		ID:             uuid.New(),
		Name:           "João Silva",
		DocumentNumber: "12345678901",
		CustomerType:   "individual",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Act
	err := useCase.Process(customer)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	// Verifica se os logs corretos foram chamados
	if len(loggedInfo) < 2 {
		t.Errorf("Expected at least 2 info logs, got %d", len(loggedInfo))
	}

	expectedInfoLogs := []string{
		"Processing customer creation",
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

	if loggedErrors[0] != "Database error creating customer" {
		t.Errorf("Expected error log 'Database error creating customer', got '%s'", loggedErrors[0])
	}
}

func TestCreateCustomer_SaveCustomerToDB_Success(t *testing.T) {
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

	customerRepoMock.CreateFunc = func(customer *models.Customer) error {
		return nil
	}

	useCase := &CreateCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	customer := &models.Customer{
		ID:             uuid.New(),
		Name:           "João Silva",
		DocumentNumber: "12345678901",
		CustomerType:   "individual",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Act
	err := useCase.SaveCustomerToDB(customer)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Customer created in database" {
		t.Errorf("Expected log message 'Customer created in database', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) > 0 {
		t.Errorf("Expected no error logs, got %d", len(loggedErrors))
	}
}

func TestCreateCustomer_SaveCustomerToDB_Error(t *testing.T) {
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

	expectedError := errors.New("database constraint violation")
	customerRepoMock.CreateFunc = func(customer *models.Customer) error {
		return expectedError
	}

	useCase := &CreateCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	customer := &models.Customer{
		ID:             uuid.New(),
		Name:           "João Silva",
		DocumentNumber: "12345678901",
		CustomerType:   "individual",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Act
	err := useCase.SaveCustomerToDB(customer)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	if len(loggedInfo) > 0 {
		t.Errorf("Expected no info logs, got %d", len(loggedInfo))
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Database error creating customer" {
		t.Errorf("Expected error log 'Database error creating customer', got '%s'", loggedErrors[0])
	}
}
