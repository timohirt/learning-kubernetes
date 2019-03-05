package certs

import "github.com/spf13/viper"

const (
	certsBaseDirKey = "certs.baseDir"
)

// Config contains all configuations for generating certs
type Config struct {
	BaseDir string
}

// ReadConfig reads certs configuration from config file.
func ReadConfig() Config {
	return Config{BaseDir: viper.GetString(certsBaseDirKey)}
}

// WriteDefaultConfig sets default certs config.
func WriteDefaultConfig() {
	viper.Set(certsBaseDirKey, certsBaseDir)
}
