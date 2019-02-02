package main

import (
	"reflect"
	"testing"
)

var noEnvVars = []string{}
var fakeEnvVars = []string{"SERVICE_NAME=simple-http"}
var expectedEnvVar = envVar{"SERVICE_NAME", "simple-http"}

func TestReadEnvVarsWhenThereAreNone(t *testing.T) {
	envLoader := EnvLoader{noEnvVars}
	envVars := envLoader.readEnvVars()

	if len(envVars) != 0 {
		t.Error("No EnvVars were expected, but there were some")
	}
}

func TestReadEnvVars(t *testing.T) {
	envLoader := EnvLoader{fakeEnvVars}
	envVars := envLoader.readEnvVars()

	actualEnvVar := envVars[0]
	if len(envVars) != 1 || !reflect.DeepEqual(actualEnvVar, expectedEnvVar) {
		t.Error("Expected envVar ", expectedEnvVar, ", but got ", envVars)
	}
}
