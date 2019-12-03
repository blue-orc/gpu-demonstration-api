package controllers

import (
	"gpu-demonstration-api/device-monitor"
	"gpu-demonstration-api/utilities"
	"net/http"

	"github.com/gorilla/mux"
)

//InitDeviceMontitorController Initializes device endpoints
func InitDeviceMontitorController(r *mux.Router) {
	r.HandleFunc("/devicemonitor/cpu", getCPUInfo).Methods("GET")
	r.HandleFunc("/devicemonitor/gpu", getGPUInfo).Methods("GET")
}

func getCPUInfo(w http.ResponseWriter, r *http.Request) {
	i, err := devicemonitor.GetCPUInfo()
	if err != nil {
		utilities.RespondInternalServerError(w, err.Error())
		return
	}
	utilities.RespondJSON(w, i)
}

func getGPUInfo(w http.ResponseWriter, r *http.Request) {
	d, err := devicemonitor.GetGPUInfo()
	if err != nil {
		utilities.RespondInternalServerError(w, err.Error())
		return
	}
	utilities.RespondJSON(w, d)
}
