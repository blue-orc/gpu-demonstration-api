package controllers

import (
	//"gpu-demonstration-api/python-job-runner"
	"gpu-demonstration-api/utilities"
	"net/http"

	"github.com/gorilla/mux"
)

//InitPythonJobRunnerController Initializes status endpoints
func InitPythonJobRunnerController(r *mux.Router) {
	r.HandleFunc("/pythonJobRunner", jobRunnerHandler).Methods("GET")
}

func jobRunnerHandler(w http.ResponseWriter, r *http.Request) {
	//go PythonJobRunner.Run()
	utilities.RespondOK(w)
}
