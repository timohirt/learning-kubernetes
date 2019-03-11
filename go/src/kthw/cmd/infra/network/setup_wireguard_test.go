package network_test

import (
	"kthw/cmd/infra/network"
	"kthw/cmd/infra/server"
	"kthw/cmd/sshconnect"
	"testing"
)

func TestSetupWireGuard(t *testing.T) {
	mock := sshconnect.NewSSHOperationsMock()
	hostConfigs := []server.Config{
		server.Config{ID: 1, PublicIP: "192.168.1.1"},
		server.Config{ID: 2, PublicIP: "192.168.1.2"}}

	network.SetupWireguard(mock, hostConfigs)

	sshconnect.EnsureCommandIssued(mock.RunCmdsCommands, "Upload wireguard config file of device 'wg0'", hostConfigs[0].PublicIP, t)
	sshconnect.EnsureCommandIssued(mock.RunCmdsCommands, "Open firewall for private overlay network", hostConfigs[0].PublicIP, t)
	sshconnect.EnsureCommandIssued(mock.RunCmdsCommands, "Start wireguard device 'wg0'", hostConfigs[0].PublicIP, t)

	sshconnect.EnsureCommandIssued(mock.RunCmdsCommands, "Upload wireguard config file of device 'wg0'", hostConfigs[1].PublicIP, t)
	sshconnect.EnsureCommandIssued(mock.RunCmdsCommands, "Open firewall for private overlay network", hostConfigs[1].PublicIP, t)
	sshconnect.EnsureCommandIssued(mock.RunCmdsCommands, "Start wireguard device 'wg0'", hostConfigs[1].PublicIP, t)
}
