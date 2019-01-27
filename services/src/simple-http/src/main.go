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

type SimpleHttpService struct {
	Router *mux.Router
}

func (service *SimpleHttpService) Init() {
	service.Router = mux.NewRouter()
	service.Router.HandleFunc("/health", healthHandler).Methods("GET")
	service.Router.HandleFunc("/host", hostHandler).Methods("GET")
}

func (service *SimpleHttpService) run(addr string) {
	log.Println("Running server!")
	loggingRouter := handlers.LoggingHandler(os.Stdout, service.Router)
	log.Fatal(http.ListenAndServe(addr, loggingRouter))
}

func main() {
	service := SimpleHttpService{}
	service.Init()
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
