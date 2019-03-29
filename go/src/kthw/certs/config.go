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

// InitDefaultConfig sets default certs config.
func InitDefaultConfig() Config {
	viper.Set(certsBaseDirKey, certsBaseDir)
	return ReadConfig()
}
