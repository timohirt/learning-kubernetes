package sshkey

import (
	"kthw/cmd/hcloudclient"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// CreateSSHKey creates a SSH key in hcloud
func CreateSSHKey(key SSHPublicKey, hcloudClient hcloudclient.HCloudOperations) *SSHPublicKey {
	opts := hcloud.SSHKeyCreateOpts{
		Name:      key.Name,
		PublicKey: key.PublicKey}
	result := hcloudClient.CreateSSHKey(opts)
	key.ID = result.ID
	return &key
}
