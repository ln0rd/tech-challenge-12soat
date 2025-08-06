package order

import (
	"errors"

	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/order"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db/models"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/logger"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/repository"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/persistence"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/order_status_history"
	"go.uber.org/zap"
)

type CreateOrder struct {
	OrderRepository      repository.OrderRepository
	CustomerRepository   repository.CustomerRepository
	VehicleRepository    repository.VehicleRepository
	Logger               logger.Logger
	StatusHistoryManager *order_status_history.ManageOrderStatusHistory
}

// FetchCustomerFromDB busca um customer específico do banco de dados
func (uc *CreateOrder) FetchCustomerFromDB(customerID uuid.UUID) (*models.Customer, error) {
	customer, err := uc.CustomerRepository.FindByID(customerID)
	if err != nil {
		uc.Logger.Error("Customer not found", zap.String("customerID", customerID.String()))
		return nil, errors.New("customer not found")
	}
	uc.Logger.Info("Customer found", zap.String("customerID", customerID.String()))
	return customer, nil
}

// FetchVehicleFromDB busca um vehicle específico do banco de dados
func (uc *CreateOrder) FetchVehicleFromDB(vehicleID uuid.UUID) (*models.Vehicle, error) {
	vehicle, err := uc.VehicleRepository.FindByID(vehicleID)
	if err != nil {
		uc.Logger.Error("Vehicle not found", zap.String("vehicleID", vehicleID.String()))
		return nil, errors.New("vehicle not found")
	}
	uc.Logger.Info("Vehicle found", zap.String("vehicleID", vehicleID.String()))
	return vehicle, nil
}

// ValidateVehicleOwnership valida se o vehicle pertence ao customer
func (uc *CreateOrder) ValidateVehicleOwnership(vehicle *models.Vehicle, customerID uuid.UUID) error {
	if vehicle.CustomerID != customerID {
		uc.Logger.Error("Vehicle does not belong to customer",
			zap.String("vehicleID", vehicle.ID.String()),
			zap.String("customerID", customerID.String()))
		return errors.New("vehicle does not belong to customer")
	}
	uc.Logger.Info("Vehicle belongs to customer")
	return nil
}

// SaveOrderToDB salva o order no banco de dados
func (uc *CreateOrder) SaveOrderToDB(model *models.Order) error {
	err := uc.OrderRepository.Create(model)
	if err != nil {
		uc.Logger.Error("Database error creating order", zap.Error(err))
		return err
	}

	uc.Logger.Info("Order created in database", zap.String("id", model.ID.String()))
	return nil
}

// StartOrderStatusHistory inicia o histórico de status da order
func (uc *CreateOrder) StartOrderStatusHistory(orderID uuid.UUID, status string) error {
	err := uc.StatusHistoryManager.StartNewStatus(orderID, status)
	if err != nil {
		uc.Logger.Error("Error starting status history", zap.Error(err))
		// Não retorna erro aqui, pois a order já foi criada
		// Apenas loga o erro para monitoramento
		return nil
	}

	uc.Logger.Info("Status history started successfully",
		zap.String("orderID", orderID.String()),
		zap.String("status", status))
	return nil
}

func (uc *CreateOrder) Process(entity *domain.Order) error {
	uc.Logger.Info("Processing order creation",
		zap.String("customerID", entity.CustomerID.String()),
		zap.String("vehicleID", entity.VehicleID.String()),
		zap.String("status", entity.Status))

	// Busca e valida customer
	_, err := uc.FetchCustomerFromDB(entity.CustomerID)
	if err != nil {
		return err
	}

	// Busca e valida vehicle
	vehicle, err := uc.FetchVehicleFromDB(entity.VehicleID)
	if err != nil {
		return err
	}

	// Valida se o vehicle pertence ao customer
	if err := uc.ValidateVehicleOwnership(vehicle, entity.CustomerID); err != nil {
		return err
	}

	// Mapeia entidade para modelo usando persistence
	model := persistence.OrderPersistence{}.ToModel(entity)
	uc.Logger.Info("Model created",
		zap.String("id", model.ID.String()),
		zap.String("customerID", model.CustomerID.String()),
		zap.String("vehicleID", model.VehicleID.String()),
		zap.String("status", model.Status))

	// Salva no banco
	err = uc.SaveOrderToDB(model)
	if err != nil {
		return err
	}

	// Inicia o histórico de status com o status inicial
	uc.StartOrderStatusHistory(model.ID, model.Status)

	return nil
}
