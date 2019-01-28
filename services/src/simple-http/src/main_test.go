package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func getFakeEnvVars() []string {
	return fakeEnvVars
}

func initServiceAndCall(req *http.Request) *httptest.ResponseRecorder {
	responseRecorder := httptest.NewRecorder()
	service := SimpleHttpService{}
	service.Init(getFakeEnvVars)
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

func TestEnvEndpoint(t *testing.T) {
	req, _ := http.NewRequest("GET", "/env", nil)
	response := initServiceAndCall(req)

	jsonResponse := envResponse{}
	json.Unmarshal(response.Body.Bytes(), &jsonResponse)

	expectedEnvVars := []envVar{expectedEnvVar}
	expectedResponse := envResponse{expectedEnvVars}
	if !reflect.DeepEqual(jsonResponse, expectedResponse) {
		t.Error("Service response", jsonResponse, "did not match expected response", expectedResponse)
	}
}
