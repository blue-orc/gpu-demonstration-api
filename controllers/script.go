package controllers

import (
	"gpu-demonstration-api/database/repositories"
	"gpu-demonstration-api/models"
	"gpu-demonstration-api/utilities"
	"net/http"

	"github.com/gorilla/mux"
)

//InitScriptController Initializes script endpoints
func InitScriptController(r *mux.Router) {
	r.HandleFunc("/script", selectScriptHandler).Methods("GET")
}

func selectScriptHandler(w http.ResponseWriter, r *http.Request) {
	scriptID, _ := utilities.ReadIntQueryParameter(r, "scriptID")

	var sr repositories.Script
	if scriptID > 0 {
		res, err := sr.SelectByID(scriptID)
		if err != nil {
			utilities.RespondInternalServerError(w, err.Error())
			return
		}
		utilities.RespondJSON(w, res)
		return
	}

	res, err := sr.SelectAll()
	if err != nil {
		utilities.RespondInternalServerError(w, err.Error())
		return
	}
	utilities.RespondJSON(w, res)
	return
}

func upsertScriptHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok || username != "admin" || password != "YkMg2NEYRCT5zxv8" {
		utilities.RespondUnauthorized(w, "Unauthorized")
		return
	}
	var s models.Script
	err := utilities.ReadJsonHttpBody(r, s)
	if err != nil {
		utilities.RespondInternalServerError(w, err.Error())
		return
	}
	var sr repositories.Script
	if s.ID > 0 {
		err := sr.Update(s)
		if err != nil {
			utilities.RespondInternalServerError(w, err.Error())
			return
		}
		utilities.RespondJSON(w, s)
		return
	}

	err = sr.Insert(s)
	if err != nil {
		utilities.RespondInternalServerError(w, err.Error())
		return
	}
	utilities.RespondJSON(w, s)
	return
}
