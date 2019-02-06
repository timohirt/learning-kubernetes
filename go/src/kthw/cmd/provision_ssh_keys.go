package cmd

import (
	"kthw/cmd/common"
	"kthw/cmd/hcloudclient"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// CreateSSHKey creates a SSH key in hcloud
func createSSHKey(key common.SSHPublicKey, hcloudClient hcloudclient.HCloudOperations) *common.SSHPublicKey {
	opts := hcloud.SSHKeyCreateOpts{
		Name:      key.Name,
		PublicKey: key.PublicKey}
	result := hcloudClient.CreateSSHKey(opts)
	key.ID = result.ID
	return &key
}
