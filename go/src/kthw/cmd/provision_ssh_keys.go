package cmd

import (
	"kthw/cmd/hcloudclient"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// CreateSSHKey creates a SSH key in hcloud
func createSSHKey(key SSHPublicKey, hcloudClient hcloudclient.HCloudOperations) *SSHPublicKey {
	opts := hcloud.SSHKeyCreateOpts{
		Name:      key.name,
		PublicKey: key.publicKey}
	result := hcloudClient.CreateSSHKey(opts)
	key.id = result.ID
	return &key
}
