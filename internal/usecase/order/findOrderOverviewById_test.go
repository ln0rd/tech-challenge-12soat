package order

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/test/mocks"
	"go.uber.org/zap"
)

func TestFindOrderOverviewById_Process_Success(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	orderInputRepoMock := &mocks.OrderInputRepositoryMock{}
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
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

	orderID := uuid.New()
	customerID := uuid.New()
	vehicleID := uuid.New()
	inputID1 := uuid.New()
	inputID2 := uuid.New()

	mockOrder := &models.Order{
		ID:         orderID,
		CustomerID: customerID,
		VehicleID:  vehicleID,
		Status:     "Completed",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	mockVehicle := &models.Vehicle{
		ID:                          vehicleID,
		CustomerID:                  customerID,
		NumberPlate:                 "ABC1234",
		Brand:                       "Toyota",
		Model:                       "Corolla",
		ReleaseYear:                 2020,
		VehicleIdentificationNumber: "VIN123456789",
		Color:                       "Prata",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	mockOrderInputs := []models.OrderInput{
		{
			ID:         uuid.New(),
			OrderID:    orderID,
			InputID:    inputID1,
			Quantity:   2,
			UnitPrice:  50.0,
			TotalPrice: 100.0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         uuid.New(),
			OrderID:    orderID,
			InputID:    inputID2,
			Quantity:   1,
			UnitPrice:  75.0,
			TotalPrice: 75.0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	mockInput1 := &models.Input{
		ID:          inputID1,
		Name:        "Óleo de Motor",
		Description: "Óleo de motor sintético",
		Price:       50.0,
		Quantity:    100,
		InputType:   "Lubrificante",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockInput2 := &models.Input{
		ID:          inputID2,
		Name:        "Filtro de Ar",
		Description: "Filtro de ar do motor",
		Price:       75.0,
		Quantity:    50,
		InputType:   "Filtro",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockStatusHistory := []models.OrderStatusHistory{
		{
			ID:        uuid.New(),
			OrderID:   orderID,
			Status:    "Received",
			StartedAt: time.Now().Add(-2 * time.Hour),
			EndedAt:   &[]time.Time{time.Now().Add(-1 * time.Hour)}[0],
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			OrderID:   orderID,
			Status:    "In Progress",
			StartedAt: time.Now().Add(-1 * time.Hour),
			EndedAt:   &[]time.Time{time.Now().Add(-30 * time.Minute)}[0],
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			OrderID:   orderID,
			Status:    "Completed",
			StartedAt: time.Now().Add(-30 * time.Minute),
			EndedAt:   nil, // Status atual
			CreatedAt: time.Now(),
		},
	}

	orderRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Order, error) {
		if id == orderID {
			return mockOrder, nil
		}
		return nil, errors.New("order not found")
	}

	vehicleRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Vehicle, error) {
		if id == vehicleID {
			return mockVehicle, nil
		}
		return nil, errors.New("vehicle not found")
	}

	orderInputRepoMock.FindByOrderIDFunc = func(id uuid.UUID) ([]models.OrderInput, error) {
		if id == orderID {
			return mockOrderInputs, nil
		}
		return nil, errors.New("order inputs not found")
	}

	inputRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Input, error) {
		if id == inputID1 {
			return mockInput1, nil
		}
		if id == inputID2 {
			return mockInput2, nil
		}
		return nil, errors.New("input not found")
	}

	orderStatusHistoryRepoMock.FindByOrderIDFunc = func(id uuid.UUID) ([]models.OrderStatusHistory, error) {
		if id == orderID {
			return mockStatusHistory, nil
		}
		return nil, errors.New("status history not found")
	}

	useCase := &FindOrderOverviewById{
		OrderRepository:              orderRepoMock,
		VehicleRepository:            vehicleRepoMock,
		OrderInputRepository:         orderInputRepoMock,
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		InputRepository:              inputRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	result, err := useCase.Process(orderID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected result, got nil")
	}

	if result.Order.ID != orderID {
		t.Errorf("Expected order ID %s, got %s", orderID, result.Order.ID)
	}

	if result.Vehicle.ID != vehicleID.String() {
		t.Errorf("Expected vehicle ID %s, got %s", vehicleID, result.Vehicle.ID)
	}

	if len(result.Inputs) != 2 {
		t.Errorf("Expected 2 inputs, got %d", len(result.Inputs))
	}

	if result.TotalPrice != 175.0 {
		t.Errorf("Expected total price 175.0, got %f", result.TotalPrice)
	}

	if len(result.Timeline) != 3 {
		t.Errorf("Expected 3 timeline entries, got %d", len(result.Timeline))
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing find completed order by ID",
		"Order found",
		"Vehicle found",
		"Order inputs found",
		"Added input detail",
		"Added input detail",
		"Calculated total price",
		"Found order status history",
		"Status duration calculated",
		"Status duration calculated",
		"Status not completed yet",
		"Timeline calculated",
		"Completed order with inputs and timeline retrieved successfully",
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

func TestFindOrderOverviewById_Process_OrderNotFound(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	orderInputRepoMock := &mocks.OrderInputRepositoryMock{}
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
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

	orderID := uuid.New()

	orderRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Order, error) {
		return nil, errors.New("order not found")
	}

	useCase := &FindOrderOverviewById{
		OrderRepository:              orderRepoMock,
		VehicleRepository:            vehicleRepoMock,
		OrderInputRepository:         orderInputRepoMock,
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		InputRepository:              inputRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	result, err := useCase.Process(orderID)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "order not found" {
		t.Errorf("Expected error 'order not found', got '%s'", err.Error())
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	// Verifica se os logs corretos foram chamados
	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Processing find completed order by ID" {
		t.Errorf("Expected log message 'Processing find completed order by ID', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Order not found" {
		t.Errorf("Expected error log 'Order not found', got '%s'", loggedErrors[0])
	}
}

func TestFindOrderOverviewById_Process_VehicleNotFound(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	orderInputRepoMock := &mocks.OrderInputRepositoryMock{}
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
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

	orderID := uuid.New()
	customerID := uuid.New()
	vehicleID := uuid.New()

	mockOrder := &models.Order{
		ID:         orderID,
		CustomerID: customerID,
		VehicleID:  vehicleID,
		Status:     "Completed",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	orderRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Order, error) {
		if id == orderID {
			return mockOrder, nil
		}
		return nil, errors.New("order not found")
	}

	vehicleRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Vehicle, error) {
		return nil, errors.New("vehicle not found")
	}

	useCase := &FindOrderOverviewById{
		OrderRepository:              orderRepoMock,
		VehicleRepository:            vehicleRepoMock,
		OrderInputRepository:         orderInputRepoMock,
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		InputRepository:              inputRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	result, err := useCase.Process(orderID)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "vehicle not found" {
		t.Errorf("Expected error 'vehicle not found', got '%s'", err.Error())
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing find completed order by ID",
		"Order found",
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

	if loggedErrors[0] != "Vehicle not found" {
		t.Errorf("Expected error log 'Vehicle not found', got '%s'", loggedErrors[0])
	}
}

func TestFindOrderOverviewById_Process_OrderInputsNotFound(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	orderInputRepoMock := &mocks.OrderInputRepositoryMock{}
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
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

	orderID := uuid.New()
	customerID := uuid.New()
	vehicleID := uuid.New()

	mockOrder := &models.Order{
		ID:         orderID,
		CustomerID: customerID,
		VehicleID:  vehicleID,
		Status:     "Completed",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	mockVehicle := &models.Vehicle{
		ID:                          vehicleID,
		CustomerID:                  customerID,
		NumberPlate:                 "ABC1234",
		Brand:                       "Toyota",
		Model:                       "Corolla",
		ReleaseYear:                 2020,
		VehicleIdentificationNumber: "VIN123456789",
		Color:                       "Prata",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	orderRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Order, error) {
		if id == orderID {
			return mockOrder, nil
		}
		return nil, errors.New("order not found")
	}

	vehicleRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Vehicle, error) {
		if id == vehicleID {
			return mockVehicle, nil
		}
		return nil, errors.New("vehicle not found")
	}

	orderInputRepoMock.FindByOrderIDFunc = func(id uuid.UUID) ([]models.OrderInput, error) {
		return nil, errors.New("database error")
	}

	useCase := &FindOrderOverviewById{
		OrderRepository:              orderRepoMock,
		VehicleRepository:            vehicleRepoMock,
		OrderInputRepository:         orderInputRepoMock,
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		InputRepository:              inputRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	result, err := useCase.Process(orderID)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "database error" {
		t.Errorf("Expected error 'database error', got '%s'", err.Error())
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing find completed order by ID",
		"Order found",
		"Vehicle found",
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

	if loggedErrors[0] != "Database error finding order inputs" {
		t.Errorf("Expected error log 'Database error finding order inputs', got '%s'", loggedErrors[0])
	}
}

func TestFindOrderOverviewById_FormatDurationFromSeconds(t *testing.T) {
	// Arrange
	testCases := []struct {
		name     string
		seconds  int
		expected string
	}{
		{"Zero seconds", 0, "00:00:00"},
		{"Negative seconds", -10, "00:00:00"},
		{"One hour", 3600, "01:00:00"},
		{"One hour and 30 minutes", 5400, "01:30:00"},
		{"Two hours, 15 minutes, 30 seconds", 8130, "02:15:30"},
		{"Complex time", 3661, "01:01:01"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := FormatDurationFromSeconds(tc.seconds)

			// Assert
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestFindOrderOverviewById_MapVehicleToDetails(t *testing.T) {
	// Arrange
	useCase := &FindOrderOverviewById{}
	vehicleID := uuid.New()

	mockVehicle := &models.Vehicle{
		ID:                          vehicleID,
		CustomerID:                  uuid.New(),
		NumberPlate:                 "ABC1234",
		Brand:                       "Toyota",
		Model:                       "Corolla",
		ReleaseYear:                 2020,
		VehicleIdentificationNumber: "VIN123456789",
		Color:                       "Prata",
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	// Act
	result := useCase.MapVehicleToDetails(mockVehicle)

	// Assert
	if result.ID != vehicleID.String() {
		t.Errorf("Expected vehicle ID %s, got %s", vehicleID, result.ID)
	}

	if result.Model != "Corolla" {
		t.Errorf("Expected model 'Corolla', got '%s'", result.Model)
	}

	if result.Brand != "Toyota" {
		t.Errorf("Expected brand 'Toyota', got '%s'", result.Brand)
	}

	if result.ReleaseYear != 2020 {
		t.Errorf("Expected release year 2020, got %d", result.ReleaseYear)
	}

	if result.VehicleIdentificationNumber != "VIN123456789" {
		t.Errorf("Expected VIN 'VIN123456789', got '%s'", result.VehicleIdentificationNumber)
	}

	if result.NumberPlate != "ABC1234" {
		t.Errorf("Expected number plate 'ABC1234', got '%s'", result.NumberPlate)
	}

	if result.Color != "Prata" {
		t.Errorf("Expected color 'Prata', got '%s'", result.Color)
	}
}

func TestFindOrderOverviewById_MapOrderInputToDetails(t *testing.T) {
	// Arrange
	useCase := &FindOrderOverviewById{}
	orderInputID := uuid.New()
	inputID := uuid.New()

	mockOrderInput := models.OrderInput{
		ID:         orderInputID,
		OrderID:    uuid.New(),
		InputID:    inputID,
		Quantity:   3,
		UnitPrice:  50.0,
		TotalPrice: 150.0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	inputName := "Óleo de Motor"

	// Act
	result := useCase.MapOrderInputToDetails(mockOrderInput, inputName)

	// Assert
	if result.ID != orderInputID.String() {
		t.Errorf("Expected order input ID %s, got %s", orderInputID, result.ID)
	}

	if result.InputID != inputID.String() {
		t.Errorf("Expected input ID %s, got %s", inputID, result.InputID)
	}

	if result.InputName != inputName {
		t.Errorf("Expected input name '%s', got '%s'", inputName, result.InputName)
	}

	if result.Quantity != 3 {
		t.Errorf("Expected quantity 3, got %d", result.Quantity)
	}

	if result.UnitPrice != 50.0 {
		t.Errorf("Expected unit price 50.0, got %f", result.UnitPrice)
	}

	if result.TotalPrice != 150.0 {
		t.Errorf("Expected total price 150.0, got %f", result.TotalPrice)
	}
}

func TestFindOrderOverviewById_MapOrderToDomain(t *testing.T) {
	// Arrange
	useCase := &FindOrderOverviewById{}
	orderID := uuid.New()
	customerID := uuid.New()
	vehicleID := uuid.New()

	mockOrder := &models.Order{
		ID:         orderID,
		CustomerID: customerID,
		VehicleID:  vehicleID,
		Status:     "Completed",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Act
	result := useCase.MapOrderToDomain(mockOrder)

	// Assert
	if result.ID != orderID {
		t.Errorf("Expected order ID %s, got %s", orderID, result.ID)
	}

	if result.CustomerID != customerID {
		t.Errorf("Expected customer ID %s, got %s", customerID, result.CustomerID)
	}

	if result.VehicleID != vehicleID {
		t.Errorf("Expected vehicle ID %s, got %s", vehicleID, result.VehicleID)
	}

	if result.Status != "Completed" {
		t.Errorf("Expected status 'Completed', got '%s'", result.Status)
	}
}

func TestFindOrderOverviewById_CalculateTimeline_WithHistory(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	orderInputRepoMock := &mocks.OrderInputRepositoryMock{}
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
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

	orderID := uuid.New()
	now := time.Now()

	mockStatusHistory := []models.OrderStatusHistory{
		{
			ID:        uuid.New(),
			OrderID:   orderID,
			Status:    "Received",
			StartedAt: now.Add(-2 * time.Hour),
			EndedAt:   &[]time.Time{now.Add(-1 * time.Hour)}[0],
			CreatedAt: now,
		},
		{
			ID:        uuid.New(),
			OrderID:   orderID,
			Status:    "In Progress",
			StartedAt: now.Add(-1 * time.Hour),
			EndedAt:   &[]time.Time{now.Add(-30 * time.Minute)}[0],
			CreatedAt: now,
		},
		{
			ID:        uuid.New(),
			OrderID:   orderID,
			Status:    "Completed",
			StartedAt: now.Add(-30 * time.Minute),
			EndedAt:   nil, // Status atual
			CreatedAt: now,
		},
	}

	orderStatusHistoryRepoMock.FindByOrderIDFunc = func(id uuid.UUID) ([]models.OrderStatusHistory, error) {
		if id == orderID {
			return mockStatusHistory, nil
		}
		return nil, errors.New("status history not found")
	}

	useCase := &FindOrderOverviewById{
		OrderRepository:              orderRepoMock,
		VehicleRepository:            vehicleRepoMock,
		OrderInputRepository:         orderInputRepoMock,
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		InputRepository:              inputRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	timeline, averageTime := useCase.CalculateTimeline(orderID)

	// Assert
	if len(timeline) != 3 {
		t.Errorf("Expected 3 timeline entries, got %d", len(timeline))
	}

	if timeline["Received"] == "00:00:00" {
		t.Error("Expected non-zero duration for 'Received' status")
	}

	if timeline["In Progress"] == "00:00:00" {
		t.Error("Expected non-zero duration for 'In Progress' status")
	}

	if timeline["Completed"] != "00:00:00" {
		t.Error("Expected zero duration for current 'Completed' status")
	}

	if averageTime == "00:00:00" {
		t.Error("Expected non-zero average time")
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Found order status history",
		"Status duration calculated",
		"Status duration calculated",
		"Status not completed yet",
		"Timeline calculated",
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

func TestFindOrderOverviewById_CalculateTimeline_NoHistory(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	orderInputRepoMock := &mocks.OrderInputRepositoryMock{}
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
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

	orderID := uuid.New()

	orderStatusHistoryRepoMock.FindByOrderIDFunc = func(id uuid.UUID) ([]models.OrderStatusHistory, error) {
		return []models.OrderStatusHistory{}, nil
	}

	useCase := &FindOrderOverviewById{
		OrderRepository:              orderRepoMock,
		VehicleRepository:            vehicleRepoMock,
		OrderInputRepository:         orderInputRepoMock,
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		InputRepository:              inputRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	timeline, averageTime := useCase.CalculateTimeline(orderID)

	// Assert
	if len(timeline) != 0 {
		t.Errorf("Expected 0 timeline entries, got %d", len(timeline))
	}

	if averageTime != "00:00:00" {
		t.Errorf("Expected average time '00:00:00', got '%s'", averageTime)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Found order status history",
		"Timeline calculated",
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

func TestFindOrderOverviewById_CalculateTimeline_Error(t *testing.T) {
	// Arrange
	orderRepoMock := &mocks.OrderRepositoryMock{}
	vehicleRepoMock := &mocks.VehicleRepositoryMock{}
	orderInputRepoMock := &mocks.OrderInputRepositoryMock{}
	orderStatusHistoryRepoMock := &mocks.OrderStatusHistoryRepositoryMock{}
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

	orderID := uuid.New()

	orderStatusHistoryRepoMock.FindByOrderIDFunc = func(id uuid.UUID) ([]models.OrderStatusHistory, error) {
		return nil, errors.New("database error")
	}

	useCase := &FindOrderOverviewById{
		OrderRepository:              orderRepoMock,
		VehicleRepository:            vehicleRepoMock,
		OrderInputRepository:         orderInputRepoMock,
		OrderStatusHistoryRepository: orderStatusHistoryRepoMock,
		InputRepository:              inputRepoMock,
		Logger:                       loggerMock,
	}

	// Act
	timeline, averageTime := useCase.CalculateTimeline(orderID)

	// Assert
	if len(timeline) != 0 {
		t.Errorf("Expected empty timeline, got %d entries", len(timeline))
	}

	if averageTime != "00:00:00" {
		t.Errorf("Expected average time '00:00:00', got '%s'", averageTime)
	}

	// Verifica se os logs corretos foram chamados
	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Error fetching order status history" {
		t.Errorf("Expected error log 'Error fetching order status history', got '%s'", loggedErrors[0])
	}
}
