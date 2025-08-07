package input

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/input"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/test/mocks"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestCreateInput_Process_Success(t *testing.T) {
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

	inputRepoMock.FindByNameFunc = func(name string) (*models.Input, error) {
		return nil, gorm.ErrRecordNotFound
	}

	inputRepoMock.CreateFunc = func(input *models.Input) error {
		return nil
	}

	useCase := &CreateInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	input := &domain.Input{
		ID:        uuid.New(),
		Name:      "Parafuso M6",
		Price:     2.50,
		Quantity:  100,
		InputType: "supplie",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	err := useCase.Process(input)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing input creation",
		"Input name is unique",
		"Model created",
		"Input created in database",
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

func TestCreateInput_Process_NameAlreadyExists(t *testing.T) {
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

	existingInput := &models.Input{
		ID:        uuid.New(),
		Name:      "Parafuso M6",
		Price:     2.50,
		Quantity:  50,
		InputType: "supplie",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	inputRepoMock.FindByNameFunc = func(name string) (*models.Input, error) {
		if name == "Parafuso M6" {
			return existingInput, nil
		}
		return nil, gorm.ErrRecordNotFound
	}

	useCase := &CreateInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	input := &domain.Input{
		ID:        uuid.New(),
		Name:      "Parafuso M6",
		Price:     3.00,
		Quantity:  100,
		InputType: "supplie",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	err := useCase.Process(input)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "input name already exists" {
		t.Errorf("Expected error 'input name already exists', got '%s'", err.Error())
	}

	// Verifica se os logs corretos foram chamados
	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Processing input creation" {
		t.Errorf("Expected log message 'Processing input creation', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Input name already exists" {
		t.Errorf("Expected error log 'Input name already exists', got '%s'", loggedErrors[0])
	}
}

func TestCreateInput_Process_DatabaseError(t *testing.T) {
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

	inputRepoMock.FindByNameFunc = func(name string) (*models.Input, error) {
		return nil, gorm.ErrRecordNotFound
	}

	expectedError := errors.New("database connection failed")
	inputRepoMock.CreateFunc = func(input *models.Input) error {
		return expectedError
	}

	useCase := &CreateInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	input := &domain.Input{
		ID:        uuid.New(),
		Name:      "Parafuso M6",
		Price:     2.50,
		Quantity:  100,
		InputType: "supplie",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	err := useCase.Process(input)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing input creation",
		"Input name is unique",
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

	if loggedErrors[0] != "Database error creating input" {
		t.Errorf("Expected error log 'Database error creating input', got '%s'", loggedErrors[0])
	}
}

func TestCreateInput_ValidateInputNameUniqueness_Success(t *testing.T) {
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

	inputRepoMock.FindByNameFunc = func(name string) (*models.Input, error) {
		return nil, gorm.ErrRecordNotFound
	}

	useCase := &CreateInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.ValidateInputNameUniqueness("Parafuso M6")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Input name is unique" {
		t.Errorf("Expected log message 'Input name is unique', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) > 0 {
		t.Errorf("Expected no error logs, got %d", len(loggedErrors))
	}
}

func TestCreateInput_ValidateInputNameUniqueness_AlreadyExists(t *testing.T) {
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

	existingInput := &models.Input{
		ID:        uuid.New(),
		Name:      "Parafuso M6",
		Price:     2.50,
		Quantity:  50,
		InputType: "supplie",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	inputRepoMock.FindByNameFunc = func(name string) (*models.Input, error) {
		if name == "Parafuso M6" {
			return existingInput, nil
		}
		return nil, gorm.ErrRecordNotFound
	}

	useCase := &CreateInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.ValidateInputNameUniqueness("Parafuso M6")

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "input name already exists" {
		t.Errorf("Expected error 'input name already exists', got '%s'", err.Error())
	}

	if len(loggedInfo) > 0 {
		t.Errorf("Expected no info logs, got %d", len(loggedInfo))
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Input name already exists" {
		t.Errorf("Expected error log 'Input name already exists', got '%s'", loggedErrors[0])
	}
}

func TestCreateInput_SaveInputToDB_Success(t *testing.T) {
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

	inputRepoMock.CreateFunc = func(input *models.Input) error {
		return nil
	}

	useCase := &CreateInput{
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
	err := useCase.SaveInputToDB(input)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Input created in database" {
		t.Errorf("Expected log message 'Input created in database', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) > 0 {
		t.Errorf("Expected no error logs, got %d", len(loggedErrors))
	}
}

func TestCreateInput_SaveInputToDB_Error(t *testing.T) {
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

	expectedError := errors.New("database constraint violation")
	inputRepoMock.CreateFunc = func(input *models.Input) error {
		return expectedError
	}

	useCase := &CreateInput{
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
	err := useCase.SaveInputToDB(input)

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

	if loggedErrors[0] != "Database error creating input" {
		t.Errorf("Expected error log 'Database error creating input', got '%s'", loggedErrors[0])
	}
}
