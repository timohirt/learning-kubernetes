package main

import (
	"os"
	"strings"
)

type envVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type EnvLoader struct {
	envVarsFromEnvironment []string
}

func (el *EnvLoader) readEnvVars() []envVar {
	var envVars = []envVar{}
	for _, fromEnv := range el.envVarsFromEnvironment {
		keyAndValue := strings.Split(fromEnv, "=")
		currentEnvVar := envVar{keyAndValue[0], keyAndValue[1]}
		envVars = append(envVars, currentEnvVar)
	}

	return envVars
}

func NewEnvLoader() *EnvLoader {
	envVarsFromEnvironment := os.Environ()
	newInstance := EnvLoader{envVarsFromEnvironment}
	return &newInstance
}
