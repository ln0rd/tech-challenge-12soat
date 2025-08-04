package http

import (
	"net/http"

	"github.com/ln0rd/tech_challenge_12soat/internal/interface/http/controller"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/http/middleware"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Router struct {
	router             *mux.Router
	logger             *zap.Logger
	customerController *controller.CustomerController
	userController     *controller.UserController
	authController     *controller.AuthController
	healthController   *controller.HealthController
	vehicleController  *controller.VehicleController
	inputController    *controller.InputController
	orderController    *controller.OrderController
	authMiddleware     *middleware.AuthMiddleware
	authzMiddleware    *middleware.AuthorizationMiddleware
}

func NewRouter(logger *zap.Logger, customerController *controller.CustomerController, userController *controller.UserController, authController *controller.AuthController, healthController *controller.HealthController, vehicleController *controller.VehicleController, inputController *controller.InputController, orderController *controller.OrderController, authMiddleware *middleware.AuthMiddleware, authzMiddleware *middleware.AuthorizationMiddleware) *Router {
	return &Router{
		router:             mux.NewRouter(),
		logger:             logger,
		customerController: customerController,
		userController:     userController,
		authController:     authController,
		healthController:   healthController,
		vehicleController:  vehicleController,
		inputController:    inputController,
		orderController:    orderController,
		authMiddleware:     authMiddleware,
		authzMiddleware:    authzMiddleware,
	}
}

