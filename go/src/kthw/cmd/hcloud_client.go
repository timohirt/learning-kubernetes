package cmd

import (
	"context"
	"log"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// HCloudOperations defines operations to be implemented by HCloudClient
type HCloudOperations interface {
	CreateServer(opts hcloud.ServerCreateOpts) *CreateServerResults
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
		log.Fatal("APIToken not set. Did you set the --apiToken flag?")
	}
	client := hcloud.NewClient(hcloud.WithToken(apiToken))

	return &HCloudClient{client: client, context: context.Background()}
}

// CreateServerResults groups returned data from hcloud
type CreateServerResults struct {
	PublicIP     string
	RootPassword string
	DNSName      string
}

// CreateServer creates a server using hcloud API and the provided options
func (hc *HCloudClient) CreateServer(opts hcloud.ServerCreateOpts) *CreateServerResults {
	serverCreateResult, _, err := hc.client.Server.Create(hc.context, opts)
	hc.ensureNoError(err)
	return &CreateServerResults{
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
		log.Fatal(err)
	}
}
