package cmd_test

import (
	"kthw/cmd"
	"testing"

	viper "github.com/spf13/viper"
)

func TestAddServer(t *testing.T) {
	viper.Reset()
	cmd.SetHCloudServerDefaults()
	cmd.AddServer(nil, []string{"controller-1"})

	if viper.GetString("hcloud.server.controller-1.name") != "controller-1" {
		t.Error("Server name 'controller-1' doesn't exist in config.")
	}

	serverType := viper.GetString("hcloud.server.controller-1.serverType")
	if serverType != cmd.HCloudServerType {
		t.Errorf("ServerType '%s' in configuration differs from expected default type '%s'.", serverType, cmd.HCloudServerType)
	}

	location := viper.GetString("hcloud.server.controller-1.locationName")
	if location != cmd.HCloudLocation {
		t.Errorf("Location '%s' in configuration differs from expected default location '%s'.", location, cmd.HCloudLocation)
	}

	image := viper.GetString("hcloud.server.controller-1.imageName")
	if image != cmd.HCloudImage {
		t.Errorf("Image '%s' in configuration differs from expected default image '%s'.", image, cmd.HCloudImage)
	}
}

func TestReadServiceConfigFailIfServerNameNotSet(t *testing.T) {
	viper.Reset()

	serverConfig := cmd.ServerConfig{}
	err := serverConfig.ReadFromConfig()
	if err == nil {
		t.Error("Server loaded from config, although no server name was given.")
	}
}

func TestReadInitialServerConfig(t *testing.T) {
	viper.Reset()

	viper.Set("hcloud.server.controller-1.name", "controller-1")
	viper.Set("hcloud.server.controller-1.serverType", "irrelevant")
	viper.Set("hcloud.server.controller-1.locationName", "irrelevant")
	viper.Set("hcloud.server.controller-1.imageName", "irrelevant")

	serverConfig := cmd.ServerConfig{Name: "controller-1"}
	err := serverConfig.ReadFromConfig()
	if err != nil {
		t.Fatal(err)
	}

	if serverConfig.ServerType != "irrelevant" {
		t.Errorf("ServerType was '%s' and differs from expected 'irrelevant'", serverConfig.ServerType)
	}

	if serverConfig.ImageName != "irrelevant" {
		t.Errorf("ImageName was '%s' and differs from expected 'irrelevant'", serverConfig.ImageName)
	}

	if serverConfig.LocationName != "irrelevant" {
		t.Errorf("LocationName was '%s' and differs from expected 'irrelevant'", serverConfig.LocationName)
	}

}

func TestReadServerConfigNonInitValues(t *testing.T) {
	viper.Reset()

	viper.Set("hcloud.server.controller-1.name", "controller-1")
	viper.Set("hcloud.server.controller-1.serverType", "irrelevant")
	viper.Set("hcloud.server.controller-1.locationName", "irrelevant")
	viper.Set("hcloud.server.controller-1.imageName", "irrelevant")
	viper.Set("hcloud.server.controller-1.publicIP", "irrelevant")
	viper.Set("hcloud.server.controller-1.rootPassword", "irrelevant")

	serverConfig := cmd.ServerConfig{Name: "controller-1"}
	err := serverConfig.ReadFromConfig()
	if err != nil {
		t.Fatal(err)
	}

	if serverConfig.PublicIP != "irrelevant" {
		t.Errorf("PublicIP was '%s' and differs from expected 'irrelevant'", serverConfig.ServerType)
	}

	if serverConfig.RootPassword != "irrelevant" {
		t.Errorf("RootPassword was '%s' and differs from expected 'irrelevant'", serverConfig.ImageName)
	}
}
