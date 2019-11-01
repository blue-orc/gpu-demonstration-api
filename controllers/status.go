package controllers

import (
	"gpu-demonstration-api/utilities"
	"net/http"

	"github.com/gorilla/mux"
)

//InitStatusController Initializes status endpoints
func InitStatusController(r *mux.Router) {
	r.HandleFunc("/status", statusHandler).Methods("GET")
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	utilities.RespondOK(w)
}
