package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func initServiceAndCall(req *http.Request) *httptest.ResponseRecorder {
	responseRecorder := httptest.NewRecorder()
	service := SimpleHttpService{}
	service.Init()
	service.Router.ServeHTTP(responseRecorder, req)
	return responseRecorder
}

func TestHealthEndpoint(t *testing.T) {
	req, _ := http.NewRequest("GET", "/health", nil)
	response := initServiceAndCall(req)

	jsonResponse := healthResponse{}
	json.Unmarshal(response.Body.Bytes(), &jsonResponse)

	expectedResponse := healthResponse{"OK"}

	if jsonResponse != expectedResponse {
		t.Error("Service response", jsonResponse, "did not match expected response", expectedResponse)
	}
}

func TestHostEndpoint(t *testing.T) {
	req, _ := http.NewRequest("GET", "/host", nil)
	response := initServiceAndCall(req)

	jsonResponse := hostResponse{}
	json.Unmarshal(response.Body.Bytes(), &jsonResponse)

	hostname, _ := os.Hostname()
	expectedResponse := hostResponse{hostname}

	if jsonResponse != expectedResponse {
		t.Error("Service response", jsonResponse, "did not match expected response", expectedResponse)
	}
}
