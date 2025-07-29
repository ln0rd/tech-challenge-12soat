package controller

import (
	"net/http"
)

type UserController struct{}

func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement user creation
	w.WriteHeader(http.StatusNotImplemented)
}
