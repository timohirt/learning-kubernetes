package etcd_test

import (
	"kthw/certs"
	"kthw/cmd/cluster/etcd"
	"kthw/cmd/infra/server"
	"kthw/cmd/sshconnect"
	"testing"
)

func TestFailInstallEtcdIfNoHostsWithRoleEtcdExist(t *testing.T) {
	mock := sshconnect.NewSSHOperationsMock()
	hostConfigs := []*server.Config{
		&server.Config{ID: 1, PublicIP: "192.168.1.2", Roles: []string{"controller"}}}

	generatesCerts := certs.NewGeneratesCertsMock()
	err := etcd.InstallOnHost(hostConfigs, mock, generatesCerts)

	if err == nil {
		t.Errorf("Installing etcd if there is no host with role etcd is not possible.\n")
	}
}

func TestFailInstallEtcdIfNoSingleNodeCluster(t *testing.T) {
	mock := sshconnect.NewSSHOperationsMock()
	hostConfigs := []*server.Config{
		&server.Config{ID: 1, PublicIP: "192.168.1.1", Roles: []string{"etcd"}},
		&server.Config{ID: 2, PublicIP: "192.168.1.2", Roles: []string{"etcd"}}}

	err := etcd.InstallOnHost(hostConfigs, mock, certs.NewGeneratesCertsMock())

	if err == nil {
		t.Errorf("Installing etcd if there is no host with role etcd is not possible.\n")
	}
}

func TestInstallEtcd(t *testing.T) {
	mock := sshconnect.NewSSHOperationsMock()
	generatesCerts := certs.NewGeneratesCertsMock()
	hostInEtcdRole := &server.Config{ID: 1, PublicIP: "192.168.1.1", Roles: []string{"etcd", "worker"}}
	hostConfigs := []*server.Config{
		hostInEtcdRole,
		&server.Config{ID: 2, PublicIP: "192.168.1.2", Roles: []string{"controller"}}}

	err := etcd.InstallOnHost(hostConfigs, mock, generatesCerts)
	if err != nil {
		t.Errorf("InstallEtcd returned an unexpected error: %s\n", err)
	}

	if !generatesCerts.IsEtcdCertGenerated {
		t.Errorf("etcd certificate was not generated\n")
	}

	sshconnect.EnsureCommandIssued(mock.RunCmdsCommands, "Download etcd binary", hostInEtcdRole.PublicIP, t)
	sshconnect.EnsureCommandIssued(mock.RunCmdsCommands, "Upload etcd certificate public key to /etc/etcd/pki/etcd.crt", hostInEtcdRole.PublicIP, t)
	sshconnect.EnsureCommandIssued(mock.RunCmdsCommands, "Upload etcd certificate private key to /etc/etcd/pki/etcd.key", hostInEtcdRole.PublicIP, t)
	sshconnect.EnsureCommandIssued(mock.RunCmdsCommands, "Upload CA certificate public key to /etc/etcd/pki/ca.crt", hostInEtcdRole.PublicIP, t)
	sshconnect.EnsureCommandIssued(mock.RunCmdsCommands, "Untar etcd archive and copy to /usr/local/bin", hostInEtcdRole.PublicIP, t)
	sshconnect.EnsureCommandIssued(mock.RunCmdsCommands, "Copy etcd systemd service to host", hostInEtcdRole.PublicIP, t)
	sshconnect.EnsureCommandIssued(mock.RunCmdsCommands, "Enable and start etcd service", hostInEtcdRole.PublicIP, t)

	sshconnect.EnsureNoCommandsIssued(mock.RunCmdsCommands, hostConfigs[1].PublicIP, t)
}
