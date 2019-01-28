package main

import (
	"reflect"
	"testing"
)

var fakeEnvVars = []string{"SERVICE_NAME=simple-http"}
var expectedEnvVar = envVar{"SERVICE_NAME", "simple-http"}

func provideFakeEnvVars() []string {
	return fakeEnvVars
}

func noEnvVars() []string {
	return []string{}
}

func TestReadEnvVarsWhenThereAreNone(t *testing.T) {
	envVars := readEnvVars(noEnvVars)

	if len(envVars) != 0 {
		t.Error("No EnvVars were expected, but there were some")
	}
}

func TestReadEnvVars(t *testing.T) {
	envVars := readEnvVars(provideFakeEnvVars)

	actualEnvVar := envVars[0]
	if len(envVars) != 1 || !reflect.DeepEqual(actualEnvVar, expectedEnvVar) {
		t.Error("Expected envVar ", expectedEnvVar, ", but got ", envVars)
	}
}
