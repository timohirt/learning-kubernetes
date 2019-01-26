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
	Status   string `json:"status"`
	Hostname string `json:"hostname"`
}

func main() {
	var router = mux.NewRouter()
	router.HandleFunc("/health", healthCheck).Methods("GET")

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	log.Println("Running server!")
	log.Fatal(http.ListenAndServe(":30000", loggedRouter))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	response := healthResponse{"OK", hostname}
	json.NewEncoder(w).Encode(response)
}
