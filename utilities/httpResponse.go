package utilities

import (
	"encoding/json"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Unable to parse response payload to JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func RespondBadRequest(w http.ResponseWriter, message string) {
	http.Error(w, "Bad Request: "+message, http.StatusBadRequest)
}

func RespondUnauthorized(w http.ResponseWriter, message string) {
	http.Error(w, "Unauthorized: "+message, http.StatusUnauthorized)
}

func RespondForbidden(w http.ResponseWriter) {
	http.Error(w, "Forbidden", http.StatusForbidden)
}

func RespondInternalServerError(w http.ResponseWriter, message string) {
	http.Error(w, "Internal Server Error: "+message, http.StatusInternalServerError)
}

func RespondOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}
