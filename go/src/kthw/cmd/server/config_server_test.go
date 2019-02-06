package server_test

import (
	"kthw/cmd/server"
	"kthw/cmd/sshkey"
	"testing"

	viper "github.com/spf13/viper"
)

func setupConfig(key sshkey.SSHPublicKey) {
	server.SetHCloudServerDefaults()
	key.WriteToConfig()
}

func TestAddServer(t *testing.T) {
	viper.Reset()
	key := sshkey.ASSHPublicKeyWithID
	setupConfig(key)
	server.AddServer("controller-1")

	if viper.GetString("hcloud.server.controller-1.name") != "controller-1" {
		t.Error("Server name 'controller-1' doesn't exist in config.")
	}

	serverType := viper.GetString("hcloud.server.controller-1.serverType")
	if serverType != server.HCloudServerType {
		t.Errorf("ServerType '%s' in configuration differs from expected default type '%s'.", serverType, server.HCloudServerType)
	}

	location := viper.GetString("hcloud.server.controller-1.locationName")
	if location != server.HCloudLocation {
		t.Errorf("Location '%s' in configuration differs from expected default location '%s'.", location, server.HCloudLocation)
	}

	image := viper.GetString("hcloud.server.controller-1.imageName")
	if image != server.HCloudImage {
		t.Errorf("Image '%s' in configuration differs from expected default image '%s'.", image, server.HCloudImage)
	}

	publicKeyID := viper.GetInt("hcloud.server.controller-1.publicKeyId")
	if publicKeyID != key.ID {
		t.Errorf("SSH key id '%d' in config differs from expected key id '%d'", publicKeyID, key.ID)
	}

}

func TestReadServiceConfigFailIfServerNameNotSet(t *testing.T) {
	viper.Reset()

	serverConfig := server.Config{}
	err := serverConfig.ReadFromConfig()
	if err == nil {
		t.Error("Server loaded from config, although no server name was given.")
	}
}

func TestReadInitialConfig(t *testing.T) {
	viper.Reset()

	viper.Set("hcloud.server.controller-1.name", "controller-1")
	viper.Set("hcloud.server.controller-1.serverType", "irrelevant")
	viper.Set("hcloud.server.controller-1.locationName", "irrelevant")
	viper.Set("hcloud.server.controller-1.imageName", "irrelevant")

	serverConfig := server.Config{Name: "controller-1"}
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

func TestReadConfigNonInitValues(t *testing.T) {
	viper.Reset()

	viper.Set("hcloud.server.controller-1.name", "controller-1")
	viper.Set("hcloud.server.controller-1.serverType", "irrelevant")
	viper.Set("hcloud.server.controller-1.locationName", "irrelevant")
	viper.Set("hcloud.server.controller-1.imageName", "irrelevant")
	viper.Set("hcloud.server.controller-1.publicIP", "irrelevant")
	viper.Set("hcloud.server.controller-1.rootPassword", "irrelevant")
	viper.Set("hcloud.server.controller-1.publicKeyId", 17)

	serverConfig := server.Config{Name: "controller-1"}
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

	if serverConfig.SSHPublicKeyID != 17 {
		t.Errorf("SSHPublicKeyID was '%d' and differs from expected '17'", serverConfig.SSHPublicKeyID)
	}
}
