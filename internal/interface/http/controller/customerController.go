package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

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
	Logger         *zap.Logger
	CreateCustomer *customer.CreateCustomer
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
		return errors.New("nome inválido: apenas letras e espaços, até 255 caracteres")
	}
	if !emailRegex.MatchString(dto.Email) {
		return errors.New("email inválido")
	}
	if !userIDRegex.MatchString(dto.UserID) {
		return errors.New("userID deve conter apenas letras e números")
	}
	if !documentNumberReg.MatchString(dto.DocumentNumber) {
		return errors.New("documentNumber deve conter apenas números e no máximo 50 caracteres")
	}
	if dto.CustomerType != CustomerTypeLegal && dto.CustomerType != CustomerTypeNatural {
		return errors.New("customerType deve ser 'legal_person' ou 'natural_person'")
	}
	return nil
}

func (cc *CustomerController) Create(w http.ResponseWriter, r *http.Request) {
	var dto CustomerDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		cc.Logger.Error("Erro ao decodificar JSON", zap.Error(err))
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
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
		cc.Logger.Error("Erro ao criar customer", zap.Error(err))
		http.Error(w, "Erro ao criar customer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Customer criado com sucesso"})
}
