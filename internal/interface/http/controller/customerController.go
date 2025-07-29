package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/costumer"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/customer"
	"go.uber.org/zap"
)

var (
	nameRegex         = regexp.MustCompile(`^[A-Za-zÀ-ÿãõÃÕ\s]{1,255}$`)
	emailRegex        = regexp.MustCompile(`^[\w._%+-]+@[\w.-]+\.[a-zA-Z]{2,}$`)
	userIDRegex       = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	documentNumberReg = regexp.MustCompile(`^[0-9]{1,50}$`)
)

const (
	CustomerTypeLegal   = "legal_person"
	CustomerTypeNatural = "natural_person"
)

type CustomerController struct {
	Logger             *zap.Logger
	CreateCustomer     *customer.CreateCustomer
	FindAllCustomer    *customer.FindAllCustomer
	FindByIdCustomer   *customer.FindByIdCustomer
	DeleteByIdCustomer *customer.DeleteByIdCustomer
}

type CustomerDTO struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	UserID         string `json:"user_id"`
	DocumentNumber string `json:"document_number"`
	CustomerType   string `json:"customer_type"`
}

func (dto *CustomerDTO) Validate() error {
	if !nameRegex.MatchString(dto.Name) {
		return errors.New("invalid name: only letters and spaces, up to 255 characters")
	}
	if !emailRegex.MatchString(dto.Email) {
		return errors.New("invalid email")
	}
	if !userIDRegex.MatchString(dto.UserID) {
		return errors.New("userID must contain only letters and numbers")
	}
	if !documentNumberReg.MatchString(dto.DocumentNumber) {
		return errors.New("documentNumber must contain only numbers and up to 50 characters")
	}
	if dto.CustomerType != CustomerTypeLegal && dto.CustomerType != CustomerTypeNatural {
		return errors.New("customerType must be 'legal_person' or 'natural_person'")
	}
	return nil
}

func (cc *CustomerController) Create(w http.ResponseWriter, r *http.Request) {
	var dto CustomerDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		cc.Logger.Error("Error decoding JSON", zap.Error(err))
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	if err := dto.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	entity := &domain.Customer{
		Name:           dto.Name,
		Email:          dto.Email,
		UserID:         dto.UserID,
		DocumentNumber: dto.DocumentNumber,
		CustomerType:   dto.CustomerType,
	}

	err := cc.CreateCustomer.Process(entity)
	if err != nil {
		cc.Logger.Error("Error creating customer", zap.Error(err))
		http.Error(w, "Error creating customer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Customer created successfully"})
}

func (cc *CustomerController) FindAll(w http.ResponseWriter, r *http.Request) {
	customers, err := cc.FindAllCustomer.Process()
	if err != nil {
		cc.Logger.Error("Error finding all customers", zap.Error(err))
		http.Error(w, "Error retrieving customers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customers)
}

func (cc *CustomerController) FindById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		cc.Logger.Error("Error parsing UUID", zap.Error(err))
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	customer, err := cc.FindByIdCustomer.Process(id)
	if err != nil {
		cc.Logger.Error("Error finding customer by ID", zap.Error(err), zap.String("id", id.String()))
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

func (cc *CustomerController) UpdateById(w http.ResponseWriter, r *http.Request) {

}

func (cc *CustomerController) DeleteById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		cc.Logger.Error("Error parsing UUID", zap.Error(err))
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = cc.DeleteByIdCustomer.Process(id)
	if err != nil {
		cc.Logger.Error("Error deleting customer by ID", zap.Error(err), zap.String("id", id.String()))
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
