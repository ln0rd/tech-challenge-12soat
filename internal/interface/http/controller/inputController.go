package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/input"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/input"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	inputNameRegex = regexp.MustCompile(`^[a-zA-Z0-9\s\-_]{2,50}$`)
	inputTypeRegex = regexp.MustCompile(`^(supplie|service)$`)
)

const (
	InputTypeSupplie = "supplie"
	InputTypeService = "service"
)

type InputController struct {
	Logger          *zap.Logger
	CreateInput     *input.CreateInput
	FindByIdInput   *input.FindByIdInput
	FindAllInputs   *input.FindAllInputs
	UpdateByIdInput *input.UpdateByIdInput
	DeleteByIdInput *input.DeleteByIdInput
}

type InputDTO struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	InputType   string  `json:"input_type"`
}

func (dto *InputDTO) Validate() error {
	if !inputNameRegex.MatchString(dto.Name) {
		return errors.New("name must contain only letters, numbers, spaces, hyphens and underscores, between 2 and 50 characters")
	}
	if dto.Price <= 0 {
		return errors.New("price must be greater than zero")
	}
	if !inputTypeRegex.MatchString(dto.InputType) {
		return errors.New("input_type must be 'supplie' or 'service'")
	}

	// Validação específica para cada tipo
	if dto.InputType == InputTypeSupplie {
		if dto.Quantity <= 0 {
			return errors.New("quantity must be greater than zero for supplie type")
		}
	} else if dto.InputType == InputTypeService {
		// Para service, a quantidade será sempre 1, mas validamos se foi enviado
		if dto.Quantity <= 0 {
			return errors.New("quantity must be greater than zero for service type")
		}
	}

	if len(dto.Description) > 500 {
		return errors.New("description must be less than 500 characters")
	}
	return nil
}

func (ic *InputController) Create(w http.ResponseWriter, r *http.Request) {
	ic.Logger.Info("=== INPUT CREATE ENDPOINT CALLED ===")

	var dto InputDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		ic.Logger.Error("Error decoding JSON", zap.Error(err))
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	ic.Logger.Info("Received input creation request",
		zap.String("name", dto.Name),
		zap.Float64("price", dto.Price),
		zap.Int("quantity", dto.Quantity),
		zap.String("description", dto.Description),
		zap.String("inputType", dto.InputType))

	if err := dto.Validate(); err != nil {
		ic.Logger.Error("Validation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ic.Logger.Info("Validation passed")

	// Ajusta a quantidade baseado no tipo
	finalQuantity := dto.Quantity
	if dto.InputType == InputTypeService {
		finalQuantity = 1
		ic.Logger.Info("Forcing quantity to 1 for service type",
			zap.String("inputType", dto.InputType),
			zap.Int("originalQuantity", dto.Quantity),
			zap.Int("finalQuantity", finalQuantity))
	}

	entity := &domain.Input{
		ID:          uuid.New(),
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		Quantity:    finalQuantity,
		InputType:   dto.InputType,
	}

	ic.Logger.Info("Entity created",
		zap.String("id", entity.ID.String()),
		zap.String("name", entity.Name),
		zap.Float64("price", entity.Price),
		zap.Int("quantity", entity.Quantity),
		zap.String("inputType", entity.InputType))

	ic.Logger.Info("Calling CreateInput.Process...")
	err := ic.CreateInput.Process(entity)
	if err != nil {
		ic.Logger.Error("Error creating input", zap.Error(err))

		// Tratamento específico para nome duplicado
		if err.Error() == "input name already exists" {
			http.Error(w, "Input name already exists", http.StatusConflict)
			return
		}

		http.Error(w, "Error creating input", http.StatusInternalServerError)
		return
	}

	ic.Logger.Info("Input created successfully",
		zap.String("id", entity.ID.String()),
		zap.String("name", entity.Name))

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Input created successfully",
		"id":      entity.ID.String(),
	})
}

func (ic *InputController) FindById(w http.ResponseWriter, r *http.Request) {
	ic.Logger.Info("=== INPUT FIND BY ID ENDPOINT CALLED ===")

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		ic.Logger.Error("Error parsing UUID", zap.Error(err))
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	ic.Logger.Info("Parsed input ID", zap.String("id", id.String()))

	ic.Logger.Info("Calling FindByIdInput.Process...")
	input, err := ic.FindByIdInput.Process(id)
	if err != nil {
		ic.Logger.Error("Error finding input by ID", zap.Error(err), zap.String("id", id.String()))
		http.Error(w, "Input not found", http.StatusNotFound)
		return
	}

	ic.Logger.Info("Successfully found input",
		zap.String("id", input.ID.String()),
		zap.String("name", input.Name),
		zap.String("inputType", input.InputType),
		zap.Float64("price", input.Price),
		zap.Int("quantity", input.Quantity))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(input)
}

