package main

import (
	"log"
	"os"

	authInfra "github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/auth"
	db "github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db"
	routes "github.com/ln0rd/tech_challenge_12soat/internal/interface/http"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/http/controller"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/http/middleware"
	authUseCase "github.com/ln0rd/tech_challenge_12soat/internal/usecase/auth"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/customer"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/input"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/order"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/order_input"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/user"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/vehicle"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	err    error
	logger *zap.Logger
)

func main() {
	config := zap.NewProductionConfig()
	if os.Getenv("ENVIRONMENT_LEVEL") == "development" {
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	logger, err = config.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Define o logger global
	zap.ReplaceGlobals(logger)

	if err := godotenv.Load(); err != nil {
		logger.Error("Error loading .env file",
			zap.Error(err),
			zap.String("current_dir", getCurrentDir()))
	} else {
		logger.Info("Successfully loaded .env file")
		logger.Debug("Environment variables",
			zap.String("DB_HOST", os.Getenv("DATABASE_HOST")),
			zap.String("DB_NAME", os.Getenv("DATABASE_NAME")),
			zap.String("DB_PORT", os.Getenv("DATABASE_PORT")))
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		logger.Warn("HTTP_PORT not set, using default 8080")
		httpPort = "8080"
	}

	db.InitDB(logger)

	logger.Info("Initializing the application...")
	r := mux.NewRouter()

	customerController, healthController, userController, authController, vehicleController, inputController, orderController, authMiddleware := InitInstances()

	rt := routes.NewRouter(logger, customerController, userController, authController, healthController, vehicleController, inputController, orderController, authMiddleware)
	rt.SetupRouter(r)

	logger.Info("Server starting", zap.String("port", httpPort))

	log.Fatal(http.ListenAndServe(":"+httpPort, r))
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "error getting current dir"
	}
	return dir
}

func InitInstances() (*controller.CustomerController, *controller.HealthController, *controller.UserController, *controller.AuthController, *controller.VehicleController, *controller.InputController, *controller.OrderController, *middleware.AuthMiddleware) {
	createCustomerUC := &customer.CreateCustomer{DB: db.DB, Logger: logger}
	findAllCustomerUC := &customer.FindAllCustomer{DB: db.DB, Logger: logger}
	findByIdCustomerUC := &customer.FindByIdCustomer{DB: db.DB, Logger: logger}
	deleteByIdCustomerUC := &customer.DeleteByIdCustomer{DB: db.DB, Logger: logger}
	updateByIdCustomerUC := &customer.UpdateByIdCustomer{DB: db.DB, Logger: logger}

	customerController := &controller.CustomerController{
		Logger:             logger,
		CreateCustomer:     createCustomerUC,
		FindAllCustomer:    findAllCustomerUC,
		FindByIdCustomer:   findByIdCustomerUC,
		DeleteByIdCustomer: deleteByIdCustomerUC,
		UpdateByIdCustomer: updateByIdCustomerUC,
	}

	createUserUC := &user.CreateUser{DB: db.DB, Logger: logger}
	userController := &controller.UserController{
		Logger:     logger,
		CreateUser: createUserUC,
	}

	createVehicleUC := &vehicle.CreateVehicle{DB: db.DB, Logger: logger}
	findByIdVehicleUC := &vehicle.FindByIdVehicle{DB: db.DB, Logger: logger}
	findByCustomerIdVehicleUC := &vehicle.FindByCustomerIdVehicle{DB: db.DB, Logger: logger}
	updateByIdVehicleUC := &vehicle.UpdateByIdVehicle{DB: db.DB, Logger: logger}
	deleteByIdVehicleUC := &vehicle.DeleteByIdVehicle{DB: db.DB, Logger: logger}
	vehicleController := &controller.VehicleController{
		Logger:                  logger,
		CreateVehicle:           createVehicleUC,
		FindByIdVehicle:         findByIdVehicleUC,
		FindByCustomerIdVehicle: findByCustomerIdVehicleUC,
		UpdateByIdVehicle:       updateByIdVehicleUC,
		DeleteByIdVehicle:       deleteByIdVehicleUC,
	}

	createInputUC := &input.CreateInput{DB: db.DB, Logger: logger}
	findByIdInputUC := &input.FindByIdInput{DB: db.DB, Logger: logger}
	findAllInputsUC := &input.FindAllInputs{DB: db.DB, Logger: logger}
	updateByIdInputUC := &input.UpdateByIdInput{DB: db.DB, Logger: logger}
	deleteByIdInputUC := &input.DeleteByIdInput{DB: db.DB, Logger: logger}
	inputController := &controller.InputController{
		Logger:          logger,
		CreateInput:     createInputUC,
		FindByIdInput:   findByIdInputUC,
		FindAllInputs:   findAllInputsUC,
		UpdateByIdInput: updateByIdInputUC,
		DeleteByIdInput: deleteByIdInputUC,
	}

	createOrderUC := &order.CreateOrder{DB: db.DB, Logger: logger}
	findCompletedOrderByIdUC := &order.FindCompletedOrderById{DB: db.DB, Logger: logger}
	updateOrderStatusUC := &order.UpdateOrderStatus{DB: db.DB, Logger: logger}

	// Order Input usecases
	decreaseQuantityInputUC := &input.DecreaseQuantityInput{DB: db.DB, Logger: logger}
	increaseQuantityInputUC := &input.IncreaseQuantityInput{DB: db.DB, Logger: logger}

	addInputToOrderUC := &order_input.AddInputToOrder{
		DB:                    db.DB,
		Logger:                logger,
		DecreaseQuantityInput: decreaseQuantityInputUC,
	}

	removeInputFromOrderUC := &order_input.RemoveInputFromOrder{
		DB:                    db.DB,
		Logger:                logger,
		IncreaseQuantityInput: increaseQuantityInputUC,
	}

	orderController := &controller.OrderController{
		Logger:                   logger,
		CreateOrder:              createOrderUC,
		FindCompletedOrderByIdUC: findCompletedOrderByIdUC,
		UpdateOrderStatusUC:      updateOrderStatusUC,
		AddInputToOrderUC:        addInputToOrderUC,
		RemoveInputFromOrderUC:   removeInputFromOrderUC,
	}

	// Auth components
	authRepository := authInfra.NewAuthRepository(db.DB, logger)
	jwtService := authInfra.NewJWTService(logger)
	loginUseCase := authUseCase.NewLoginUseCase(authRepository, jwtService, logger)

	authController := &controller.AuthController{
		Logger:       logger,
		LoginUseCase: loginUseCase,
	}

	// Auth middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService, logger)

	healthController := &controller.HealthController{}

	return customerController, healthController, userController, authController, vehicleController, inputController, orderController, authMiddleware
}
