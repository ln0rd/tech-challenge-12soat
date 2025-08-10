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

func TestRemoveInputFromOrder_Process_Success(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	inputRepoMock := &mocks.InputRepositoryMock{}
	orderInputRepoMock := &mocks.OrderInputRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}
	increaseQuantityInputMock := &input.IncreaseQuantityInput{
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
	quantityToRemove := 2

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
		Quantity:  5,
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

	// Mock OrderInput
	orderInput := &models.OrderInput{
		ID:         uuid.New(),
		OrderID:    orderID,
		InputID:    inputID,
		Quantity:   5,
		UnitPrice:  15.50,
		TotalPrice: 77.50,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	orderInputRepoMock.FindByOrderIDFunc = func(id uuid.UUID) ([]models.OrderInput, error) {
		if id == orderID {
			return []models.OrderInput{*orderInput}, nil
		}
		return []models.OrderInput{}, nil
	}

	// Mock Update OrderInput
	orderInputRepoMock.UpdateFunc = func(orderInput *models.OrderInput) error {
		return nil
	}

	// Mock Update Input quantity
	inputRepoMock.UpdateFunc = func(input *models.Input) error {
		return nil
	}

	useCase := &RemoveInputFromOrder{
		OrderRepository:       orderRepoMock,
		InputRepository:       inputRepoMock,
		OrderInputRepository:  orderInputRepoMock,
		Logger:                loggerMock,
		IncreaseQuantityInput: increaseQuantityInputMock,
	}

	// Act
	err := useCase.Process(orderID, inputID, quantityToRemove)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verifica se os logs de sucesso foram chamados
	expectedInfoLogs := []string{
		"Order found",
		"Input found",
		"Order input found",
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

func TestRemoveInputFromOrder_Process_OrderNotFound(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	inputRepoMock := &mocks.InputRepositoryMock{}
	orderInputRepoMock := &mocks.OrderInputRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}
	increaseQuantityInputMock := &input.IncreaseQuantityInput{
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
	quantityToRemove := 2

	orderRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Order, error) {
		return nil, errors.New("order not found")
	}

	useCase := &RemoveInputFromOrder{
		OrderRepository:       orderRepoMock,
		InputRepository:       inputRepoMock,
		OrderInputRepository:  orderInputRepoMock,
		Logger:                loggerMock,
		IncreaseQuantityInput: increaseQuantityInputMock,
	}

	// Act
	err := useCase.Process(orderID, inputID, quantityToRemove)

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

func TestRemoveInputFromOrder_Process_InputNotFound(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	inputRepoMock := &mocks.InputRepositoryMock{}
	orderInputRepoMock := &mocks.OrderInputRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}
	increaseQuantityInputMock := &input.IncreaseQuantityInput{
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
	quantityToRemove := 2

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

	useCase := &RemoveInputFromOrder{
		OrderRepository:       orderRepoMock,
		InputRepository:       inputRepoMock,
		OrderInputRepository:  orderInputRepoMock,
		Logger:                loggerMock,
		IncreaseQuantityInput: increaseQuantityInputMock,
	}

	// Act
	err := useCase.Process(orderID, inputID, quantityToRemove)

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

func TestRemoveInputFromOrder_Process_OrderInputNotFound(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	inputRepoMock := &mocks.InputRepositoryMock{}
	orderInputRepoMock := &mocks.OrderInputRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}
	increaseQuantityInputMock := &input.IncreaseQuantityInput{
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
	quantityToRemove := 2

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

	// Mock Input - sucesso
	input := &models.Input{
		ID:        inputID,
		Name:      "Test Input",
		InputType: "material",
		Quantity:  5,
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

	// Mock OrderInput - não encontrado
	orderInputRepoMock.FindByOrderIDFunc = func(id uuid.UUID) ([]models.OrderInput, error) {
		return []models.OrderInput{}, nil
	}

	useCase := &RemoveInputFromOrder{
		OrderRepository:       orderRepoMock,
		InputRepository:       inputRepoMock,
		OrderInputRepository:  orderInputRepoMock,
		Logger:                loggerMock,
		IncreaseQuantityInput: increaseQuantityInputMock,
	}

	// Act
	err := useCase.Process(orderID, inputID, quantityToRemove)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "order input not found" {
		t.Errorf("Expected error 'order input not found', got %v", err)
	}

	// Verifica se o log de erro foi chamado
	expectedErrorLog := "Order input not found"
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

func TestRemoveInputFromOrder_ValidateQuantityToRemove_Valid(t *testing.T) {
	// Arrange
	loggerMock := &mocks.LoggerMock{}
	useCase := &RemoveInputFromOrder{
		Logger: loggerMock,
	}

	// Act
	err := useCase.ValidateQuantityToRemove(3)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestRemoveInputFromOrder_ValidateQuantityToRemove_Invalid(t *testing.T) {
	// Arrange
	loggerMock := &mocks.LoggerMock{}
	var loggedErrors []string

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	useCase := &RemoveInputFromOrder{
		Logger: loggerMock,
	}

	// Act
	err := useCase.ValidateQuantityToRemove(0)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "quantity to remove must be greater than zero" {
		t.Errorf("Expected error 'quantity to remove must be greater than zero', got %v", err)
	}

	// Verifica se o log de erro foi chamado
	expectedErrorLog := "Invalid quantity to remove"
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

func TestRemoveInputFromOrder_ValidateOrderInputQuantity_Valid(t *testing.T) {
	// Arrange
	loggerMock := &mocks.LoggerMock{}
	useCase := &RemoveInputFromOrder{
		Logger: loggerMock,
	}

	orderInput := &models.OrderInput{
		ID:         uuid.New(),
		OrderID:    uuid.New(),
		InputID:    uuid.New(),
		Quantity:   5,
		UnitPrice:  15.50,
		TotalPrice: 77.50,
	}

	// Act
	err := useCase.ValidateOrderInputQuantity(orderInput, 3)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestRemoveInputFromOrder_ValidateOrderInputQuantity_InvalidQuantity(t *testing.T) {
	// Arrange
	loggerMock := &mocks.LoggerMock{}
	var loggedErrors []string

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	useCase := &RemoveInputFromOrder{
		Logger: loggerMock,
	}

	orderInput := &models.OrderInput{
		ID:         uuid.New(),
		OrderID:    uuid.New(),
		InputID:    uuid.New(),
		Quantity:   0,
		UnitPrice:  15.50,
		TotalPrice: 0.0,
	}

	// Act
	err := useCase.ValidateOrderInputQuantity(orderInput, 3)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "invalid quantity in order input" {
		t.Errorf("Expected error 'invalid quantity in order input', got %v", err)
	}

	// Verifica se o log de erro foi chamado
	expectedErrorLog := "Invalid quantity in order input"
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

func TestRemoveInputFromOrder_ValidateOrderInputQuantity_InsufficientQuantity(t *testing.T) {
	// Arrange
	loggerMock := &mocks.LoggerMock{}
	var loggedErrors []string

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	useCase := &RemoveInputFromOrder{
		Logger: loggerMock,
	}

	orderInput := &models.OrderInput{
		ID:         uuid.New(),
		OrderID:    uuid.New(),
		InputID:    uuid.New(),
		Quantity:   2,
		UnitPrice:  15.50,
		TotalPrice: 31.00,
	}

	// Act
	err := useCase.ValidateOrderInputQuantity(orderInput, 5)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "insufficient quantity in order input" {
		t.Errorf("Expected error 'insufficient quantity in order input', got %v", err)
	}

	// Verifica se o log de erro foi chamado
	expectedErrorLog := "Insufficient quantity in order input"
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

func TestRemoveInputFromOrder_CalculateNewOrderInputValues(t *testing.T) {
	// Arrange
	loggerMock := &mocks.LoggerMock{}
	useCase := &RemoveInputFromOrder{
		Logger: loggerMock,
	}

	orderInput := &models.OrderInput{
		ID:         uuid.New(),
		OrderID:    uuid.New(),
		InputID:    uuid.New(),
		Quantity:   10,
		UnitPrice:  15.50,
		TotalPrice: 155.00,
	}

	quantityToRemove := 3

	// Act
	newQuantity, newTotalPrice := useCase.CalculateNewOrderInputValues(orderInput, quantityToRemove)

	// Assert
	expectedQuantity := 7
	expectedTotalPrice := 108.50

	if newQuantity != expectedQuantity {
		t.Errorf("Expected quantity %d, got %d", expectedQuantity, newQuantity)
	}

	if newTotalPrice != expectedTotalPrice {
		t.Errorf("Expected total price %f, got %f", expectedTotalPrice, newTotalPrice)
	}
}

func TestRemoveInputFromOrder_FetchOrderInputFromDB_Success(t *testing.T) {
	// Arrange
	orderInputRepoMock := &mocks.OrderInputRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}
	var loggedInfo []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	orderID := uuid.New()
	inputID := uuid.New()

	orderInput := &models.OrderInput{
		ID:         uuid.New(),
		OrderID:    orderID,
		InputID:    inputID,
		Quantity:   5,
		UnitPrice:  15.50,
		TotalPrice: 77.50,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	orderInputRepoMock.FindByOrderIDFunc = func(id uuid.UUID) ([]models.OrderInput, error) {
		if id == orderID {
			return []models.OrderInput{*orderInput}, nil
		}
		return []models.OrderInput{}, nil
	}

	useCase := &RemoveInputFromOrder{
		OrderInputRepository: orderInputRepoMock,
		Logger:               loggerMock,
	}

	// Act
	result, err := useCase.FetchOrderInputFromDB(orderID, inputID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil { //nolint:staticcheck // SA5011: this is a test, we expect result to be non-nil
		t.Error("Expected order input, got nil")
	}

	if result.ID != orderInput.ID { //nolint:staticcheck // SA5011: this is a test, we already checked result is not nil above
		t.Errorf("Expected order input ID %s, got %s", orderInput.ID, result.ID)
	}

	// Verifica se o log de sucesso foi chamado
	expectedInfoLog := "Order input found"
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

func TestRemoveInputFromOrder_FetchOrderInputFromDB_NotFound(t *testing.T) {
	// Arrange
	orderInputRepoMock := &mocks.OrderInputRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}
	var loggedErrors []string

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	orderID := uuid.New()
	inputID := uuid.New()

	orderInputRepoMock.FindByOrderIDFunc = func(id uuid.UUID) ([]models.OrderInput, error) {
		return []models.OrderInput{}, nil
	}

	useCase := &RemoveInputFromOrder{
		OrderInputRepository: orderInputRepoMock,
		Logger:               loggerMock,
	}

	// Act
	result, err := useCase.FetchOrderInputFromDB(orderID, inputID)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result, got order input")
	}

	if err.Error() != "order input not found" {
		t.Errorf("Expected error 'order input not found', got %v", err)
	}

	// Verifica se o log de erro foi chamado
	expectedErrorLog := "Order input not found"
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