func (ic *InputController) FindAll(w http.ResponseWriter, r *http.Request) {
	ic.Logger.Info("=== INPUT FIND ALL ENDPOINT CALLED ===")

	ic.Logger.Info("Calling FindAllInputs.Process...")
	inputs, err := ic.FindAllInputs.Process()
	if err != nil {
		ic.Logger.Error("Error finding all inputs", zap.Error(err))
		http.Error(w, "Error finding inputs", http.StatusInternalServerError)
		return
	}

	ic.Logger.Info("Successfully found inputs", zap.Int("count", len(inputs)))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(inputs)
}

func (ic *InputController) UpdateById(w http.ResponseWriter, r *http.Request) {
	ic.Logger.Info("=== INPUT UPDATE BY ID ENDPOINT CALLED ===")

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		ic.Logger.Error("Error parsing UUID", zap.Error(err))
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	ic.Logger.Info("Parsed input ID", zap.String("id", id.String()))

	var dto InputDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		ic.Logger.Error("Error decoding JSON", zap.Error(err))
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	ic.Logger.Info("Received input update request",
		zap.String("id", id.String()),
		zap.String("name", dto.Name),
		zap.String("inputType", dto.InputType),
		zap.Float64("price", dto.Price),
		zap.Int("quantity", dto.Quantity),
		zap.String("description", dto.Description))

	if err := dto.Validate(); err != nil {
		ic.Logger.Error("Validation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ic.Logger.Info("Validation passed")

	// Ajusta a quantidade baseado no tipo
	finalQuantity := dto.Quantity
	if dto.InputType == InputTypeService {
		finalQuantity = 1
		ic.Logger.Info("Forcing quantity to 1 for service type",
			zap.String("inputType", dto.InputType),
			zap.Int("originalQuantity", dto.Quantity),
			zap.Int("finalQuantity", finalQuantity))
	}

	entity := &domain.Input{
		ID:          id,
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		Quantity:    finalQuantity,
		InputType:   dto.InputType,
	}

	ic.Logger.Info("Entity created for update",
		zap.String("id", entity.ID.String()),
		zap.String("name", entity.Name),
		zap.String("inputType", entity.InputType),
		zap.Float64("price", entity.Price),
		zap.Int("quantity", entity.Quantity))

	ic.Logger.Info("Calling UpdateByIdInput.Process...")
	err = ic.UpdateByIdInput.Process(id, entity)
	if err != nil {
		ic.Logger.Error("Error updating input", zap.Error(err))

		// Tratamento específico para nome duplicado
		if err.Error() == "input name already exists" {
			http.Error(w, "Input name already exists", http.StatusConflict)
			return
		}

		// Tratamento específico para input não encontrado
		if err.Error() == "record not found" {
			http.Error(w, "Input not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Error updating input", http.StatusInternalServerError)
		return
	}

	ic.Logger.Info("Input updated successfully",
		zap.String("id", entity.ID.String()),
		zap.String("name", entity.Name))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Input updated successfully",
		"id":      entity.ID.String(),
	})
}

func (ic *InputController) DeleteById(w http.ResponseWriter, r *http.Request) {
	ic.Logger.Info("=== INPUT DELETE BY ID ENDPOINT CALLED ===")

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		ic.Logger.Error("Error parsing UUID", zap.Error(err))
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	ic.Logger.Info("Parsed input ID", zap.String("id", id.String()))

	ic.Logger.Info("Calling DeleteByIdInput.Process...")
	err = ic.DeleteByIdInput.Process(id)
	if err != nil {
		ic.Logger.Error("Error deleting input by ID", zap.Error(err), zap.String("id", id.String()))

		// Tratamento específico para input não encontrado
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Input not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Error deleting input", http.StatusInternalServerError)
		return
	}

	ic.Logger.Info("Successfully deleted input", zap.String("id", id.String()))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Input deleted successfully",
		"id":      id.String(),
	})
}
