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

	if len(mock.WrittenReadOnlyFiles) != 2 {
		t.Errorf("Expected two files to be written, but '%d' were written.", len(mock.WrittenReadOnlyFiles))
	}

	expectedFilePath := "/etc/wireguard/wg0.conf"
	for _, writtenReadOnlyFile := range mock.WrittenReadOnlyFiles {
		if writtenReadOnlyFile.FilePathOnHost != expectedFilePath {
			t.Errorf("Expected file path '%s', but was '%s' for host '%s'.", expectedFilePath, writtenReadOnlyFile.FilePathOnHost, writtenReadOnlyFile.Host)
		}
	}

	if len(mock.RunCmdCommands) != 4 {
		t.Errorf("Expected 4 issued commands during setup, but only '%d' were issued.", len(mock.RunCmdCommands))
	}

	expectedStartWireguardCommand := "systemctl enable wg-quick@wg0 && systemctl restart wg-quick@wg0"
	expectedUfwCommand := "ufw allow in on wg0"
	for _, issuedCommand := range mock.RunCmdCommands {
		if issuedCommand.Command != expectedStartWireguardCommand && issuedCommand.Command != expectedUfwCommand {
			t.Errorf("Unexpected command '%s' issued on host '%s'", issuedCommand.Command, issuedCommand.Host)
		}
	}
}
