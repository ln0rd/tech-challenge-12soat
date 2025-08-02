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
	authMiddleware     *middleware.AuthMiddleware
}

func NewRouter(logger *zap.Logger, customerController *controller.CustomerController, userController *controller.UserController, authController *controller.AuthController, healthController *controller.HealthController, authMiddleware *middleware.AuthMiddleware) *Router {
	return &Router{
		router:             mux.NewRouter(),
		logger:             logger,
		customerController: customerController,
		userController:     userController,
		authController:     authController,
		healthController:   healthController,
		authMiddleware:     authMiddleware,
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

	// Rotas protegidas (com autenticação)
	router.Handle("/customer", r.authMiddleware.Authenticate(http.HandlerFunc(r.customerController.Create))).Methods("POST")
	r.logger.Info("Route registered: POST /customer (PROTECTED)")

	router.Handle("/customer", r.authMiddleware.Authenticate(http.HandlerFunc(r.customerController.FindAll))).Methods("GET")
	r.logger.Info("Route registered: GET /customer (PROTECTED)")

	router.Handle("/customer/{id}", r.authMiddleware.Authenticate(http.HandlerFunc(r.customerController.FindById))).Methods("GET")
	r.logger.Info("Route registered: GET /customer/{id} (PROTECTED)")

	router.Handle("/customer/{id}", r.authMiddleware.Authenticate(http.HandlerFunc(r.customerController.UpdateById))).Methods("PUT")
	r.logger.Info("Route registered: PUT /customer/{id} (PROTECTED)")

	router.Handle("/customer/{id}", r.authMiddleware.Authenticate(http.HandlerFunc(r.customerController.DeleteById))).Methods("DELETE")
	r.logger.Info("Route registered: DELETE /customer/{id} (PROTECTED)")

	r.logger.Info("All routes registered successfully")
}
