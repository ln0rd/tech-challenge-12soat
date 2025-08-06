package customer

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/costumer"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/test/mocks"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestUpdateByIdCustomer_Process_Success(t *testing.T) {
	// Arrange
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	customerID := uuid.New()
	existingCustomer := &models.Customer{
		ID:             customerID,
		Name:           "João Silva",
		DocumentNumber: "12345678901",
		CustomerType:   "individual",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	updateEntity := &domain.Customer{
		ID:             customerID,
		Name:           "João Silva Santos",
		DocumentNumber: "12345678902",
		CustomerType:   "individual",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		if id == customerID {
			return existingCustomer, nil
		}
		return nil, gorm.ErrRecordNotFound
	}

	customerRepoMock.UpdateFunc = func(customer *models.Customer) error {
		return nil
	}

	useCase := &UpdateByIdCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.Process(customerID, updateEntity)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing update customer by ID",
		"Successfully fetched customer from database",
		"Customer fields updated",
		"Customer updated in database",
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

func TestUpdateByIdCustomer_Process_CustomerNotFound(t *testing.T) {
	// Arrange
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	customerID := uuid.New()
	updateEntity := &domain.Customer{
		ID:             customerID,
		Name:           "João Silva Santos",
		DocumentNumber: "12345678902",
		CustomerType:   "individual",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		return nil, gorm.ErrRecordNotFound
	}

	useCase := &UpdateByIdCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.Process(customerID, updateEntity)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err != gorm.ErrRecordNotFound {
		t.Errorf("Expected error %v, got %v", gorm.ErrRecordNotFound, err)
	}

	// Verifica se os logs corretos foram chamados
	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Processing update customer by ID" {
		t.Errorf("Expected log message 'Processing update customer by ID', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Customer not found" {
		t.Errorf("Expected error log 'Customer not found', got '%s'", loggedErrors[0])
	}
}

func TestUpdateByIdCustomer_Process_DatabaseError(t *testing.T) {
	// Arrange
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	customerID := uuid.New()
	updateEntity := &domain.Customer{
		ID:             customerID,
		Name:           "João Silva Santos",
		DocumentNumber: "12345678902",
		CustomerType:   "individual",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	expectedError := errors.New("database connection failed")
	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		return nil, expectedError
	}

	useCase := &UpdateByIdCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.Process(customerID, updateEntity)

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

	if loggedInfo[0] != "Processing update customer by ID" {
		t.Errorf("Expected log message 'Processing update customer by ID', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(loggedErrors))
	}

	if loggedErrors[0] != "Database error fetching customer" {
		t.Errorf("Expected error log 'Database error fetching customer', got '%s'", loggedErrors[0])
	}
}

func TestUpdateByIdCustomer_Process_UpdateError(t *testing.T) {
	// Arrange
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	customerID := uuid.New()
	existingCustomer := &models.Customer{
		ID:             customerID,
		Name:           "João Silva",
		DocumentNumber: "12345678901",
		CustomerType:   "individual",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	updateEntity := &domain.Customer{
		ID:             customerID,
		Name:           "João Silva Santos",
		DocumentNumber: "12345678902",
		CustomerType:   "individual",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	expectedError := errors.New("update constraint violation")
	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		if id == customerID {
			return existingCustomer, nil
		}
		return nil, gorm.ErrRecordNotFound
	}

	customerRepoMock.UpdateFunc = func(customer *models.Customer) error {
		return expectedError
	}

	useCase := &UpdateByIdCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.Process(customerID, updateEntity)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	// Verifica se os logs corretos foram chamados
	expectedInfoLogs := []string{
		"Processing update customer by ID",
		"Successfully fetched customer from database",
		"Customer fields updated",
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

	if loggedErrors[0] != "Database error updating customer" {
		t.Errorf("Expected error log 'Database error updating customer', got '%s'", loggedErrors[0])
	}
}

func TestUpdateByIdCustomer_FetchCustomerFromDB_Success(t *testing.T) {
	// Arrange
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	customerID := uuid.New()
	mockCustomer := &models.Customer{
		ID:             customerID,
		Name:           "João Silva",
		DocumentNumber: "12345678901",
		CustomerType:   "individual",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	customerRepoMock.FindByIDFunc = func(id uuid.UUID) (*models.Customer, error) {
		if id == customerID {
			return mockCustomer, nil
		}
		return nil, gorm.ErrRecordNotFound
	}

	useCase := &UpdateByIdCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	result, err := useCase.FetchCustomerFromDB(customerID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected customer result, got nil")
	}

	if result.ID != customerID {
		t.Errorf("Expected customer ID %s, got %s", customerID, result.ID)
	}

	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Successfully fetched customer from database" {
		t.Errorf("Expected log message 'Successfully fetched customer from database', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) > 0 {
		t.Errorf("Expected no error logs, got %d", len(loggedErrors))
	}
}

func TestUpdateByIdCustomer_UpdateCustomerFields(t *testing.T) {
	// Arrange
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	existingCustomer := &models.Customer{
		ID:             uuid.New(),
		Name:           "João Silva",
		DocumentNumber: "12345678901",
		CustomerType:   "individual",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	updateEntity := &domain.Customer{
		ID:             existingCustomer.ID,
		Name:           "João Silva Santos",
		DocumentNumber: "12345678902",
		CustomerType:   "corporate",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	useCase := &UpdateByIdCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	useCase.UpdateCustomerFields(existingCustomer, updateEntity)

	// Assert
	if existingCustomer.Name != updateEntity.Name {
		t.Errorf("Expected name %s, got %s", updateEntity.Name, existingCustomer.Name)
	}

	if existingCustomer.DocumentNumber != updateEntity.DocumentNumber {
		t.Errorf("Expected document number %s, got %s", updateEntity.DocumentNumber, existingCustomer.DocumentNumber)
	}

	if existingCustomer.CustomerType != updateEntity.CustomerType {
		t.Errorf("Expected customer type %s, got %s", updateEntity.CustomerType, existingCustomer.CustomerType)
	}

	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Customer fields updated" {
		t.Errorf("Expected log message 'Customer fields updated', got '%s'", loggedInfo[0])
	}
}

func TestUpdateByIdCustomer_SaveCustomerToDB_Success(t *testing.T) {
	// Arrange
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	customer := &models.Customer{
		ID:             uuid.New(),
		Name:           "João Silva Santos",
		DocumentNumber: "12345678902",
		CustomerType:   "individual",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	customerRepoMock.UpdateFunc = func(customer *models.Customer) error {
		return nil
	}

	useCase := &UpdateByIdCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.SaveCustomerToDB(customer)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(loggedInfo) != 1 {
		t.Errorf("Expected 1 info log, got %d", len(loggedInfo))
	}

	if loggedInfo[0] != "Customer updated in database" {
		t.Errorf("Expected log message 'Customer updated in database', got '%s'", loggedInfo[0])
	}

	if len(loggedErrors) > 0 {
		t.Errorf("Expected no error logs, got %d", len(loggedErrors))
	}
}

func TestUpdateByIdCustomer_SaveCustomerToDB_Error(t *testing.T) {
	// Arrange
	customerRepoMock := &mocks.CustomerRepositoryMock{}
	loggerMock := &mocks.LoggerMock{}

	var loggedInfo []string
	var loggedErrors []string

	loggerMock.InfoFunc = func(msg string, fields ...zap.Field) {
		loggedInfo = append(loggedInfo, msg)
	}

	loggerMock.ErrorFunc = func(msg string, fields ...zap.Field) {
		loggedErrors = append(loggedErrors, msg)
	}

	customer := &models.Customer{
		ID:             uuid.New(),
		Name:           "João Silva Santos",
		DocumentNumber: "12345678902",
		CustomerType:   "individual",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	expectedError := errors.New("update constraint violation")
	customerRepoMock.UpdateFunc = func(customer *models.Customer) error {
		return expectedError
	}

	useCase := &UpdateByIdCustomer{
		CustomerRepository: customerRepoMock,
		Logger:             loggerMock,
	}

	// Act
	err := useCase.SaveCustomerToDB(customer)

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

	if loggedErrors[0] != "Database error updating customer" {
		t.Errorf("Expected error log 'Database error updating customer', got '%s'", loggedErrors[0])
	}
}
