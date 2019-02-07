package hcloudclient

import (
	"context"
	"fmt"
	"os"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// HCloudOperations defines operations to be implemented by HCloudClient
type HCloudOperations interface {
	Create(opts hcloud.ServerCreateOpts) *CreateServerResults
	CreateSSHKey(opts hcloud.SSHKeyCreateOpts) *CreateSSHKeyResults
}

// HCloudClient talks to the hcloud API
type HCloudClient struct {
	client  *hcloud.Client
	context context.Context
}

// NewHCloudClient creates a new HetznerClient using the APIToken provided as a flag
func NewHCloudClient(apiToken string) *HCloudClient {
	if apiToken == "" {
		fmt.Println("APIToken not set. Did you set the --apiToken flag?")
		os.Exit(1)
	}
	client := hcloud.NewClient(hcloud.WithToken(apiToken))

	return &HCloudClient{client: client, context: context.Background()}
}

// CreateServerResults groups returned data from hcloud
type CreateServerResults struct {
	ID           int
	PublicIP     string
	RootPassword string
	DNSName      string
}

// Create creates a server using hcloud API and the provided options
func (hc *HCloudClient) Create(opts hcloud.ServerCreateOpts) *CreateServerResults {
	serverCreateResult, _, err := hc.client.Server.Create(hc.context, opts)
	hc.ensureNoError(err)
	return &CreateServerResults{
		ID:           serverCreateResult.Server.ID,
		PublicIP:     serverCreateResult.Server.PublicNet.IPv4.IP.String(),
		RootPassword: serverCreateResult.RootPassword,
		DNSName:      serverCreateResult.Server.PublicNet.IPv4.DNSPtr}
}

// CreateSSHKeyResults groups returned data from hcloud
type CreateSSHKeyResults struct {
	ID int
}

// CreateSSHKey creates a SSH key in hcloud
func (hc *HCloudClient) CreateSSHKey(opts hcloud.SSHKeyCreateOpts) *CreateSSHKeyResults {
	sshKey, _, err := hc.client.SSHKey.Create(hc.context, opts)
	hc.ensureNoError(err)

	return &CreateSSHKeyResults{
		ID: sshKey.ID}
}

func (hc *HCloudClient) ensureNoError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
