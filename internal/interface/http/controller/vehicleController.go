package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/vehicle"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/vehicle"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	modelRegex = regexp.MustCompile(`^[a-zA-Z0-9\s\-]{2,50}$`)
	brandRegex = regexp.MustCompile(`^[a-zA-Z0-9\s\-]{2,30}$`)
	vinRegex   = regexp.MustCompile(`^[A-HJ-NPR-Z0-9]{17}$`)
	plateRegex = regexp.MustCompile(`^[A-Z]{3}[0-9][0-9A-Z][0-9]{2}$`)
	colorRegex = regexp.MustCompile(`^[a-zA-Z\s]{2,20}$`)
)

type VehicleController struct {
	Logger                  *zap.Logger
	CreateVehicle           *vehicle.CreateVehicle
	FindByIdVehicle         *vehicle.FindByIdVehicle
	FindByCustomerIdVehicle *vehicle.FindByCustomerIdVehicle
	UpdateByIdVehicle       *vehicle.UpdateByIdVehicle
	DeleteByIdVehicle       *vehicle.DeleteByIdVehicle
}

type VehicleDTO struct {
	Model                       string `json:"model"`
	Brand                       string `json:"brand"`
	ReleaseYear                 int    `json:"release_year"`
	VehicleIdentificationNumber string `json:"vehicle_identification_number"`
	NumberPlate                 string `json:"number_plate"`
	Color                       string `json:"color"`
	CustomerID                  string `json:"customer_id"`
}

func (dto *VehicleDTO) Validate() error {
	if !modelRegex.MatchString(dto.Model) {
		return errors.New("model must contain only letters, numbers, spaces and hyphens, between 2 and 50 characters")
	}
	if !brandRegex.MatchString(dto.Brand) {
		return errors.New("brand must contain only letters, numbers, spaces and hyphens, between 2 and 30 characters")
	}
	if !vinRegex.MatchString(dto.VehicleIdentificationNumber) {
		return errors.New("vehicle identification number must be exactly 17 characters (A-Z, 0-9, excluding I, O, Q)")
	}
	if !plateRegex.MatchString(dto.NumberPlate) {
		return errors.New("number plate must follow Brazilian format: ABC1D23")
	}
	if !colorRegex.MatchString(dto.Color) {
		return errors.New("color must contain only letters and spaces, between 2 and 20 characters")
	}
	if dto.ReleaseYear < 1900 || dto.ReleaseYear > 2024 {
		return errors.New("release year must be between 1900 and 2024")
	}
	if dto.CustomerID == "" {
		return errors.New("customer_id is required")
	}
	return nil
}

func (vc *VehicleController) Create(w http.ResponseWriter, r *http.Request) {
	vc.Logger.Info("=== VEHICLE CREATE ENDPOINT CALLED ===")

	var dto VehicleDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		vc.Logger.Error("Error decoding JSON", zap.Error(err))
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	vc.Logger.Info("Received vehicle creation request",
		zap.String("model", dto.Model),
		zap.String("brand", dto.Brand),
		zap.String("numberPlate", dto.NumberPlate),
		zap.Int("releaseYear", dto.ReleaseYear),
		zap.Any("customerID", dto.CustomerID))

	if err := dto.Validate(); err != nil {
		vc.Logger.Error("Validation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vc.Logger.Info("Validation passed")

	parsedCustomerID, err := uuid.Parse(dto.CustomerID)
	if err != nil {
		vc.Logger.Error("Error parsing customer ID", zap.Error(err))
		http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
		return
	}
	customerID := parsedCustomerID
	vc.Logger.Info("Customer ID parsed successfully", zap.String("customerID", customerID.String()))

	entity := &domain.Vehicle{
		ID:                          uuid.New(),
		Model:                       dto.Model,
		Brand:                       dto.Brand,
		ReleaseYear:                 dto.ReleaseYear,
		VehicleIdentificationNumber: dto.VehicleIdentificationNumber,
		NumberPlate:                 dto.NumberPlate,
		Color:                       dto.Color,
		CustomerID:                  customerID,
	}

	vc.Logger.Info("Entity created",
		zap.String("id", entity.ID.String()),
		zap.String("model", entity.Model),
		zap.String("brand", entity.Brand),
		zap.String("numberPlate", entity.NumberPlate))

	vc.Logger.Info("Calling CreateVehicle.Process...")
	err = vc.CreateVehicle.Process(entity)
	if err != nil {
		vc.Logger.Error("Error creating vehicle", zap.Error(err))

		// Tratamento específico para placa duplicada
		if err.Error() == "number plate already exists" {
			http.Error(w, "Number plate already exists", http.StatusConflict)
			return
		}

		// Tratamento específico para customer não encontrado
		if err.Error() == "customer not found" {
			http.Error(w, "Customer not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Error creating vehicle", http.StatusInternalServerError)
		return
	}

	vc.Logger.Info("Vehicle created successfully",
		zap.String("id", entity.ID.String()),
		zap.String("numberPlate", entity.NumberPlate))

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Vehicle created successfully",
		"id":      entity.ID.String(),
	})
}

func (vc *VehicleController) FindById(w http.ResponseWriter, r *http.Request) {
	vc.Logger.Info("=== VEHICLE FIND BY ID ENDPOINT CALLED ===")

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		vc.Logger.Error("Error parsing UUID", zap.Error(err))
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	vc.Logger.Info("Parsed vehicle ID", zap.String("id", id.String()))

	vc.Logger.Info("Calling FindByIdVehicle.Process...")
	vehicle, err := vc.FindByIdVehicle.Process(id)
	if err != nil {
		vc.Logger.Error("Error finding vehicle by ID", zap.Error(err), zap.String("id", id.String()))
		http.Error(w, "Vehicle not found", http.StatusNotFound)
		return
	}

	vc.Logger.Info("Successfully found vehicle",
		zap.String("id", vehicle.ID.String()),
		zap.String("model", vehicle.Model),
		zap.String("brand", vehicle.Brand),
		zap.String("numberPlate", vehicle.NumberPlate))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vehicle)
}

