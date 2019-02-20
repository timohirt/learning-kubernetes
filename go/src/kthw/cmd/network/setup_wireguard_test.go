package network_test

import (
	"fmt"
	"io"
	"kthw/cmd/network"
	"kthw/cmd/server"
	"kthw/cmd/sshconnect"
	"testing"
)

type SSHOperationsMock struct {
	WrittenReadOnlyFiles []ReadOnlyFiles
	IssuedCommands       []IssuedCommand
	sshconnect.SSHOperations
}

type ReadOnlyFiles struct {
	Host           string
	FilePathOnHost string
}

type IssuedCommand struct {
	Host    string
	Command string
}

func NewSSHOperationsMock() *SSHOperationsMock {
	return &SSHOperationsMock{}
}

func (s *SSHOperationsMock) RunCmd(host string, command string) (string, error) {
	s.IssuedCommands = append(s.IssuedCommands, IssuedCommand{Host: host, Command: command})
	return "", nil
}
func (s *SSHOperationsMock) WriteReadOnlyFileTo(host string, contentReader io.Reader, filePathOnHost string) error {
	s.WrittenReadOnlyFiles = append(s.WrittenReadOnlyFiles, ReadOnlyFiles{Host: host, FilePathOnHost: filePathOnHost})
	return nil
}

func (s *SSHOperationsMock) WriteExecutableFileTo(host string, contentReader io.Reader, filePathOnHost string) error {
	return fmt.Errorf("Not implemented")
}

func TestSetupWireGuard(t *testing.T) {
	mock := NewSSHOperationsMock()
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

	if len(mock.IssuedCommands) != 4 {
		t.Errorf("Expected 4 issued commands during setup, but only '%d' were issued.", len(mock.IssuedCommands))
	}

	expectedStartWireguardCommand := "systemctl enable wg-quick@wg0 && systemctl restart wg-quick@wg0"
	expectedUfwCommand := "ufw allow in on wg0"
	for _, issuedCommand := range mock.IssuedCommands {
		if issuedCommand.Command != expectedStartWireguardCommand && issuedCommand.Command != expectedUfwCommand {
			t.Errorf("Unexpected command '%s' issued on host '%s'", issuedCommand.Command, issuedCommand.Host)
		}
	}
}
