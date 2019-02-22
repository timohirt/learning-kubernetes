package server

import (
	"fmt"
	"kthw/cmd/common"
	"kthw/cmd/infra/sshkey"

	viper "github.com/spf13/viper"
)

const (
	confHCloudServersKey           = "hcloud.server"
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
	PrivateIP      string
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

	fmt.Printf("Private ip %s, host %s", sc.PrivateIP, sc.Name)
	if sc.PrivateIP != "" {
		viper.Set(sc.confPrivateIPKey(), sc.PrivateIP)
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
	privateIP := viper.GetString(sc.confPrivateIPKey())
	if privateIP != "" {
		sc.PrivateIP = privateIP
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

func (sc *Config) confPrivateIPKey() string {
	return fmt.Sprintf("hcloud.server.%s.privateIP", sc.Name)
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

// AllFromConfig reads Config of all servers from configuration.
func AllFromConfig() ([]Config, error) {
	serverNames := viper.GetStringMapString(confHCloudServersKey)
	if len(serverNames) < 1 {
		return nil, fmt.Errorf("no servers fond in config")
	}

	serverConfigs := make([]Config, 0)
	for name := range serverNames {
		current := FromConfig(name)
		serverConfigs = append(serverConfigs, current)
	}
	return serverConfigs, nil
}

// SetHCloudServerDefaults sets default server type, image and location, which are used to add servers.
func SetHCloudServerDefaults() {
	viper.Set(confHCloudDefaultServerTypeKey, HCloudServerType)
	viper.Set(confHCloudDefaultImageNameKey, HCloudImage)
	viper.Set(confHCloudLocationNameKey, HCloudLocation)
}

// AddServer uses the first argument as server name and adds this server to the configuration.
func AddServer(serverName string) error {
	sshKey, err := sshkey.ReadSSHPublicKeyFromConf()
	if err != nil {
		return err
	}
	if !sshKey.IsProvisioned() {
		return fmt.Errorf(
			"Could not add SSH key to server '%s'. The SSH key '%s' is not available in hcloud. Use the provision command first", serverName, sshKey.Name)
	}
	serverConf := Config{
		Name:           serverName,
		SSHPublicKeyID: sshKey.ID,
		ServerType:     viper.GetString(confHCloudDefaultServerTypeKey),
		ImageName:      viper.GetString(confHCloudDefaultImageNameKey),
		LocationName:   viper.GetString(confHCloudLocationNameKey),
	}
	serverConf.UpdateConfig()
	return nil
}
