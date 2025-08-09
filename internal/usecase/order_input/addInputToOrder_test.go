package order_input

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/test/mocks"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/input"
	"go.uber.org/zap"
)

func TestAddInputToOrder_Process_Success(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	inputRepoMock := &mocks.InputRepositoryMock{}
	orderInputRepoMock := &mocks.OrderInputRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}
	decreaseQuantityInputMock := &input.DecreaseQuantityInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	orderID := uuid.New()
	inputID := uuid.New()
	quantity := 3

	// Mock Order
	order := &models.Order{
		ID:        orderID,
		Status:    "open",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	orderRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Order, error) {
		if id == orderID {
			return order, nil
		}
		return nil, errors.New("order not found")
	}

	// Mock Input
	input := &models.Input{
		ID:        inputID,
		Name:      "Test Input",
		InputType: "material",
		Quantity:  10,
		Price:     15.50,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	inputRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Input, error) {
		if id == inputID {
			return input, nil
		}
		return nil, errors.New("input not found")
	}

	// Mock OrderInput - não existe ainda
	orderInputRepoMock.FindByOrderIDFunc = func(id uuid.UUID) ([]models.OrderInput, error) {
		return []models.OrderInput{}, nil
	}

	// Mock Create OrderInput
	orderInputRepoMock.CreateFunc = func(orderInput *models.OrderInput) error {
		return nil
	}

	// Mock Update Input quantity
	inputRepoMock.UpdateFunc = func(input *models.Input) error {
		return nil
	}

	useCase := &AddInputToOrder{
		OrderRepository:       orderRepoMock,
		InputRepository:       inputRepoMock,
		OrderInputRepository:  orderInputRepoMock,
		Logger:                loggerMock,
		DecreaseQuantityInput: decreaseQuantityInputMock,
	}

	// Act
	err := useCase.Process(orderID, inputID, quantity)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verifica se os logs de sucesso foram chamados
	expectedInfoLogs := []string{
		"Order found",
		"Input found",
		"Input price retrieved",
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

func TestAddInputToOrder_Process_OrderNotFound(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	inputRepoMock := &mocks.InputRepositoryMock{}
	orderInputRepoMock := &mocks.OrderInputRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}
	decreaseQuantityInputMock := &input.DecreaseQuantityInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	orderID := uuid.New()
	inputID := uuid.New()
	quantity := 3

	orderRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Order, error) {
		return nil, errors.New("order not found")
	}

	useCase := &AddInputToOrder{
		OrderRepository:       orderRepoMock,
		InputRepository:       inputRepoMock,
		OrderInputRepository:  orderInputRepoMock,
		Logger:                loggerMock,
		DecreaseQuantityInput: decreaseQuantityInputMock,
	}

	// Act
	err := useCase.Process(orderID, inputID, quantity)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "order not found" {
		t.Errorf("Expected error 'order not found', got %v", err)
	}

	// Verifica se o log de erro foi chamado
	expectedErrorLog := "Order not found"
	found := false
	for _, actualLog := range loggedErrors {
		if actualLog == expectedErrorLog {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected error log message '%s' not found", expectedErrorLog)
	}
}

func TestAddInputToOrder_Process_InputNotFound(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	inputRepoMock := &mocks.InputRepositoryMock{}
	orderInputRepoMock := &mocks.OrderInputRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}
	decreaseQuantityInputMock := &input.DecreaseQuantityInput{
		InputRepository: inputRepoMock,
		Logger:          loggerMock,
	}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	orderID := uuid.New()
	inputID := uuid.New()
	quantity := 3

	// Mock Order - sucesso
	order := &models.Order{
		ID:        orderID,
		Status:    "open",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	orderRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Order, error) {
		if id == orderID {
			return order, nil
		}
		return nil, errors.New("order not found")
	}

	// Mock Input - erro
	inputRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Input, error) {
		return nil, errors.New("input not found")
	}

	useCase := &AddInputToOrder{
		OrderRepository:       orderRepoMock,
		InputRepository:       inputRepoMock,
		OrderInputRepository:  orderInputRepoMock,
		Logger:                loggerMock,
		DecreaseQuantityInput: decreaseQuantityInputMock,
	}

	// Act
	err := useCase.Process(orderID, inputID, quantity)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "input not found" {
		t.Errorf("Expected error 'input not found', got %v", err)
	}

	// Verifica se o log de erro foi chamado
	expectedErrorLog := "Input not found"
	found := false
	for _, actualLog := range loggedErrors {
		if actualLog == expectedErrorLog {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected error log message '%s' not found", expectedErrorLog)
	}
}

