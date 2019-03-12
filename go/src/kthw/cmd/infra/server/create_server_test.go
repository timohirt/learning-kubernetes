package server_test

import (
	"kthw/cmd/hcloudclient"
	"kthw/cmd/infra/server"
	"kthw/cmd/infra/sshkey"
	"testing"

	viper "github.com/spf13/viper"
)

func setupTestCreateServer() (*hcloudclient.CreateServerResults, *hcloudclient.MockHCloudOperations, server.Config) {
	createServerResult := &hcloudclient.CreateServerResults{
		ID:       42,
		PublicIP: "10.0.0.1",
		DNSName:  "m1.hetzner.com"}
	hcloudClient := &hcloudclient.MockHCloudOperations{
		CreateServerResults: createServerResult}
	config := server.Config{
		Name:         "m1",
		ServerType:   "cx21",
		ImageName:    "ubuntu",
		LocationName: "nbg1"}
	return createServerResult, hcloudClient, config
}

func TestCreateServer(t *testing.T) {
	viper.Reset()
	sshKey := sshkey.ASSHPublicKeyWithIDInConfig()
	createServerResult, hcloudClient, serverConfig := setupTestCreateServer()

	err := server.Create(&serverConfig, hcloudClient)
	if err != nil {
		t.Errorf("Error while creating server: %s", err)
	}

	if serverConfig.PublicIP != createServerResult.PublicIP {
		t.Errorf("No public IP set, but expected IP '%s'", createServerResult.PublicIP)
	}
	if serverConfig.SSHPublicKeyID != sshKey.ID {
		t.Errorf("No ssh key id set, but expected ID '%d'", sshKey.ID)
	}

	if serverConfig.ID != createServerResult.ID {
		t.Errorf("No server ID set, but expected ID '%d'", createServerResult.ID)
	}
}

func TestCreateServerWhenThereIsNoSSHPublicKeyInConfig(t *testing.T) {
	viper.Reset()
	_, hcloudClient, serverConfig := setupTestCreateServer()

	err := server.Create(&serverConfig, hcloudClient)
	if err == nil {
		t.Errorf("A error should be returned as there is no SSH public key in config")
	}
}
