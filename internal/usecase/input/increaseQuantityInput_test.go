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

func TestIncreaseQuantityInput_Process_Success(t *testing.T) {
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

	useCase := &IncreaseQuantityInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.Process(inputID, 50)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing increase quantity for input",
		"Found input",
		"Calculated new quantity",
		"Input quantity increased successfully",
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

func TestIncreaseQuantityInput_Process_InvalidQuantity(t *testing.T) {
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

	useCase := &IncreaseQuantityInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.Process(inputID, 0)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "quantity to increase must be greater than zero" {
		t.Errorf("Expected error 'quantity to increase must be greater than zero', got '%s'", err.Error())
	}

	// Verifica se os logs corretos foram chamados
	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Processing increase quantity for input" {
		t.Errorf("Expected log message 'Processing increase quantity for input', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Invalid quantity to increase" {
		t.Errorf("Expected error log 'Invalid quantity to increase', got '%s'", loggedErrors[0])
	}
}

func TestIncreaseQuantityInput_Process_InputNotFound(t *testing.T) {
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

	useCase := &IncreaseQuantityInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.Process(inputID, 50)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "input not found" {
		t.Errorf("Expected error 'input not found', got '%s'", err.Error())
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing increase quantity for input",
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

func TestIncreaseQuantityInput_Process_DatabaseError(t *testing.T) {
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

	useCase := &IncreaseQuantityInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.Process(inputID, 50)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing increase quantity for input",
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

func TestIncreaseQuantityInput_ValidateQuantityToIncrease_Success(t *testing.T) {
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

	useCase := &IncreaseQuantityInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.ValidateQuantityToIncrease(50)

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

func TestIncreaseQuantityInput_ValidateQuantityToIncrease_Invalid(t *testing.T) {
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

	useCase := &IncreaseQuantityInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.ValidateQuantityToIncrease(0)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "quantity to increase must be greater than zero" {
		t.Errorf("Expected error 'quantity to increase must be greater than zero', got '%s'", err.Error())
	}

	if len(loggedInfo) > 0 {
		t.Errorf("Expected no info logs, got %d", len(loggedInfo))
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Invalid quantity to increase" {
		t.Errorf("Expected error log 'Invalid quantity to increase', got '%s'", loggedErrors[0])
	}
}

func TestIncreaseQuantityInput_CalculateNewQuantity(t *testing.T) {
	// Arrange
	inputRepoMock := &mocks.InputRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	useCase := &IncreaseQuantityInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	result := useCase.CalculateNewQuantity(100, 50)

	// Assert
	if result != 150 {
		t.Errorf("Expected quantity 150, got %d", result)
	}

	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Calculated new quantity" {
		t.Errorf("Expected log message 'Calculated new quantity', got '%s'", loggedInfo[0])
	}
}

func TestIncreaseQuantityInput_UpdateInputQuantity_Success(t *testing.T) {
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

	useCase := &IncreaseQuantityInput{
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
	err := useCase.UpdateInputQuantity(input, 150)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if input.Quantity != 150 {
		t.Errorf("Expected quantity 150, got %d", input.Quantity)
	}

	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Input quantity increased successfully" {
		t.Errorf("Expected log message 'Input quantity increased successfully', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) > 0 {
		t.Errorf("Expected no error logs, got %d", len(loggedErrors))
	}
}

func TestIncreaseQuantityInput_UpdateInputQuantity_Error(t *testing.T) {
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

	useCase := &IncreaseQuantityInput{
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
	err := useCase.UpdateInputQuantity(input, 150)

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
