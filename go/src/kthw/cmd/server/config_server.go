package server

import (
	"fmt"
	"kthw/cmd/common"
	"kthw/cmd/sshkey"

	viper "github.com/spf13/viper"
)

const (
	confHCloudDefaultServerTypeKey = "hcloud.default.serverType"
	confHCloudDefaultImageNameKey  = "hcloud.default.imageName"
	confHCloudLocationNameKey      = "hcloud.default.locationName"

	// HCloudServerType is the default server type used when adding servers.
	HCloudServerType = "cx21"
	// HCloudImage is the default image used when adding servers.
	HCloudImage = "ubuntu-18.04"
	// HCloudLocation is the default location (like datacenter) where a added server is created at.
	HCloudLocation = "nbg1"
)

// Config from config file
type Config struct {
	ID             int
	Name           string
	ServerType     string
	ImageName      string
	LocationName   string
	PublicIP       string
	RootPassword   string
	SSHPublicKeyID int
}

// UpdateConfig updates the configuration with the current field values. Changes are not persisted.
func (sc *Config) UpdateConfig() {
	viper.Set(sc.confServerNameKey(), sc.Name)
	viper.Set(sc.confServerTypeKey(), sc.ServerType)
	viper.Set(sc.confLocationNameKey(), sc.LocationName)
	viper.Set(sc.confImageNameKey(), sc.ImageName)
	viper.Set(sc.confSSKPublicKeyID(), sc.SSHPublicKeyID)

	if sc.ID != 0 {
		viper.Set(sc.confIDKey(), sc.ID)
	}

	if sc.PublicIP != "" {
		viper.Set(sc.confPublicIPKey(), sc.PublicIP)
	}

	if sc.RootPassword != "" {
		viper.Set(sc.confRootPasswordKey(), sc.RootPassword)
	}
}

// ReadFromConfig reads the config of a server from the configuration file.
// Name field of Config must be set.
func (sc *Config) ReadFromConfig() error {
	if sc.Name == "" {
		return fmt.Errorf("Could not read server from config. Server name not set")
	}
	// TODO verify server name exists in config and return error if not

	publicIP := viper.GetString(sc.confPublicIPKey())
	if publicIP != "" {
		sc.PublicIP = publicIP
	}
	rootPassword := viper.GetString(sc.confRootPasswordKey())
	if rootPassword != "" {
		sc.RootPassword = rootPassword
	}

	id := viper.GetInt(sc.confIDKey())
	if id != 0 {
		sc.ID = id
	}

	sc.SSHPublicKeyID = viper.GetInt(sc.confSSKPublicKeyID())
	sc.ServerType = viper.GetString(sc.confServerTypeKey())
	sc.ImageName = viper.GetString(sc.confImageNameKey())
	sc.LocationName = viper.GetString(sc.confLocationNameKey())
	return nil
}

func (sc *Config) confSSKPublicKeyID() string {
	return fmt.Sprintf("hcloud.server.%s.publicKeyId", sc.Name)
}

func (sc *Config) confServerNameKey() string {
	return fmt.Sprintf("hcloud.server.%s.name", sc.Name)
}

func (sc *Config) confServerTypeKey() string {
	return fmt.Sprintf("hcloud.server.%s.serverType", sc.Name)
}

func (sc *Config) confImageNameKey() string {
	return fmt.Sprintf("hcloud.server.%s.imageName", sc.Name)
}

func (sc *Config) confLocationNameKey() string {
	return fmt.Sprintf("hcloud.server.%s.locationName", sc.Name)
}

func (sc *Config) confPublicIPKey() string {
	return fmt.Sprintf("hcloud.server.%s.publicIP", sc.Name)
}

func (sc *Config) confRootPasswordKey() string {
	return fmt.Sprintf("hcloud.server.%s.rootPassword", sc.Name)
}

func (sc *Config) confIDKey() string {
	return fmt.Sprintf("hcloud.server.%s.id", sc.Name)
}

// FromConfig reads settings of a specific server from config.
func FromConfig(serverName string) Config {
	serverConfig := Config{Name: serverName}
	err := serverConfig.ReadFromConfig()
	common.WhenErrPrintAndExit(err)
	return serverConfig
}

// SetHCloudServerDefaults sets default server type, image and location, which are used to add servers.
func SetHCloudServerDefaults() {
	viper.Set(confHCloudDefaultServerTypeKey, HCloudServerType)
	viper.Set(confHCloudDefaultImageNameKey, HCloudImage)
	viper.Set(confHCloudLocationNameKey, HCloudLocation)
}

// AddServer uses the first argument as server name and adds this server to the configuration.
func AddServer(serverName string) {
	sshKey, err := sshkey.ReadSSHPublicKeyFromConf()
	common.WhenErrPrintAndExit(err)
	// TODO Ensure sshkey is provisioned
	serverConf := Config{
		Name:           serverName,
		SSHPublicKeyID: sshKey.ID,
		ServerType:     viper.GetString(confHCloudDefaultServerTypeKey),
		ImageName:      viper.GetString(confHCloudDefaultImageNameKey),
		LocationName:   viper.GetString(confHCloudLocationNameKey),
	}
	serverConf.UpdateConfig()
}
