package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/order"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/order"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/order_input"
	"go.uber.org/zap"
)

var (
	orderStatusRegex = regexp.MustCompile(`^(Received|Undergoing diagnosis|Awaiting approval|In progress|Completed|Delivered)$`)
)

const (
	OrderStatusReceived            = "Received"
	OrderStatusUndergoingDiagnosis = "Undergoing diagnosis"
	OrderStatusAwaitingApproval    = "Awaiting approval"
	OrderStatusInProgress          = "In progress"
	OrderStatusCompleted           = "Completed"
	OrderStatusDelivered           = "Delivered"
)

type OrderController struct {
	Logger                   *zap.Logger
	CreateOrder              *order.CreateOrder
	AddInputToOrderUC        *order_input.AddInputToOrder
	RemoveInputFromOrderUC   *order_input.RemoveInputFromOrder
	FindCompletedOrderByIdUC *order.FindCompletedOrderById
	UpdateOrderStatusUC      *order.UpdateOrderStatus
}

type OrderDTO struct {
	CustomerID string `json:"customer_id"`
	VehicleID  string `json:"vehicle_id"`
}

func (dto *OrderDTO) Validate() error {
	if dto.CustomerID == "" {
		return errors.New("customer_id is required")
	}
	if dto.VehicleID == "" {
		return errors.New("vehicle_id is required")
	}
	return nil
}

type AddInputToOrderDTO struct {
	InputID  string `json:"input_id"`
	Quantity int    `json:"quantity"`
}

func (dto *AddInputToOrderDTO) Validate() error {
	if dto.InputID == "" {
		return errors.New("input_id is required")
	}
	if dto.Quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	return nil
}

type RemoveInputFromOrderDTO struct {
	InputID  string `json:"input_id"`
	Quantity int    `json:"quantity"`
}

func (dto *RemoveInputFromOrderDTO) Validate() error {
	if dto.InputID == "" {
		return errors.New("input_id is required")
	}
	if dto.Quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	return nil
}

type UpdateOrderStatusDTO struct {
	Status string `json:"status"`
}

func (dto *UpdateOrderStatusDTO) Validate() error {
	if dto.Status == "" {
		return errors.New("status is required")
	}
	return nil
}

func (oc *OrderController) Create(w http.ResponseWriter, r *http.Request) {
	oc.Logger.Info("=== ORDER CREATE ENDPOINT CALLED ===")

	var dto OrderDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		oc.Logger.Error("Error decoding JSON", zap.Error(err))
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	oc.Logger.Info("Received order creation request",
		zap.String("customerID", dto.CustomerID),
		zap.String("vehicleID", dto.VehicleID))

	if err := dto.Validate(); err != nil {
		oc.Logger.Error("Validation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	oc.Logger.Info("Validation passed")

	// Parse customer ID
	customerID, err := uuid.Parse(dto.CustomerID)
	if err != nil {
		oc.Logger.Error("Error parsing customer ID", zap.Error(err))
		http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
		return
	}
	oc.Logger.Info("Customer ID parsed successfully", zap.String("customerID", customerID.String()))

	// Parse vehicle ID
	vehicleID, err := uuid.Parse(dto.VehicleID)
	if err != nil {
		oc.Logger.Error("Error parsing vehicle ID", zap.Error(err))
		http.Error(w, "Invalid vehicle ID format", http.StatusBadRequest)
		return
	}
	oc.Logger.Info("Vehicle ID parsed successfully", zap.String("vehicleID", vehicleID.String()))

	entity := &domain.Order{
		ID:         uuid.New(),
		CustomerID: customerID,
		VehicleID:  vehicleID,
		Status:     OrderStatusReceived, // Status inicial automático
	}

	oc.Logger.Info("Entity created",
		zap.String("id", entity.ID.String()),
		zap.String("customerID", entity.CustomerID.String()),
		zap.String("vehicleID", entity.VehicleID.String()),
		zap.String("status", entity.Status))

	oc.Logger.Info("Calling CreateOrder.Process...")
	err = oc.CreateOrder.Process(entity)
	if err != nil {
		oc.Logger.Error("Error creating order", zap.Error(err))

		// Tratamento específico para customer não encontrado
		if err.Error() == "customer not found" {
			http.Error(w, "Customer not found", http.StatusNotFound)
			return
		}

		// Tratamento específico para vehicle não encontrado
		if err.Error() == "vehicle not found" {
			http.Error(w, "Vehicle not found", http.StatusNotFound)
			return
		}

		// Tratamento específico para vehicle não pertencer ao customer
		if err.Error() == "vehicle does not belong to customer" {
			http.Error(w, "Vehicle does not belong to customer", http.StatusBadRequest)
			return
		}

		http.Error(w, "Error creating order", http.StatusInternalServerError)
		return
	}

	oc.Logger.Info("Order created successfully",
		zap.String("id", entity.ID.String()),
		zap.String("status", entity.Status))

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Order created successfully",
		"id":      entity.ID.String(),
		"status":  entity.Status,
	})
}

