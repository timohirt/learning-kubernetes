package cmd

import (
	"kthw/cmd/config"
	"kthw/cmd/hcloudclient"
	"reflect"
	"testing"

	viper "github.com/spf13/viper"
)

func setupTestCreateServer() (*hcloudclient.CreateServerResults, *MockHCloudOperations, config.ServerConfig) {
	createServerResult := &hcloudclient.CreateServerResults{
		PublicIP:     "10.0.0.1",
		RootPassword: "Passw0rt",
		DNSName:      "m1.hetzner.com"}
	hcloudClient := &MockHCloudOperations{
		createServerResults: createServerResult}
	config := config.ServerConfig{
		Name:         "m1",
		ServerType:   "cx21",
		ImageName:    "ubuntu",
		LocationName: "nbg1"}
	return createServerResult, hcloudClient, config
}

func TestCreateServer(t *testing.T) {
	viper.Reset()
	sshKey := ASSHPublicKeyWithIDInConfig()
	createServerResult, hcloudClient, serverConfig := setupTestCreateServer()

	updatedConfig, err := createServer(serverConfig, hcloudClient)
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

	_, err := createServer(serverConfig, hcloudClient)
	if err == nil {
		t.Errorf("A error should be returned as there is no SSH public key in config")
	}
}
