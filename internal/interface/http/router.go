package http

import (
	"encoding/json"
	"net/http"

	"github.com/ln0rd/tech_challenge_12soat/internal/interface/http/middleware"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Router struct {
	router *mux.Router
	logger *zap.Logger
}

func NewRouter(logger *zap.Logger) *Router {
	return &Router{router: mux.NewRouter(), logger: logger}
}

func (r *Router) SetupRouter(router *mux.Router) {
	router.Use(middleware.SetHeaders)

	router.HandleFunc("/", r.HomeHandler)
}
func (rt *Router) HomeHandler(w http.ResponseWriter, r *http.Request) {
	rt.logger.Info("Arrived in HomeHandler")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"response": "Hello World"})
}
