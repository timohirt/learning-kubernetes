package cmd

import (
	"kthw/cmd/hcloudclient"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// CreateSSHKey creates a SSH key in hcloud
func createSSHKey(key sshPublicKey, hcloudClient hcloudclient.HCloudOperations) *sshPublicKey {
	opts := hcloud.SSHKeyCreateOpts{
		Name:      key.name,
		PublicKey: key.publicKey}
	result := hcloudClient.CreateSSHKey(opts)
	key.id = result.ID
	return &key
}
