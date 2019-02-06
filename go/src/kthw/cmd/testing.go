package cmd

import "github.com/hetznercloud/hcloud-go/hcloud"

// MockHCloudOperations mock object to test code which depends on HCloudOperations
type MockHCloudOperations struct {
	createServerResults *CreateServerResults
	createSSHKeyResults *CreateSSHKeyResults
	err                 error
}

// CreateServer returns createServerResults defiend in MockHCloudOperations
func (m *MockHCloudOperations) CreateServer(opts hcloud.ServerCreateOpts) *CreateServerResults {
	return m.createServerResults
}

// CreateSSHKey returns createSSHKeyResults defiend in MockHCloudOperations
func (m *MockHCloudOperations) CreateSSHKey(opts hcloud.SSHKeyCreateOpts) *CreateSSHKeyResults {
	return m.createSSHKeyResults
}

// ASSHPublicKeyWithID fixture to be used in tests
var ASSHPublicKeyWithID = sshPublicKey{id: 17, publicKey: "publicKey", name: "name"}

func ASSHPublicKeyWithIDInConfig() sshPublicKey {
	ASSHPublicKeyWithID.WriteToConfig()
	return ASSHPublicKeyWithID
}
