package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type healthResponse struct {
	Status string `json:"status"`
}

type hostResponse struct {
	Hostname string `json:"hostname"`
}

type envResponse struct {
	EnvVars []envVar `json:"envVars"`
}

type SimpleHttpService struct {
	Router *mux.Router
}

func (service *SimpleHttpService) Init(el *EnvLoader) {
	service.Router = mux.NewRouter()
	service.Router.HandleFunc("/health", healthHandler).Methods("GET")
	service.Router.HandleFunc("/host", hostHandler).Methods("GET")
	service.Router.HandleFunc("/env", envHandler(el)).Methods("GET")
}

func (service *SimpleHttpService) run(addr string) {
	log.Println("Running server!")
	loggingRouter := handlers.LoggingHandler(os.Stdout, service.Router)
	log.Fatal(http.ListenAndServe(addr, loggingRouter))
}

func main() {
	var envLoader *EnvLoader = NewEnvLoader()
	service := SimpleHttpService{}
	service.Init(envLoader)
	service.run(":30000")
}

func hostHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	response := hostResponse{hostname}
	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := healthResponse{"OK"}
	json.NewEncoder(w).Encode(response)
}

func envHandler(el *EnvLoader) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		envVars := el.readEnvVars()
		response := envResponse{envVars}
		json.NewEncoder(w).Encode(response)
	}
}
