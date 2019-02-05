package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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

// ServerConfig from config file
type ServerConfig struct {
	Name         string
	ServerType   string
	ImageName    string
	LocationName string
	PublicIP     string
	RootPassword string
}

// UpdateConfig updates the configuration with the current field values. Changes are not persisted.
func (sc *ServerConfig) UpdateConfig() {
	viper.Set(sc.confServerNameKey(), sc.Name)
	viper.Set(sc.confServerTypeKey(), sc.ServerType)
	viper.Set(sc.confLocationNameKey(), sc.LocationName)
	viper.Set(sc.confImageNameKey(), sc.ImageName)

	if sc.PublicIP != "" {
		viper.Set(sc.confPublicIPKey(), sc.PublicIP)
	}

	if sc.RootPassword != "" {
		viper.Set(sc.confRootPasswordKey(), sc.RootPassword)
	}
}

// ReadFromConfig reads the config of a server from the configuration file.
// Name field of ServerConfig must be set.
func (sc *ServerConfig) ReadFromConfig() error {
	if sc.Name == "" {
		return fmt.Errorf("Could not read server from config. Server name not set")
	}

	publicIP := viper.GetString(sc.confPublicIPKey())
	if publicIP != "" {
		sc.PublicIP = publicIP
	}
	rootPassword := viper.GetString(sc.confRootPasswordKey())
	if rootPassword != "" {
		sc.RootPassword = rootPassword
	}

	sc.ServerType = viper.GetString(sc.confServerTypeKey())
	sc.ImageName = viper.GetString(sc.confImageNameKey())
	sc.LocationName = viper.GetString(sc.confLocationNameKey())
	return nil
}

func (sc *ServerConfig) confServerNameKey() string {
	return fmt.Sprintf("hcloud.server.%s.name", sc.Name)
}

func (sc *ServerConfig) confServerTypeKey() string {
	return fmt.Sprintf("hcloud.server.%s.serverType", sc.Name)
}

func (sc *ServerConfig) confImageNameKey() string {
	return fmt.Sprintf("hcloud.server.%s.imageName", sc.Name)
}

func (sc *ServerConfig) confLocationNameKey() string {
	return fmt.Sprintf("hcloud.server.%s.locationName", sc.Name)
}

func (sc *ServerConfig) confPublicIPKey() string {
	return fmt.Sprintf("hcloud.server.%s.publicIP", sc.Name)
}

func (sc *ServerConfig) confRootPasswordKey() string {
	return fmt.Sprintf("hcloud.server.%s.rootPassword", sc.Name)
}

func serverConfigFromConfig(serverName string) ServerConfig {
	serverConfig := ServerConfig{Name: serverName}
	err := serverConfig.ReadFromConfig()
	whenErrPrintAndExit(err)
	return serverConfig
}

// SetHCloudServerDefaults sets default server type, image and location, which are used to add servers.
func SetHCloudServerDefaults() {
	viper.Set(confHCloudDefaultServerTypeKey, HCloudServerType)
	viper.Set(confHCloudDefaultImageNameKey, HCloudImage)
	viper.Set(confHCloudLocationNameKey, HCloudLocation)
}

// AddServer uses the first argument as server name and adds this server to the configuration.
func AddServer(cmd *cobra.Command, args []string) {
	serverName := args[0]
	serverConf := ServerConfig{
		Name:         serverName,
		ServerType:   viper.GetString(confHCloudDefaultServerTypeKey),
		ImageName:    viper.GetString(confHCloudDefaultImageNameKey),
		LocationName: viper.GetString(confHCloudLocationNameKey),
	}
	serverConf.UpdateConfig()
}

var addServerCommand = &cobra.Command{
	Use:   "add-server <name>",
	Short: "Adds a new server to the config file.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		AddServer(cmd, args)
		err := viper.WriteConfig()
		whenErrPrintAndExit(err)
	}}
