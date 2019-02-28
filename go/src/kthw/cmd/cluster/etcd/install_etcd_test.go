package etcd_test

import (
	"kthw/cmd/cluster/etcd"
	"kthw/cmd/infra/server"
	"kthw/cmd/sshconnect"
	"testing"
)

func TestFailInstallEtcdIfNoHostsWithRoleEtcdExist(t *testing.T) {
	mock := sshconnect.NewSSHOperationsMock()
	hostConfigs := []server.Config{
		server.Config{ID: 1, PublicIP: "192.168.1.2", Roles: []string{"controller"}}}

	err := etcd.InstallOnHost(hostConfigs, mock)

	if err == nil {
		t.Errorf("Installing etcd if there is no host with role etcd is not possible.\n")
	}
}

func TestFailInstallEtcdIfNoSingleNodeCluster(t *testing.T) {
	mock := sshconnect.NewSSHOperationsMock()
	hostConfigs := []server.Config{
		server.Config{ID: 1, PublicIP: "192.168.1.1", Roles: []string{"etcd"}},
		server.Config{ID: 2, PublicIP: "192.168.1.2", Roles: []string{"etcd"}}}

	err := etcd.InstallOnHost(hostConfigs, mock)

	if err == nil {
		t.Errorf("Installing etcd if there is no host with role etcd is not possible.\n")
	}
}

func TestInstallEtcd(t *testing.T) {
	mock := sshconnect.NewSSHOperationsMock()
	hostInEtcdRole := server.Config{ID: 1, PublicIP: "192.168.1.1", Roles: []string{"etcd", "worker"}}
	hostConfigs := []server.Config{
		hostInEtcdRole,
		server.Config{ID: 2, PublicIP: "192.168.1.2", Roles: []string{"controller"}}}

	err := etcd.InstallOnHost(hostConfigs, mock)
	if err != nil {
		t.Errorf("InstallEtcd returned an unexpected error: %s\n", err)
	}

	ensureCommandIssued(mock.RunCmdsCommands, "Download etcd binary", hostInEtcdRole.PublicIP, t)
	ensureCommandIssued(mock.RunCmdsCommands, "Untar etcd archive and copy to /usr/local/bin", hostInEtcdRole.PublicIP, t)
	ensureCommandIssued(mock.RunCmdsCommands, "Copy etcd systemd service to host", hostInEtcdRole.PublicIP, t)

	ensureNoCommandsIssued(mock.RunCmdsCommands, hostConfigs[1].PublicIP, t)
}

func ensureNoCommandsIssued(issuedCommands []sshconnect.Command, host string, t *testing.T) {
	for _, issuedCommand := range issuedCommands {
		if issuedCommand.GetHost() == host {
			t.Errorf("No commands for host '%s' expected, but found some.", host)
		}
	}
}

func ensureCommandIssued(issuedCommands []sshconnect.Command, commandDescription string, host string, t *testing.T) {
	for _, issuedCommand := range issuedCommands {
		if issuedCommand.GetHost() == host && issuedCommand.GetDescription() == commandDescription {
			return
		}
	}
	t.Errorf("Command '%s' was not executed on host '%s'", commandDescription, host)
}