func (oc *OrderController) AddInputToOrder(w http.ResponseWriter, r *http.Request) {
	// Extrai o order ID da URL
	vars := mux.Vars(r)
	orderIDStr := vars["orderId"]

	oc.Logger.Info("Received add input to order request", zap.String("orderID", orderIDStr))

	// Parse order ID
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		oc.Logger.Error("Error parsing order ID", zap.Error(err))
		http.Error(w, "Invalid order ID format", http.StatusBadRequest)
		return
	}
	oc.Logger.Info("Order ID parsed successfully", zap.String("orderID", orderID.String()))

	var dto AddInputToOrderDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		oc.Logger.Error("Error decoding JSON", zap.Error(err))
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	oc.Logger.Info("Received add input request",
		zap.String("inputID", dto.InputID),
		zap.Int("quantity", dto.Quantity))

	if err := dto.Validate(); err != nil {
		oc.Logger.Error("Validation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	oc.Logger.Info("Validation passed")

	// Parse input ID
	inputID, err := uuid.Parse(dto.InputID)
	if err != nil {
		oc.Logger.Error("Error parsing input ID", zap.Error(err))
		http.Error(w, "Invalid input ID format", http.StatusBadRequest)
		return
	}
	oc.Logger.Info("Input ID parsed successfully", zap.String("inputID", inputID.String()))

	oc.Logger.Info("Calling AddInputToOrder.Process...")
	err = oc.AddInputToOrderUC.Process(orderID, inputID, dto.Quantity)
	if err != nil {
		oc.Logger.Error("Error adding input to order", zap.Error(err))

		// Tratamento específico para order não encontrado
		if err.Error() == "order not found" {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		// Tratamento específico para input não encontrado
		if err.Error() == "input not found" {
			http.Error(w, "Input not found", http.StatusNotFound)
			return
		}

		// Tratamento específico para quantidade insuficiente
		if err.Error() == "insufficient input quantity" {
			http.Error(w, "Insufficient input quantity", http.StatusBadRequest)
			return
		}

		// Tratamento específico para quantidade inválida
		if err.Error() == "quantity must be greater than zero" {
			http.Error(w, "Quantity must be greater than zero", http.StatusBadRequest)
			return
		}

		http.Error(w, "Error adding input to order", http.StatusInternalServerError)
		return
	}

	oc.Logger.Info("Input added to order successfully",
		zap.String("orderID", orderID.String()),
		zap.String("inputID", inputID.String()),
		zap.Int("quantity", dto.Quantity))

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Input added to order successfully",
		"order_id": orderID.String(),
		"input_id": inputID.String(),
	})
}

func (oc *OrderController) RemoveInputFromOrder(w http.ResponseWriter, r *http.Request) {
	oc.Logger.Info("=== ORDER REMOVE INPUT ENDPOINT CALLED ===")

	// Extrai o order ID da URL
	vars := mux.Vars(r)
	orderIDStr := vars["orderId"]

	oc.Logger.Info("Received remove input from order request",
		zap.String("orderID", orderIDStr))

	// Parse order ID
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		oc.Logger.Error("Error parsing order ID", zap.Error(err))
		http.Error(w, "Invalid order ID format", http.StatusBadRequest)
		return
	}
	oc.Logger.Info("Order ID parsed successfully", zap.String("orderID", orderID.String()))

	// Decodifica o body para obter input_id e quantidade
	var dto RemoveInputFromOrderDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		oc.Logger.Error("Error decoding JSON", zap.Error(err))
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	oc.Logger.Info("Received remove input request",
		zap.String("inputID", dto.InputID),
		zap.Int("quantity", dto.Quantity))

	if err := dto.Validate(); err != nil {
		oc.Logger.Error("Validation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	oc.Logger.Info("Validation passed")

	// Parse input ID
	inputID, err := uuid.Parse(dto.InputID)
	if err != nil {
		oc.Logger.Error("Error parsing input ID", zap.Error(err))
		http.Error(w, "Invalid input ID format", http.StatusBadRequest)
		return
	}
	oc.Logger.Info("Input ID parsed successfully", zap.String("inputID", inputID.String()))

	oc.Logger.Info("Calling RemoveInputFromOrder.Process...")
	err = oc.RemoveInputFromOrderUC.Process(orderID, inputID, dto.Quantity)
	if err != nil {
		oc.Logger.Error("Error removing input from order", zap.Error(err))

		// Tratamento específico para order não encontrado
		if err.Error() == "order not found" {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		// Tratamento específico para input não encontrado
		if err.Error() == "input not found" {
			http.Error(w, "Input not found", http.StatusNotFound)
			return
		}

		// Tratamento específico para order input não encontrado
		if err.Error() == "order input not found" {
			http.Error(w, "Order input not found", http.StatusNotFound)
			return
		}

		// Tratamento específico para quantidade inválida no order input
		if err.Error() == "invalid quantity in order input" {
			http.Error(w, "Invalid quantity in order input", http.StatusBadRequest)
			return
		}

		// Tratamento específico para quantidade insuficiente no order input
		if err.Error() == "insufficient quantity in order input" {
			http.Error(w, "Insufficient quantity in order input", http.StatusBadRequest)
			return
		}

		// Tratamento específico para quantidade inválida a remover
		if err.Error() == "quantity to remove must be greater than zero" {
			http.Error(w, "Quantity to remove must be greater than zero", http.StatusBadRequest)
			return
		}

		http.Error(w, "Error removing input from order", http.StatusInternalServerError)
		return
	}

	oc.Logger.Info("Input removed from order successfully",
		zap.String("orderID", orderID.String()),
		zap.String("inputID", inputID.String()))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Input removed from order successfully",
		"order_id": orderID.String(),
		"input_id": inputID.String(),
	})
}

