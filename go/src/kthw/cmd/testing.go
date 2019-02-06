package cmd

import (
	"kthw/cmd/hcloudclient"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// MockHCloudOperations mock object to test code which depends on HCloudOperations
type MockHCloudOperations struct {
	createServerResults *hcloudclient.CreateServerResults
	createSSHKeyResults *hcloudclient.CreateSSHKeyResults
	err                 error
}

// CreateServer returns createServerResults defiend in MockHCloudOperations
func (m *MockHCloudOperations) CreateServer(opts hcloud.ServerCreateOpts) *hcloudclient.CreateServerResults {
	return m.createServerResults
}

// CreateSSHKey returns createSSHKeyResults defiend in MockHCloudOperations
func (m *MockHCloudOperations) CreateSSHKey(opts hcloud.SSHKeyCreateOpts) *hcloudclient.CreateSSHKeyResults {
	return m.createSSHKeyResults
}

// ASSHPublicKeyWithID fixture to be used in tests
var ASSHPublicKeyWithID = sshPublicKey{id: 17, publicKey: "publicKey", name: "name"}

// ASSHPublicKeyWithIDInConfig writes key to config in scope and return sshPublicKey
func ASSHPublicKeyWithIDInConfig() sshPublicKey {
	ASSHPublicKeyWithID.WriteToConfig()
	return ASSHPublicKeyWithID
}