func (vc *VehicleController) FindByCustomerId(w http.ResponseWriter, r *http.Request) {
	vc.Logger.Info("=== VEHICLE FIND BY CUSTOMER ID ENDPOINT CALLED ===")

	vars := mux.Vars(r)
	customerID, err := uuid.Parse(vars["customerId"])
	if err != nil {
		vc.Logger.Error("Error parsing customer UUID", zap.Error(err))
		http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
		return
	}

	vc.Logger.Info("Parsed customer ID", zap.String("customerID", customerID.String()))

	vc.Logger.Info("Calling FindByCustomerIdVehicle.Process...")
	vehicles, err := vc.FindByCustomerIdVehicle.Process(customerID)
	if err != nil {
		vc.Logger.Error("Error finding vehicles by customer ID", zap.Error(err), zap.String("customerID", customerID.String()))
		http.Error(w, "Error finding vehicles", http.StatusInternalServerError)
		return
	}

	vc.Logger.Info("Successfully found vehicles",
		zap.String("customerID", customerID.String()),
		zap.Int("count", len(vehicles)))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vehicles)
}

func (vc *VehicleController) UpdateById(w http.ResponseWriter, r *http.Request) {
	vc.Logger.Info("=== VEHICLE UPDATE BY ID ENDPOINT CALLED ===")

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		vc.Logger.Error("Error parsing UUID", zap.Error(err))
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	vc.Logger.Info("Parsed vehicle ID", zap.String("id", id.String()))

	var dto VehicleDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		vc.Logger.Error("Error decoding JSON", zap.Error(err))
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	vc.Logger.Info("Received vehicle update request",
		zap.String("id", id.String()),
		zap.String("model", dto.Model),
		zap.String("brand", dto.Brand),
		zap.String("numberPlate", dto.NumberPlate),
		zap.Int("releaseYear", dto.ReleaseYear),
		zap.Any("customerID", dto.CustomerID))

	if err := dto.Validate(); err != nil {
		vc.Logger.Error("Validation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vc.Logger.Info("Validation passed")

	parsedCustomerID, err := uuid.Parse(dto.CustomerID)
	if err != nil {
		vc.Logger.Error("Error parsing customer ID", zap.Error(err))
		http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
		return
	}
	customerID := parsedCustomerID
	vc.Logger.Info("Customer ID parsed successfully", zap.String("customerID", customerID.String()))

	entity := &domain.Vehicle{
		ID:                          id,
		Model:                       dto.Model,
		Brand:                       dto.Brand,
		ReleaseYear:                 dto.ReleaseYear,
		VehicleIdentificationNumber: dto.VehicleIdentificationNumber,
		NumberPlate:                 dto.NumberPlate,
		Color:                       dto.Color,
		CustomerID:                  customerID,
	}

	vc.Logger.Info("Entity created for update",
		zap.String("id", entity.ID.String()),
		zap.String("model", entity.Model),
		zap.String("brand", entity.Brand),
		zap.String("numberPlate", entity.NumberPlate))

	vc.Logger.Info("Calling UpdateByIdVehicle.Process...")
	err = vc.UpdateByIdVehicle.Process(id, entity)
	if err != nil {
		vc.Logger.Error("Error updating vehicle", zap.Error(err))

		// Tratamento específico para placa duplicada
		if err.Error() == "number plate already exists" {
			http.Error(w, "Number plate already exists", http.StatusConflict)
			return
		}

		// Tratamento específico para customer não encontrado
		if err.Error() == "customer not found" {
			http.Error(w, "Customer not found", http.StatusNotFound)
			return
		}

		// Tratamento específico para vehicle não encontrado
		if err.Error() == "record not found" {
			http.Error(w, "Vehicle not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Error updating vehicle", http.StatusInternalServerError)
		return
	}

	vc.Logger.Info("Vehicle updated successfully",
		zap.String("id", entity.ID.String()),
		zap.String("numberPlate", entity.NumberPlate))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Vehicle updated successfully",
		"id":      entity.ID.String(),
	})
}

func (vc *VehicleController) DeleteById(w http.ResponseWriter, r *http.Request) {
	vc.Logger.Info("=== VEHICLE DELETE BY ID ENDPOINT CALLED ===")

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		vc.Logger.Error("Error parsing UUID", zap.Error(err))
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	vc.Logger.Info("Parsed vehicle ID", zap.String("id", id.String()))

	vc.Logger.Info("Calling DeleteByIdVehicle.Process...")
	err = vc.DeleteByIdVehicle.Process(id)
	if err != nil {
		vc.Logger.Error("Error deleting vehicle by ID", zap.Error(err), zap.String("id", id.String()))

		// Tratamento específico para vehicle não encontrado
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Vehicle not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Error deleting vehicle", http.StatusInternalServerError)
		return
	}

	vc.Logger.Info("Successfully deleted vehicle", zap.String("id", id.String()))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Vehicle deleted successfully",
		"id":      id.String(),
	})
}
