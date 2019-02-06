package hcloudclient

import (
	"github.com/hetznercloud/hcloud-go/hcloud"
)

// MockHCloudOperations mock object to test code which depends on HCloudOperations
type MockHCloudOperations struct {
	CreateServerResults *CreateServerResults
	CreateSSHKeyResults *CreateSSHKeyResults
	Err                 error
}

// CreateServer returns createServerResults defiend in MockHCloudOperations
func (m *MockHCloudOperations) CreateServer(opts hcloud.ServerCreateOpts) *CreateServerResults {
	return m.CreateServerResults
}

// CreateSSHKey returns createSSHKeyResults defiend in MockHCloudOperations
func (m *MockHCloudOperations) CreateSSHKey(opts hcloud.SSHKeyCreateOpts) *CreateSSHKeyResults {
	return m.CreateSSHKeyResults
}
