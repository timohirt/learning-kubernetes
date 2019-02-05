package cmd_test

import (
	"kthw/cmd"
	"testing"

	"github.com/cloudflare/cfssl/log"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

type MockHCloudClient struct {
	createServerResults *cmd.CreateServerResults
	err                 error
}

func (m *MockHCloudClient) CreateServer(opts hcloud.ServerCreateOpts) (*cmd.CreateServerResults, error) {
	log.Info("Client: ", m, " opts: ", opts)
	return m.createServerResults, m.err
}

func TestCreateServer(t *testing.T) {
	createServerResult := &cmd.CreateServerResults{
		PublicIP:     "10.0.0.1",
		RootPassword: "Passw0rt",
		DNSName:      "m1.hetzner.com"}
	hcloudClient := &MockHCloudClient{
		createServerResults: createServerResult}

	config := cmd.ServerConfig{
		Name:         "m1",
		ServerType:   "cx21",
		ImageName:    "ubuntu",
		LocationName: "nbg1"}

	updatedConfig, _ := cmd.CreateServer(config, hcloudClient)

	config.RootPassword = createServerResult.RootPassword
	config.PublicIP = createServerResult.PublicIP

	if config != *updatedConfig {
		t.Errorf("Expected config differs from actual config")
	}

}
