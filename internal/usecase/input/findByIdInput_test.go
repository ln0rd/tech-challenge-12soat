package input

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/test/mocks"
	"go.uber.org/zap"
)

func TestFindByIdInput_Process_Success(t *testing.T) {
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
	mockInput := &models.Input{
		ID:        inputID,
		Name:      "Parafuso M6",
		Price:     2.50,
		Quantity:  100,
		InputType: "supplie",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	inputRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Input, error) {
		if id == inputID {
			return mockInput, nil
		}
		return nil, errors.New("input not found")
	}

	useCase := &FindByIdInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	result, err := useCase.Process(inputID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected result, got nil")
	}

	if result.ID != inputID {
		t.Errorf("Expected input ID %s, got %s", inputID, result.ID)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing find input by ID",
		"Found input in database",
		"Successfully mapped input to domain",
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

func TestFindByIdInput_Process_InputNotFound(t *testing.T) {
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

	inputRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Input, error) {
		return nil, errors.New("input not found")
	}

	useCase := &FindByIdInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	result, err := useCase.Process(inputID)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "input not found" {
		t.Errorf("Expected error 'input not found', got '%s'", err.Error())
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	// Verifica se os logs corretos foram chamados
	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Processing find input by ID" {
		t.Errorf("Expected log message 'Processing find input by ID', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Database error finding input by ID" {
		t.Errorf("Expected error log 'Database error finding input by ID', got '%s'", loggedErrors[0])
	}
}

func TestFindByIdInput_Process_DatabaseError(t *testing.T) {
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

	inputRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Input, error) {
		return nil, expectedError
	}

	useCase := &FindByIdInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	result, err := useCase.Process(inputID)

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

	if loggedInfo[0] != "Processing find input by ID" {
		t.Errorf("Expected log message 'Processing find input by ID', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Database error finding input by ID" {
		t.Errorf("Expected error log 'Database error finding input by ID', got '%s'", loggedErrors[0])
	}
}

func TestFindByIdInput_FetchInputFromDB_Success(t *testing.T) {
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
	mockInput := &models.Input{
		ID:        inputID,
		Name:      "Parafuso M6",
		Price:     2.50,
		Quantity:  100,
		InputType: "supplie",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	inputRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Input, error) {
		if id == inputID {
			return mockInput, nil
		}
		return nil, errors.New("input not found")
	}

	useCase := &FindByIdInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	result, err := useCase.FetchInputFromDB(inputID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected result, got nil")
	}

	if result.ID != inputID {
		t.Errorf("Expected input ID %s, got %s", inputID, result.ID)
	}

	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Found input in database" {
		t.Errorf("Expected log message 'Found input in database', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) > 0 {
		t.Errorf("Expected no error logs, got %d", len(loggedErrors))
	}
}

func TestFindByIdInput_FetchInputFromDB_Error(t *testing.T) {
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

	inputRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Input, error) {
		return nil, expectedError
	}

	useCase := &FindByIdInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	result, err := useCase.FetchInputFromDB(inputID)

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

	if loggedErrors[0] != "Database error finding input by ID" {
		t.Errorf("Expected error log 'Database error finding input by ID', got '%s'", loggedErrors[0])
	}
}