func TestAddInputToOrder_ValidateQuantity_Valid(t *testing.T) {
	// Arrange
	loggerMock := &mocks.LoggerMock{}
	useCase := &AddInputToOrder{
		Logger: loggerMock,
	}

	// Act
	err := useCase.ValidateQuantity(5)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestAddInputToOrder_ValidateQuantity_Invalid(t *testing.T) {
	// Arrange
	loggerMock := &mocks.LoggerMock{}
	var loggedErrors []string

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	useCase := &AddInputToOrder{
		Logger: loggerMock,
	}

	// Act
	err := useCase.ValidateQuantity(0)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "quantity must be greater than zero" {
		t.Errorf("Expected error 'quantity must be greater than zero', got %v", err)
	}

	// Verifica se o log de erro foi chamado
	expectedErrorLog := "Invalid quantity"
	found := false
	for _, actualLog := range loggedErrors {
		if actualLog == expectedErrorLog {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected error log message '%s' not found", expectedErrorLog)
	}
}

func TestAddInputToOrder_ValidateInputAvailability_SufficientQuantity(t *testing.T) {
	// Arrange
	loggerMock := &mocks.LoggerMock{}
	useCase := &AddInputToOrder{
		Logger: loggerMock,
	}

	input := &models.Input{
		ID:        uuid.New(),
		Name:      "Test Input",
		InputType: "material",
		Quantity:  10,
		Price:     15.50,
	}

	// Act
	err := useCase.ValidateInputAvailability(input, 5)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestAddInputToOrder_ValidateInputAvailability_InsufficientQuantity(t *testing.T) {
	// Arrange
	loggerMock := &mocks.LoggerMock{}
	var loggedErrors []string

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	useCase := &AddInputToOrder{
		Logger: loggerMock,
	}

	input := &models.Input{
		ID:        uuid.New(),
		Name:      "Test Input",
		InputType: "material",
		Quantity:  3,
		Price:     15.50,
	}

	// Act
	err := useCase.ValidateInputAvailability(input, 5)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "insufficient input quantity" {
		t.Errorf("Expected error 'insufficient input quantity', got %v", err)
	}

	// Verifica se o log de erro foi chamado
	expectedErrorLog := "Insufficient input quantity"
	found := false
	for _, actualLog := range loggedErrors {
		if actualLog == expectedErrorLog {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected error log message '%s' not found", expectedErrorLog)
	}
}

func TestAddInputToOrder_ValidateInputAvailability_ServiceType(t *testing.T) {
	// Arrange
	loggerMock := &mocks.LoggerMock{}
	var loggedInfo []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	useCase := &AddInputToOrder{
		Logger: loggerMock,
	}

	input := &models.Input{
		ID:        uuid.New(),
		Name:      "Test Service",
		InputType: "service",
		Quantity:  0,
		Price:     50.00,
	}

	// Act
	err := useCase.ValidateInputAvailability(input, 5)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verifica se o log de info foi chamado
	expectedInfoLog := "Input is service type, skipping stock control"
	found := false
	for _, actualLog := range loggedInfo {
		if actualLog == expectedInfoLog {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected info log message '%s' not found", expectedInfoLog)
	}
}

func TestAddInputToOrder_ValidateInputPrice_Valid(t *testing.T) {
	// Arrange
	loggerMock := &mocks.LoggerMock{}
	var loggedInfo []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	useCase := &AddInputToOrder{
		Logger: loggerMock,
	}

	input := &models.Input{
		ID:        uuid.New(),
		Name:      "Test Input",
		InputType: "material",
		Quantity:  10,
		Price:     25.50,
	}

	// Act
	price, err := useCase.ValidateInputPrice(input)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if price != 25.50 {
		t.Errorf("Expected price 25.50, got %f", price)
	}

	// Verifica se o log de info foi chamado
	expectedInfoLog := "Input price retrieved"
	found := false
	for _, actualLog := range loggedInfo {
		if actualLog == expectedInfoLog {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected info log message '%s' not found", expectedInfoLog)
	}
}

func TestAddInputToOrder_ValidateInputPrice_Invalid(t *testing.T) {
	// Arrange
	loggerMock := &mocks.LoggerMock{}
	var loggedErrors []string

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	useCase := &AddInputToOrder{
		Logger: loggerMock,
	}

	input := &models.Input{
		ID:        uuid.New(),
		Name:      "Test Input",
		InputType: "material",
		Quantity:  10,
		Price:     0.0,
	}

	// Act
	price, err := useCase.ValidateInputPrice(input)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if price != 0 {
		t.Errorf("Expected price 0, got %f", price)
	}

	if err.Error() != "input has invalid price" {
		t.Errorf("Expected error 'input has invalid price', got %v", err)
	}

	// Verifica se o log de erro foi chamado
	expectedErrorLog := "Input has invalid price"
	found := false
	for _, actualLog := range loggedErrors {
		if actualLog == expectedErrorLog {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected error log message '%s' not found", expectedErrorLog)
	}
}
