package http

import (
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/http/controller"
	"github.com/ln0rd/tech_challenge_12soat/internal/interface/http/middleware"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Router struct {
	router             *mux.Router
	logger             *zap.Logger
	customerController *controller.CustomerController
	healthController   *controller.HealthController
}

func NewRouter(logger *zap.Logger, customerController *controller.CustomerController, healthController *controller.HealthController) *Router {
	return &Router{
		router:             mux.NewRouter(),
		logger:             logger,
		customerController: customerController,
		healthController:   healthController,
	}
}

func (r *Router) SetupRouter(router *mux.Router) {
	router.Use(middleware.SetHeaders)

	router.HandleFunc("/healthz", r.healthController.Healthz).Methods("GET")
	router.HandleFunc("/customer", r.customerController.Create).Methods("POST")
	router.HandleFunc("/customer", r.customerController.FindAll).Methods("GET")
	router.HandleFunc("/customer/{id}", r.customerController.FindById).Methods("GET")
	router.HandleFunc("/customer/{id}", r.customerController.DeleteById).Methods("DELETE")
}
