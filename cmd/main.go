package main

import (
	"log"
	"os"

	db "github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db"
	loggerAdapter "github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/logger"
	"github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/repository"
	routes "github.com/ln0rd/tech_challenge_12soat/internal/interface/http"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/http/controller"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/http/middleware"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/customer"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/input"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/order"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/order_input"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/order_status_history"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/user"
	"github.com/ln0rd/tech_challenge_12soat/internal/usecase/vehicle"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	authInfra "github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/auth"
	authUseCase "github.com/ln0rd/tech_challenge_12soat/internal/usecase/auth"
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

	customerController, healthController, userController, authController, vehicleController, inputController, orderController, authMiddleware, authzMiddleware := InitInstances()

	rt := routes.NewRouter(logger, customerController, userController, authController, healthController, vehicleController, inputController, orderController, authMiddleware, authzMiddleware)
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

func InitInstances() (*controller.CustomerController, *controller.HealthController, *controller.UserController, *controller.AuthController, *controller.VehicleController, *controller.InputController, *controller.OrderController, *middleware.AuthMiddleware, *middleware.AuthorizationMiddleware) {
	// Cria os repositories
	customerRepository := repository.NewCustomerRepositoryAdapter(db.DB)
	userRepository := repository.NewUserRepositoryAdapter(db.DB)
	vehicleRepository := repository.NewVehicleRepositoryAdapter(db.DB)
	inputRepository := repository.NewInputRepositoryAdapter(db.DB)
	orderRepository := repository.NewOrderRepositoryAdapter(db.DB)
	orderInputRepository := repository.NewOrderInputRepositoryAdapter(db.DB)
	orderStatusHistoryRepository := repository.NewOrderStatusHistoryRepositoryAdapter(db.DB)

	// Cria o logger adapter
	loggerAdapter := loggerAdapter.NewZapAdapter(logger)

	createCustomerUC := &customer.CreateCustomer{CustomerRepository: customerRepository, Logger: loggerAdapter}
	findAllCustomerUC := &customer.FindAllCustomer{CustomerRepository: customerRepository, Logger: loggerAdapter}
	findByIdCustomerUC := &customer.FindByIdCustomer{CustomerRepository: customerRepository, Logger: loggerAdapter}
	deleteByIdCustomerUC := &customer.DeleteByIdCustomer{CustomerRepository: customerRepository, Logger: loggerAdapter}
	updateByIdCustomerUC := &customer.UpdateByIdCustomer{CustomerRepository: customerRepository, Logger: loggerAdapter}

	customerController := &controller.CustomerController{
		Logger:             logger,
		CreateCustomer:     createCustomerUC,
		FindAllCustomer:    findAllCustomerUC,
		FindByIdCustomer:   findByIdCustomerUC,
		DeleteByIdCustomer: deleteByIdCustomerUC,
		UpdateByIdCustomer: updateByIdCustomerUC,
	}

	createUserUC := &user.CreateUser{UserRepository: userRepository, Logger: loggerAdapter}
	userController := &controller.UserController{
		Logger:     logger,
		CreateUser: createUserUC,
	}

	createVehicleUC := &vehicle.CreateVehicle{VehicleRepository: vehicleRepository, CustomerRepository: customerRepository, Logger: loggerAdapter}
	findByIdVehicleUC := &vehicle.FindByIdVehicle{VehicleRepository: vehicleRepository, Logger: loggerAdapter}
	findByCustomerIdVehicleUC := &vehicle.FindByCustomerIdVehicle{VehicleRepository: vehicleRepository, Logger: loggerAdapter}
	updateByIdVehicleUC := &vehicle.UpdateByIdVehicle{VehicleRepository: vehicleRepository, CustomerRepository: customerRepository, Logger: loggerAdapter}
	deleteByIdVehicleUC := &vehicle.DeleteByIdVehicle{VehicleRepository: vehicleRepository, Logger: loggerAdapter}

	vehicleController := &controller.VehicleController{
		Logger:                  logger,
		CreateVehicle:           createVehicleUC,
		FindByIdVehicle:         findByIdVehicleUC,
		FindByCustomerIdVehicle: findByCustomerIdVehicleUC,
		UpdateByIdVehicle:       updateByIdVehicleUC,
		DeleteByIdVehicle:       deleteByIdVehicleUC,
	}

	createInputUC := &input.CreateInput{InputRepository: inputRepository, Logger: loggerAdapter}
	findAllInputsUC := &input.FindAllInputs{InputRepository: inputRepository, Logger: loggerAdapter}
	findByIdInputUC := &input.FindByIdInput{InputRepository: inputRepository, Logger: loggerAdapter}
	updateByIdInputUC := &input.UpdateByIdInput{InputRepository: inputRepository, Logger: loggerAdapter}
	deleteByIdInputUC := &input.DeleteByIdInput{InputRepository: inputRepository, Logger: loggerAdapter}

	inputController := &controller.InputController{
		Logger:          logger,
		CreateInput:     createInputUC,
		FindAllInputs:   findAllInputsUC,
		FindByIdInput:   findByIdInputUC,
		UpdateByIdInput: updateByIdInputUC,
		DeleteByIdInput: deleteByIdInputUC,
	}

	// Order usecases
	statusHistoryManager := &order_status_history.ManageOrderStatusHistory{
		OrderStatusHistoryRepository: orderStatusHistoryRepository,
		Logger:                       loggerAdapter,
	}

	createOrderUC := &order.CreateOrder{
		OrderRepository:      orderRepository,
		CustomerRepository:   customerRepository,
		VehicleRepository:    vehicleRepository,
		Logger:               loggerAdapter,
		StatusHistoryManager: statusHistoryManager,
	}

	findOrderOverviewByIdUC := &order.FindOrderOverviewById{
		OrderRepository:              orderRepository,
		VehicleRepository:            vehicleRepository,
		OrderInputRepository:         orderInputRepository,
		OrderStatusHistoryRepository: orderStatusHistoryRepository,
		InputRepository:              inputRepository,
		Logger:                       loggerAdapter,
	}

	// Input quantity usecases
	increaseQuantityInputUC := &input.IncreaseQuantityInput{InputRepository: inputRepository, Logger: loggerAdapter}
	decreaseQuantityInputUC := &input.DecreaseQuantityInput{InputRepository: inputRepository, Logger: loggerAdapter}

	// Order input usecases
	addInputToOrderUC := &order_input.AddInputToOrder{
		OrderRepository:       orderRepository,
		InputRepository:       inputRepository,
		OrderInputRepository:  orderInputRepository,
		Logger:                loggerAdapter,
		DecreaseQuantityInput: decreaseQuantityInputUC,
	}

	removeInputFromOrderUC := &order_input.RemoveInputFromOrder{
		OrderRepository:       orderRepository,
		InputRepository:       inputRepository,
		OrderInputRepository:  orderInputRepository,
		Logger:                loggerAdapter,
		IncreaseQuantityInput: increaseQuantityInputUC,
	}

	// Order status usecase
	updateOrderStatusUC := &order.UpdateOrderStatus{
		OrderRepository:      orderRepository,
		Logger:               loggerAdapter,
		StatusHistoryManager: statusHistoryManager,
	}

	orderController := &controller.OrderController{
		Logger:                  logger,
		CreateOrder:             createOrderUC,
		FindOrderOverviewByIdUC: findOrderOverviewByIdUC,
		AddInputToOrderUC:       addInputToOrderUC,
		RemoveInputFromOrderUC:  removeInputFromOrderUC,
		UpdateOrderStatusUC:     updateOrderStatusUC,
	}

	healthController := &controller.HealthController{}

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
	authzMiddleware := middleware.NewAuthorizationMiddleware(logger)

	return customerController, healthController, userController, authController, vehicleController, inputController, orderController, authMiddleware, authzMiddleware
}