func (oc *OrderController) FindCompletedOrderById(w http.ResponseWriter, r *http.Request) {
	// Extrai o order ID da URL
	vars := mux.Vars(r)
	orderIDStr := vars["orderId"]

	oc.Logger.Info("Received find completed order by ID request", zap.String("orderID", orderIDStr))

	// Parse order ID
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		oc.Logger.Error("Error parsing order ID", zap.Error(err))
		http.Error(w, "Invalid order ID format", http.StatusBadRequest)
		return
	}
	oc.Logger.Info("Order ID parsed successfully", zap.String("orderID", orderID.String()))

	oc.Logger.Info("Calling FindCompletedOrderById.Process...")
	result, err := oc.FindCompletedOrderByIdUC.Process(orderID)
	if err != nil {
		oc.Logger.Error("Error finding completed order by ID", zap.Error(err))

		// Tratamento específico para order não encontrado
		if err.Error() == "order not found" {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		// Tratamento específico para vehicle não encontrado
		if err.Error() == "vehicle not found" {
			http.Error(w, "Vehicle not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Error finding completed order by ID", http.StatusInternalServerError)
		return
	}

	oc.Logger.Info("Completed order with inputs found successfully",
		zap.String("orderID", orderID.String()),
		zap.Int("inputsCount", len(result.Inputs)))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (oc *OrderController) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	oc.Logger.Info("=== ORDER UPDATE STATUS ENDPOINT CALLED ===")

	// Extrai o order ID da URL
	vars := mux.Vars(r)
	orderIDStr := vars["orderId"]

	oc.Logger.Info("Received update order status request", zap.String("orderID", orderIDStr))

	// Parse order ID
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		oc.Logger.Error("Error parsing order ID", zap.Error(err))
		http.Error(w, "Invalid order ID format", http.StatusBadRequest)
		return
	}
	oc.Logger.Info("Order ID parsed successfully", zap.String("orderID", orderID.String()))

	// Decodifica o body para obter o novo status
	var dto UpdateOrderStatusDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		oc.Logger.Error("Error decoding JSON", zap.Error(err))
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	oc.Logger.Info("Received update status request", zap.String("newStatus", dto.Status))

	if err := dto.Validate(); err != nil {
		oc.Logger.Error("Validation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	oc.Logger.Info("Validation passed")

	oc.Logger.Info("Calling UpdateOrderStatus.Process...")
	err = oc.UpdateOrderStatusUC.Process(orderID, dto.Status)
	if err != nil {
		oc.Logger.Error("Error updating order status", zap.Error(err))

		// Tratamento específico para order não encontrado
		if err.Error() == "order not found" {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		// Tratamento específico para status inválido
		if err.Error() == "invalid order status" {
			http.Error(w, "Invalid order status", http.StatusBadRequest)
			return
		}

		http.Error(w, "Error updating order status", http.StatusInternalServerError)
		return
	}

	oc.Logger.Info("Order status updated successfully",
		zap.String("orderID", orderID.String()),
		zap.String("newStatus", dto.Status))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Order status updated successfully",
		"order_id": orderID.String(),
		"status":   dto.Status,
	})
}
