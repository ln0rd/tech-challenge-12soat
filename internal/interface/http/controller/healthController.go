package controller

import (
	"encoding/json"
	"net/http"
)

type HealthController struct{}

func (hc *HealthController) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
