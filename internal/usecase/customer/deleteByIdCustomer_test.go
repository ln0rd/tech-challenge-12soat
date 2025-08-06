package customer

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/test/mocks"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestDeleteByIdCustomer_Process_Success(t *testing.T) {
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

	customerRepoMock.DeleteFunc = func(id uuid.UUID) error {
		if id == customerID {
			return nil
		}
		return gorm.ErrRecordNotFound
	}

	useCase := &DeleteByIdCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.Process(customerID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing delete customer by ID",
		"Customer deleted from database",
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

func TestDeleteByIdCustomer_Process_CustomerNotFound(t *testing.T) {
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

	customerRepoMock.DeleteFunc = func(id uuid.UUID) error {
		return gorm.ErrRecordNotFound
	}

	useCase := &DeleteByIdCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.Process(customerID)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err != gorm.ErrRecordNotFound {
		t.Errorf("Expected error %v, got %v", gorm.ErrRecordNotFound, err)
	}

	// Verifica se os logs corretos foram chamados
	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Processing delete customer by ID" {
		t.Errorf("Expected log message 'Processing delete customer by ID', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Customer not found for deletion" {
		t.Errorf("Expected error log 'Customer not found for deletion', got '%s'", loggedErrors[0])
	}
}

func TestDeleteByIdCustomer_Process_DatabaseError(t *testing.T) {
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

	customerRepoMock.DeleteFunc = func(id uuid.UUID) error {
		return expectedError
	}

	useCase := &DeleteByIdCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.Process(customerID)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	// Verifica se os logs corretos foram chamados
	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Processing delete customer by ID" {
		t.Errorf("Expected log message 'Processing delete customer by ID', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Database error deleting customer" {
		t.Errorf("Expected error log 'Database error deleting customer', got '%s'", loggedErrors[0])
	}
}

func TestDeleteByIdCustomer_DeleteCustomerFromDB_Success(t *testing.T) {
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

	customerRepoMock.DeleteFunc = func(id uuid.UUID) error {
		if id == customerID {
			return nil
		}
		return gorm.ErrRecordNotFound
	}

	useCase := &DeleteByIdCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.DeleteCustomerFromDB(customerID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Customer deleted from database" {
		t.Errorf("Expected log message 'Customer deleted from database', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) > 0 {
		t.Errorf("Expected no error logs, got %d", len(loggedErrors))
	}
}

func TestDeleteByIdCustomer_DeleteCustomerFromDB_NotFound(t *testing.T) {
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

	customerRepoMock.DeleteFunc = func(id uuid.UUID) error {
		return gorm.ErrRecordNotFound
	}

	useCase := &DeleteByIdCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.DeleteCustomerFromDB(customerID)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err != gorm.ErrRecordNotFound {
		t.Errorf("Expected error %v, got %v", gorm.ErrRecordNotFound, err)
	}

	if len(loggedInfo) > 0 {
		t.Errorf("Expected no info logs, got %d", len(loggedInfo))
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Customer not found for deletion" {
		t.Errorf("Expected error log 'Customer not found for deletion', got '%s'", loggedErrors[0])
	}
}

func TestDeleteByIdCustomer_DeleteCustomerFromDB_DatabaseError(t *testing.T) {
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

	customerRepoMock.DeleteFunc = func(id uuid.UUID) error {
		return expectedError
	}

	useCase := &DeleteByIdCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.DeleteCustomerFromDB(customerID)

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

	if loggedErrors[0] != "Database error deleting customer" {
		t.Errorf("Expected error log 'Database error deleting customer', got '%s'", loggedErrors[0])
	}
}
