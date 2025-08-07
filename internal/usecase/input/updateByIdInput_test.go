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

func TestUpdateByIdInput_Process_Success(t *testing.T) {
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
	existingInput := &models.Input{
		ID:          inputID,
		Name:        "Parafuso M6",
		Description: "Parafuso sextavado M6",
		Price:       2.50,
		Quantity:    100,
		InputType:   "supplie",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	updateEntity := &domain.Input{
		ID:          inputID,
		Name:        "Parafuso M6 Atualizado",
		Description: "Parafuso sextavado M6 atualizado",
		Price:       3.00,
		Quantity:    150,
		InputType:   "supplie",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	inputRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Input, error) {
		if id == inputID {
			return existingInput, nil
		}
		return nil, errors.New("input not found")
	}

	inputRepoMock.FindByNameFunc = func(name string) (*models.Input, error) {
		return nil, gorm.ErrRecordNotFound
	}

	inputRepoMock.UpdateFunc = func(input *models.Input) error {
		return nil
	}

	useCase := &UpdateByIdInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.Process(inputID, updateEntity)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing update input by ID",
		"Found existing input",
		"Input name is unique",
		"Updated input fields",
		"Input updated successfully",
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

func TestUpdateByIdInput_Process_InputNotFound(t *testing.T) {
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
	updateEntity := &domain.Input{
		ID:          inputID,
		Name:        "Parafuso M6 Atualizado",
		Description: "Parafuso sextavado M6 atualizado",
		Price:       3.00,
		Quantity:    150,
		InputType:   "supplie",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	inputRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Input, error) {
		return nil, errors.New("input not found")
	}

	useCase := &UpdateByIdInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.Process(inputID, updateEntity)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "input not found" {
		t.Errorf("Expected error 'input not found', got '%s'", err.Error())
	}

	// Verifica se os logs corretos foram chamados
	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Processing update input by ID" {
		t.Errorf("Expected log message 'Processing update input by ID', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Database error finding input to update" {
		t.Errorf("Expected error log 'Database error finding input to update', got '%s'", loggedErrors[0])
	}
}

func TestUpdateByIdInput_Process_NameAlreadyExists(t *testing.T) {
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
	existingInput := &models.Input{
		ID:          inputID,
		Name:        "Parafuso M6",
		Description: "Parafuso sextavado M6",
		Price:       2.50,
		Quantity:    100,
		InputType:   "supplie",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	updateEntity := &domain.Input{
		ID:          inputID,
		Name:        "Parafuso M8", // Nome diferente
		Description: "Parafuso sextavado M6 atualizado",
		Price:       3.00,
		Quantity:    150,
		InputType:   "supplie",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	existingInputWithSameName := &models.Input{
		ID:          uuid.New(), // ID diferente
		Name:        "Parafuso M8",
		Description: "Parafuso M8 existente",
		Price:       4.00,
		Quantity:    50,
		InputType:   "supplie",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	inputRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Input, error) {
		if id == inputID {
			return existingInput, nil
		}
		return nil, errors.New("input not found")
	}

	inputRepoMock.FindByNameFunc = func(name string) (*models.Input, error) {
		if name == "Parafuso M8" {
			return existingInputWithSameName, nil
		}
		return nil, gorm.ErrRecordNotFound
	}

	useCase := &UpdateByIdInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.Process(inputID, updateEntity)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "input name already exists" {
		t.Errorf("Expected error 'input name already exists', got '%s'", err.Error())
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing update input by ID",
		"Found existing input",
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

	if loggedErrors[0] != "Input name already exists" {
		t.Errorf("Expected error log 'Input name already exists', got '%s'", loggedErrors[0])
	}
}

func TestUpdateByIdInput_Process_ServiceTypeQuantityAdjustment(t *testing.T) {
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
	existingInput := &models.Input{
		ID:          inputID,
		Name:        "Serviço de Troca de Óleo",
		Description: "Serviço completo de troca de óleo",
		Price:       50.00,
		Quantity:    1,
		InputType:   "service",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	updateEntity := &domain.Input{
		ID:          inputID,
		Name:        "Serviço de Troca de Óleo Atualizado",
		Description: "Serviço completo de troca de óleo atualizado",
		Price:       60.00,
		Quantity:    5, // Deveria ser forçado para 1
		InputType:   "service",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	inputRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Input, error) {
		if id == inputID {
			return existingInput, nil
		}
		return nil, errors.New("input not found")
	}

	inputRepoMock.FindByNameFunc = func(name string) (*models.Input, error) {
		return nil, gorm.ErrRecordNotFound
	}

	inputRepoMock.UpdateFunc = func(input *models.Input) error {
		return nil
	}

	useCase := &UpdateByIdInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	err := useCase.Process(inputID, updateEntity)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing update input by ID",
		"Found existing input",
		"Input name is unique",
		"Forcing quantity to 1 for service type",
		"Updated input fields",
		"Input updated successfully",
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

func TestUpdateByIdInput_AdjustQuantityForInputType_Service(t *testing.T) {
	// Arrange
	inputRepoMock := &mocks.InputRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	useCase := &UpdateByIdInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	result := useCase.AdjustQuantityForInputType(5, "service")

	// Assert
	if result != 1 {
		t.Errorf("Expected quantity 1 for service type, got %d", result)
	}

	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Forcing quantity to 1 for service type" {
		t.Errorf("Expected log message 'Forcing quantity to 1 for service type', got '%s'", loggedInfo[0])
	}
}

func TestUpdateByIdInput_AdjustQuantityForInputType_Supplie(t *testing.T) {
	// Arrange
	inputRepoMock := &mocks.InputRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	useCase := &UpdateByIdInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	// Act
	result := useCase.AdjustQuantityForInputType(100, "supplie")

	// Assert
	if result != 100 {
		t.Errorf("Expected quantity 100 for supplie type, got %d", result)
	}

	if len(loggedInfo) > 0 {
		t.Errorf("Expected no info logs, got %d", len(loggedInfo))
	}
}

func TestUpdateByIdInput_SaveInputToDB_Success(t *testing.T) {
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

	useCase := &UpdateByIdInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	input := &models.Input{
		ID:          uuid.New(),
		Name:        "Parafuso M6",
		Description: "Parafuso sextavado M6",
		Price:       2.50,
		Quantity:    100,
		InputType:   "supplie",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
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

	if loggedInfo[0] != "Input updated successfully" {
		t.Errorf("Expected log message 'Input updated successfully', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) > 0 {
		t.Errorf("Expected no error logs, got %d", len(loggedErrors))
	}
}

func TestUpdateByIdInput_SaveInputToDB_Error(t *testing.T) {
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

	useCase := &UpdateByIdInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	input := &models.Input{
		ID:          uuid.New(),
		Name:        "Parafuso M6",
		Description: "Parafuso sextavado M6",
		Price:       2.50,
		Quantity:    100,
		InputType:   "supplie",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
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

	if loggedErrors[0] != "Database error updating input" {
		t.Errorf("Expected error log 'Database error updating input', got '%s'", loggedErrors[0])
	}
}
