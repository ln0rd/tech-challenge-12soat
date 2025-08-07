package input

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/test/mocks"
	"go.uber.org/zap"
)

func TestDeleteByIdInput_Process_Success(t *testing.T) {
	// Arrange
	inputRepoMock := &mocks.InputRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	inputID := uuid.New()

	inputRepoMock.DeleteFunc = func(id uuid.UUID) error {
		if id == inputID {
			return nil
		}
		return errors.New("input not found")
	}

	useCase := &DeleteByIdInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.Process(inputID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing delete input by ID",
		"Input deleted successfully",
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

func TestDeleteByIdInput_Process_DatabaseError(t *testing.T) {
	// Arrange
	inputRepoMock := &mocks.InputRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	inputID := uuid.New()
	expectedError := errors.New("database connection failed")

	inputRepoMock.DeleteFunc = func(id uuid.UUID) error {
		return expectedError
	}

	useCase := &DeleteByIdInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.Process(inputID)

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

	if loggedInfo[0] != "Processing delete input by ID" {
		t.Errorf("Expected log message 'Processing delete input by ID', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Database error deleting input" {
		t.Errorf("Expected error log 'Database error deleting input', got '%s'", loggedErrors[0])
	}
}

func TestDeleteByIdInput_DeleteInputFromDB_Success(t *testing.T) {
	// Arrange
	inputRepoMock := &mocks.InputRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	inputID := uuid.New()

	inputRepoMock.DeleteFunc = func(id uuid.UUID) error {
		if id == inputID {
			return nil
		}
		return errors.New("input not found")
	}

	useCase := &DeleteByIdInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.DeleteInputFromDB(inputID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Input deleted successfully" {
		t.Errorf("Expected log message 'Input deleted successfully', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) > 0 {
		t.Errorf("Expected no error logs, got %d", len(loggedErrors))
	}
}

func TestDeleteByIdInput_DeleteInputFromDB_Error(t *testing.T) {
	// Arrange
	inputRepoMock := &mocks.InputRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	inputID := uuid.New()
	expectedError := errors.New("database timeout")

	inputRepoMock.DeleteFunc = func(id uuid.UUID) error {
		return expectedError
	}

	useCase := &DeleteByIdInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.DeleteInputFromDB(inputID)

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

	if loggedErrors[0] != "Database error deleting input" {
		t.Errorf("Expected error log 'Database error deleting input', got '%s'", loggedErrors[0])
	}
}
