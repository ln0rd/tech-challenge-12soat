package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/user"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/user"
	"go.uber.org/zap"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	passwordRegex = regexp.MustCompile(`^.{6,}$`)
	userTypeRegex = regexp.MustCompile(`^(admin|mechanic|vehicle_owner)$`)
)

const (
	UserTypeAdmin    = "admin"
	UserTypeCustomer = "customer"
)

type UserController struct {
	Logger     *zap.Logger
	CreateUser *user.CreateUser
}

type UserDTO struct {
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	Username   string  `json:"username"`
	UserType   string  `json:"user_type"`
	CustomerID *string `json:"customer_id,omitempty"`
}

func (dto *UserDTO) Validate() error {
	if !emailRegex.MatchString(dto.Email) {
		return errors.New("invalid email format")
	}
	if !usernameRegex.MatchString(dto.Username) {
		return errors.New("username must contain only letters, numbers and underscore, between 3 and 20 characters")
	}
	if !passwordRegex.MatchString(dto.Password) {
		return errors.New("password must be at least 6 characters long")
	}
	if !userTypeRegex.MatchString(dto.UserType) {
		return errors.New("userType must be 'admin', 'mechanic', or 'vehicle_owner'")
	}
	return nil
}

func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) {
	uc.Logger.Info("=== USER CREATE ENDPOINT CALLED ===")

	var dto UserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		uc.Logger.Error("Error decoding JSON", zap.Error(err))
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	uc.Logger.Info("Received user creation request",
		zap.String("email", dto.Email),
		zap.String("username", dto.Username),
		zap.String("userType", dto.UserType),
		zap.Any("customerID", dto.CustomerID))

	if err := dto.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uc.Logger.Info("Validation passed")

	var customerID *uuid.UUID
	if dto.CustomerID != nil {
		parsedCustomerID, err := uuid.Parse(*dto.CustomerID)
		if err != nil {
			uc.Logger.Error("Error parsing customer ID", zap.Error(err))
			http.Error(w, "Invalid customer ID format", http.StatusBadRequest)
			return
		}
		customerID = &parsedCustomerID
		uc.Logger.Info("Customer ID parsed successfully", zap.String("customerID", customerID.String()))
	} else {
		uc.Logger.Info("No customer ID provided")
	}

	entity := &domain.User{
		Email:      dto.Email,
		Password:   dto.Password,
		Username:   dto.Username,
		UserType:   dto.UserType,
		CustomerID: customerID,
	}

	uc.Logger.Info("Entity created",
		zap.String("email", entity.Email),
		zap.String("username", entity.Username),
		zap.String("userType", entity.UserType))

	uc.Logger.Info("Calling CreateUser.Process...")
	err := uc.CreateUser.Process(entity)
	if err != nil {
		uc.Logger.Error("Error creating user", zap.Error(err))

		// Tratamento específico para email duplicado
		if err.Error() == "email already exists" {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}

		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	uc.Logger.Info("User created successfully", zap.String("email", entity.Email))

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}
