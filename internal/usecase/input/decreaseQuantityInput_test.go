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

func TestDecreaseQuantityInput_Process_Success(t *testing.T) {
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

	inputRepoMock.UpdateFunc = func(input *models.Input) error {
		return nil
	}

	useCase := &DecreaseQuantityInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.Process(inputID, 30)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing decrease quantity for input",
		"Found input",
		"Calculated new quantity",
		"Input quantity decreased successfully",
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

func TestDecreaseQuantityInput_Process_InvalidQuantity(t *testing.T) {
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

	useCase := &DecreaseQuantityInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.Process(inputID, 0)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "quantity to decrease must be greater than zero" {
		t.Errorf("Expected error 'quantity to decrease must be greater than zero', got '%s'", err.Error())
	}

	// Verifica se os logs corretos foram chamados
	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Processing decrease quantity for input" {
		t.Errorf("Expected log message 'Processing decrease quantity for input', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Invalid quantity to decrease" {
		t.Errorf("Expected error log 'Invalid quantity to decrease', got '%s'", loggedErrors[0])
	}
}

func TestDecreaseQuantityInput_Process_InsufficientQuantity(t *testing.T) {
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
		Quantity:  50,
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

	useCase := &DecreaseQuantityInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.Process(inputID, 100)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "insufficient quantity" {
		t.Errorf("Expected error 'insufficient quantity', got '%s'", err.Error())
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing decrease quantity for input",
		"Found input",
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

	if loggedErrors[0] != "Insufficient quantity" {
		t.Errorf("Expected error log 'Insufficient quantity', got '%s'", loggedErrors[0])
	}
}

func TestDecreaseQuantityInput_Process_InputNotFound(t *testing.T) {
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

	useCase := &DecreaseQuantityInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.Process(inputID, 30)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "input not found" {
		t.Errorf("Expected error 'input not found', got '%s'", err.Error())
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing decrease quantity for input",
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

	if loggedErrors[0] != "Input not found" {
		t.Errorf("Expected error log 'Input not found', got '%s'", loggedErrors[0])
	}
}

func TestDecreaseQuantityInput_Process_DatabaseError(t *testing.T) {
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

	expectedError := errors.New("database connection failed")
	inputRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Input, error) {
		if id == inputID {
			return mockInput, nil
		}
		return nil, errors.New("input not found")
	}

	inputRepoMock.UpdateFunc = func(input *models.Input) error {
		return expectedError
	}

	useCase := &DecreaseQuantityInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.Process(inputID, 30)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing decrease quantity for input",
		"Found input",
		"Calculated new quantity",
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

	if loggedErrors[0] != "Database error updating input quantity" {
		t.Errorf("Expected error log 'Database error updating input quantity', got '%s'", loggedErrors[0])
	}
}

func TestDecreaseQuantityInput_ValidateQuantityToDecrease_Success(t *testing.T) {
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

	useCase := &DecreaseQuantityInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.ValidateQuantityToDecrease(30)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(loggedInfo) > 0 {
		t.Errorf("Expected no info logs, got %d", len(loggedInfo))
	}

	if len(loggedErrors) > 0 {
		t.Errorf("Expected no error logs, got %d", len(loggedErrors))
	}
}

func TestDecreaseQuantityInput_ValidateQuantityToDecrease_Invalid(t *testing.T) {
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

	useCase := &DecreaseQuantityInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.ValidateQuantityToDecrease(0)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "quantity to decrease must be greater than zero" {
		t.Errorf("Expected error 'quantity to decrease must be greater than zero', got '%s'", err.Error())
	}

	if len(loggedInfo) > 0 {
		t.Errorf("Expected no info logs, got %d", len(loggedInfo))
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Invalid quantity to decrease" {
		t.Errorf("Expected error log 'Invalid quantity to decrease', got '%s'", loggedErrors[0])
	}
}

func TestDecreaseQuantityInput_CalculateNewQuantity_Success(t *testing.T) {
	// Arrange
	inputRepoMock := &mocks.InputRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	useCase := &DecreaseQuantityInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	result, err := useCase.CalculateNewQuantity(100, 30)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != 70 {
		t.Errorf("Expected quantity 70, got %d", result)
	}

	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Calculated new quantity" {
		t.Errorf("Expected log message 'Calculated new quantity', got '%s'", loggedInfo[0])
	}
}

func TestDecreaseQuantityInput_CalculateNewQuantity_Insufficient(t *testing.T) {
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

	useCase := &DecreaseQuantityInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	result, err := useCase.CalculateNewQuantity(50, 100)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "insufficient quantity" {
		t.Errorf("Expected error 'insufficient quantity', got '%s'", err.Error())
	}

	if result != 0 {
		t.Errorf("Expected quantity 0, got %d", result)
	}

	if len(loggedInfo) > 0 {
		t.Errorf("Expected no info logs, got %d", len(loggedInfo))
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Insufficient quantity" {
		t.Errorf("Expected error log 'Insufficient quantity', got '%s'", loggedErrors[0])
	}
}

func TestDecreaseQuantityInput_UpdateInputQuantity_Success(t *testing.T) {
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

	inputRepoMock.UpdateFunc = func(input *models.Input) error {
		return nil
	}

	useCase := &DecreaseQuantityInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	input := &models.Input{
		ID:        uuid.New(),
		Name:      "Parafuso M6",
		Price:     2.50,
		Quantity:  100,
		InputType: "supplie",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	err := useCase.UpdateInputQuantity(input, 70)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if input.Quantity != 70 {
		t.Errorf("Expected quantity 70, got %d", input.Quantity)
	}

	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Input quantity decreased successfully" {
		t.Errorf("Expected log message 'Input quantity decreased successfully', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) > 0 {
		t.Errorf("Expected no error logs, got %d", len(loggedErrors))
	}
}

func TestDecreaseQuantityInput_UpdateInputQuantity_Error(t *testing.T) {
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

	expectedError := errors.New("update constraint violation")
	inputRepoMock.UpdateFunc = func(input *models.Input) error {
		return expectedError
	}

	useCase := &DecreaseQuantityInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	input := &models.Input{
		ID:        uuid.New(),
		Name:      "Parafuso M6",
		Price:     2.50,
		Quantity:  100,
		InputType: "supplie",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	err := useCase.UpdateInputQuantity(input, 70)

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

	if loggedErrors[0] != "Database error updating input quantity" {
		t.Errorf("Expected error log 'Database error updating input quantity', got '%s'", loggedErrors[0])
	}
}
