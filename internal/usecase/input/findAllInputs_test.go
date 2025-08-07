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

func TestFindAllInputs_Process_Success(t *testing.T) {
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

	// Mock data
	mockInputs := []models.Input{
		{
			ID:        uuid.New(),
			Name:      "Parafuso M6",
			Price:     2.50,
			Quantity:  100,
			InputType: "supplie",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Name:      "Serviço de Troca de Óleo",
			Price:     50.00,
			Quantity:  1,
			InputType: "service",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	inputRepoMock.FindAllFunc = func() ([]models.Input, error) {
		return mockInputs, nil
	}

	useCase := &FindAllInputs{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	result, err := useCase.Process()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 inputs, got %d", len(result))
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing find all inputs",
		"Found inputs in database",
		"Successfully mapped inputs to domain",
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

func TestFindAllInputs_Process_DatabaseError(t *testing.T) {
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

	expectedError := errors.New("database connection failed")
	inputRepoMock.FindAllFunc = func() ([]models.Input, error) {
		return nil, expectedError
	}

	useCase := &FindAllInputs{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
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

	if len(result) != 0 {
		t.Errorf("Expected empty result, got %d", len(result))
	}

	// Verifica se os logs corretos foram chamados
	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Processing find all inputs" {
		t.Errorf("Expected log message 'Processing find all inputs', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Database error finding all inputs" {
		t.Errorf("Expected error log 'Database error finding all inputs', got '%s'", loggedErrors[0])
	}
}

func TestFindAllInputs_Process_EmptyResult(t *testing.T) {
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

	// Mock empty result
	inputRepoMock.FindAllFunc = func() ([]models.Input, error) {
		return []models.Input{}, nil
	}

	useCase := &FindAllInputs{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	result, err := useCase.Process()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 inputs, got %d", len(result))
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing find all inputs",
		"Found inputs in database",
		"Successfully mapped inputs to domain",
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

func TestFindAllInputs_FetchInputsFromDB_Success(t *testing.T) {
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

	mockInputs := []models.Input{
		{
			ID:        uuid.New(),
			Name:      "Parafuso M6",
			Price:     2.50,
			Quantity:  100,
			InputType: "supplie",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	inputRepoMock.FindAllFunc = func() ([]models.Input, error) {
		return mockInputs, nil
	}

	useCase := &FindAllInputs{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	result, err := useCase.FetchInputsFromDB()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 input, got %d", len(result))
	}

	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Found inputs in database" {
		t.Errorf("Expected log message 'Found inputs in database', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) > 0 {
		t.Errorf("Expected no error logs, got %d", len(loggedErrors))
	}
}

func TestFindAllInputs_FetchInputsFromDB_Error(t *testing.T) {
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

	expectedError := errors.New("database timeout")
	inputRepoMock.FindAllFunc = func() ([]models.Input, error) {
		return nil, expectedError
	}

	useCase := &FindAllInputs{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	result, err := useCase.FetchInputsFromDB()

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty result, got %d", len(result))
	}

	if len(loggedInfo) > 0 {
		t.Errorf("Expected no info logs, got %d", len(loggedInfo))
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Database error finding all inputs" {
		t.Errorf("Expected error log 'Database error finding all inputs', got '%s'", loggedErrors[0])
	}
}
