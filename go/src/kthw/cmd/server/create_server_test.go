package server_test

import (
	"kthw/cmd/hcloudclient"
	"kthw/cmd/server"
	"kthw/cmd/sshkey"
	"reflect"
	"testing"

	viper "github.com/spf13/viper"
)

func setupTestCreateServer() (*hcloudclient.CreateServerResults, *hcloudclient.MockHCloudOperations, server.Config) {
	createServerResult := &hcloudclient.CreateServerResults{
		PublicIP:     "10.0.0.1",
		RootPassword: "Passw0rt",
		DNSName:      "m1.hetzner.com"}
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

	updatedConfig, err := server.Create(serverConfig, hcloudClient)
	if err != nil {
		t.Errorf("Error while creating server: %s", err)
	}

	serverConfig.RootPassword = createServerResult.RootPassword
	serverConfig.PublicIP = createServerResult.PublicIP
	serverConfig.SSHPublicKeyID = sshKey.ID

	if !reflect.DeepEqual(serverConfig, *updatedConfig) {
		t.Errorf("Expected config differs from actual config")
	}
}

func TestCreateServerWhenThereIsNoSSHPublicKeyInConfig(t *testing.T) {
	viper.Reset()
	_, hcloudClient, serverConfig := setupTestCreateServer()

	_, err := server.Create(serverConfig, hcloudClient)
	if err == nil {
		t.Errorf("A error should be returned as there is no SSH public key in config")
	}
}
