package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// HetznerClient talks to the hcloud API
type HetznerClient struct {
	apiToken string
}

// NewHetznerClient creates a new HetznerClient using the APIToken provided as a flag
func NewHetznerClient() *HetznerClient {
	fmt.Println("apitoken: ", APIToken)
	return &HetznerClient{apiToken: APIToken}
}

func (hc *HetznerClient) getServerByName(name string) (*hcloud.Server, error) {
	fmt.Println("Dd ", hc.apiToken)
	client := hcloud.NewClient(hcloud.WithToken(hc.apiToken))

	server, _, err := client.Server.GetByName(context.Background(), name)

	if err != nil {
		log.Println("Error calling Hetzner API: ", err, " (Token: ", hc.apiToken, ")")
		log.Fatal(1)
	}
	if server == nil {
		return nil, fmt.Errorf("Server %s doesn't exist", name)
	}
	return server, nil
}