func (r *Router) SetupRouter(router *mux.Router) {
	router.Use(middleware.SetHeaders)

	r.logger.Info("Setting up routes...")

	// Rotas públicas (sem autenticação)
	router.HandleFunc("/healthz", r.healthController.Healthz).Methods("GET")
	r.logger.Info("Route registered: GET /healthz")

	router.HandleFunc("/auth/login", r.authController.Login).Methods("POST")
	r.logger.Info("Route registered: POST /auth/login")

	router.HandleFunc("/user", r.userController.Create).Methods("POST")
	r.logger.Info("Route registered: POST /user")

	// ===== ROTAS PARA MECHANIC E ADMIN =====
	// Customer routes - mechanic e admin
	router.Handle("/customer", r.authMiddleware.Authenticate(r.authzMiddleware.RequireMechanicOrAdmin(http.HandlerFunc(r.customerController.Create)))).Methods("POST")
	r.logger.Info("Route registered: POST /customer (MECHANIC & ADMIN)")

	router.Handle("/customer", r.authMiddleware.Authenticate(r.authzMiddleware.RequireMechanicOrAdmin(http.HandlerFunc(r.customerController.FindAll)))).Methods("GET")
	r.logger.Info("Route registered: GET /customer (MECHANIC & ADMIN)")

	router.Handle("/customer/{id}", r.authMiddleware.Authenticate(r.authzMiddleware.RequireMechanicOrAdmin(http.HandlerFunc(r.customerController.FindById)))).Methods("GET")
	r.logger.Info("Route registered: GET /customer/{id} (MECHANIC & ADMIN)")

	router.Handle("/customer/{id}", r.authMiddleware.Authenticate(r.authzMiddleware.RequireMechanicOrAdmin(http.HandlerFunc(r.customerController.UpdateById)))).Methods("PUT")
	r.logger.Info("Route registered: PUT /customer/{id} (MECHANIC & ADMIN)")

	router.Handle("/customer/{id}", r.authMiddleware.Authenticate(r.authzMiddleware.RequireMechanicOrAdmin(http.HandlerFunc(r.customerController.DeleteById)))).Methods("DELETE")
	r.logger.Info("Route registered: DELETE /customer/{id} (MECHANIC & ADMIN)")

	// Vehicle routes - mechanic e admin
	router.Handle("/vehicle", r.authMiddleware.Authenticate(r.authzMiddleware.RequireMechanicOrAdmin(http.HandlerFunc(r.vehicleController.Create)))).Methods("POST")
	r.logger.Info("Route registered: POST /vehicle (MECHANIC & ADMIN)")

	router.Handle("/vehicle/{id}", r.authMiddleware.Authenticate(r.authzMiddleware.RequireMechanicOrAdmin(http.HandlerFunc(r.vehicleController.FindById)))).Methods("GET")
	r.logger.Info("Route registered: GET /vehicle/{id} (MECHANIC & ADMIN)")

	router.Handle("/vehicle/{id}", r.authMiddleware.Authenticate(r.authzMiddleware.RequireMechanicOrAdmin(http.HandlerFunc(r.vehicleController.UpdateById)))).Methods("PUT")
	r.logger.Info("Route registered: PUT /vehicle/{id} (MECHANIC & ADMIN)")

	router.Handle("/vehicle/{id}", r.authMiddleware.Authenticate(r.authzMiddleware.RequireMechanicOrAdmin(http.HandlerFunc(r.vehicleController.DeleteById)))).Methods("DELETE")
	r.logger.Info("Route registered: DELETE /vehicle/{id} (MECHANIC & ADMIN)")

	router.Handle("/vehicle/customer/{customerId}", r.authMiddleware.Authenticate(r.authzMiddleware.RequireMechanicOrAdmin(http.HandlerFunc(r.vehicleController.FindByCustomerId)))).Methods("GET")
	r.logger.Info("Route registered: GET /vehicle/customer/{customerId} (MECHANIC & ADMIN)")

	// Input routes - mechanic e admin
	router.Handle("/input", r.authMiddleware.Authenticate(r.authzMiddleware.RequireMechanicOrAdmin(http.HandlerFunc(r.inputController.Create)))).Methods("POST")
	r.logger.Info("Route registered: POST /input (MECHANIC & ADMIN)")

	router.Handle("/input", r.authMiddleware.Authenticate(r.authzMiddleware.RequireMechanicOrAdmin(http.HandlerFunc(r.inputController.FindAll)))).Methods("GET")
	r.logger.Info("Route registered: GET /input (MECHANIC & ADMIN)")

	router.Handle("/input/{id}", r.authMiddleware.Authenticate(r.authzMiddleware.RequireMechanicOrAdmin(http.HandlerFunc(r.inputController.FindById)))).Methods("GET")
	r.logger.Info("Route registered: GET /input/{id} (MECHANIC & ADMIN)")

	router.Handle("/input/{id}", r.authMiddleware.Authenticate(r.authzMiddleware.RequireMechanicOrAdmin(http.HandlerFunc(r.inputController.UpdateById)))).Methods("PUT")
	r.logger.Info("Route registered: PUT /input/{id} (MECHANIC & ADMIN)")

	router.Handle("/input/{id}", r.authMiddleware.Authenticate(r.authzMiddleware.RequireMechanicOrAdmin(http.HandlerFunc(r.inputController.DeleteById)))).Methods("DELETE")
	r.logger.Info("Route registered: DELETE /input/{id} (MECHANIC & ADMIN)")

	// Order routes - mechanic e admin
	router.Handle("/order", r.authMiddleware.Authenticate(r.authzMiddleware.RequireMechanicOrAdmin(http.HandlerFunc(r.orderController.Create)))).Methods("POST")
	r.logger.Info("Route registered: POST /order (MECHANIC & ADMIN)")

	router.Handle("/order/{orderId}/input", r.authMiddleware.Authenticate(r.authzMiddleware.RequireMechanicOrAdmin(http.HandlerFunc(r.orderController.AddInputToOrder)))).Methods("POST")
	r.logger.Info("Route registered: POST /order/{orderId}/input (MECHANIC & ADMIN)")

	router.Handle("/order/{orderId}/input/remove", r.authMiddleware.Authenticate(r.authzMiddleware.RequireMechanicOrAdmin(http.HandlerFunc(r.orderController.RemoveInputFromOrder)))).Methods("POST")
	r.logger.Info("Route registered: POST /order/{orderId}/input/remove (MECHANIC & ADMIN)")

	router.Handle("/order/{orderId}/status", r.authMiddleware.Authenticate(r.authzMiddleware.RequireMechanicOrAdmin(http.HandlerFunc(r.orderController.UpdateOrderStatus)))).Methods("PUT")
	r.logger.Info("Route registered: PUT /order/{orderId}/status (MECHANIC & ADMIN)")

	// Order overview - todos os tipos de usuário podem acessar
	router.Handle("/order/{orderId}/overview", r.authMiddleware.Authenticate(http.HandlerFunc(r.orderController.FindOrderOverviewById))).Methods("GET")
	r.logger.Info("Route registered: GET /order/{orderId}/overview (ALL AUTHENTICATED USERS)")

	r.logger.Info("All routes registered successfully")
}
