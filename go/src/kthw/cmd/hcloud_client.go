package cmd

import (
	"context"
	"log"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// HCloudOperations defines operations to be implemented by HCloudClient
type HCloudOperations interface {
	CreateServer(opts hcloud.ServerCreateOpts) (*CreateServerResults, error)
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

// CreateServerResults groups data from server
type CreateServerResults struct {
	PublicIP     string
	RootPassword string
	DNSName      string
}

// Creates a server using hcloud API and the provided options
func (hc *HCloudClient) CreateServer(opts hcloud.ServerCreateOpts) (*CreateServerResults, error) {
	serverCreateResult, _, err := hc.client.Server.Create(hc.context, opts)
	hc.ensureNoError(err)
	return &CreateServerResults{
		PublicIP:     serverCreateResult.Server.PublicNet.IPv4.IP.String(),
		RootPassword: serverCreateResult.RootPassword,
		DNSName:      serverCreateResult.Server.PublicNet.IPv4.DNSPtr}, nil
}

func (hc *HCloudClient) ensureNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
