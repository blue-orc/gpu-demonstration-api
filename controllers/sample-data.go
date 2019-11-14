package controllers

import (
	"gpu-demonstration-api/database/repositories"
	"gpu-demonstration-api/utilities"
	"net/http"

	"github.com/gorilla/mux"
)

//InitSampleDataController Initializes sample data endpoints
func InitSampleDataController(r *mux.Router) {
	r.HandleFunc("/sampledata", selectSampleDatahandler).Methods("GET")
}

func selectSampleDatahandler(w http.ResponseWriter, r *http.Request) {
	scriptName, err := utilities.ReadStringQueryParameter(r, "scriptName")
	if err != nil {
		utilities.RespondBadRequest(w, err.Error())
		return
	}
	var repo repositories.BatteryRepo
	switch scriptName {
	case "Battery Remaining Useful Life":
		res, err := repo.SelectSampleDischargeTop(100)
		if err != nil {
			utilities.RespondInternalServerError(w, err.Error())
			return
		}
		utilities.RespondJSON(w, res)
		return
	default:
		utilities.RespondBadRequest(w, "Invalid script name")
		return
	}
}
