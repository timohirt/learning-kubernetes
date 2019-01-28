package main

import (
	"os"
	"strings"
)

type envVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type providesEnvVariables func() []string

func readEnvVars(p providesEnvVariables) []envVar {

	var envVars = []envVar{}
	for _, fromEnv := range p() {
		keyAndValue := strings.Split(fromEnv, "=")
		currentEnvVar := envVar{keyAndValue[0], keyAndValue[1]}
		envVars = append(envVars, currentEnvVar)
	}

	return envVars
}

func getEnvVarsFromEnvironment() []string {
	return os.Environ()
}
