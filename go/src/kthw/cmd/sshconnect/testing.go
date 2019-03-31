package sshconnect

import (
	"fmt"
	"io"
	"testing"
)

type SSHOperationsMock struct {
	WrittenReadOnlyFiles []ReadOnlyFiles
	RunCmdsCommands      []Command
	SSHOperations
}

type ReadOnlyFiles struct {
	Host           string
	FilePathOnHost string
}

func NewSSHOperationsMock() *SSHOperationsMock {
	return &SSHOperationsMock{}
}

func (s *SSHOperationsMock) RunCmds(commands *Commands) error {
	for _, command := range commands.Commands {
		s.RunCmdsCommands = append(s.RunCmdsCommands, command)
	}
	return nil
}

func (s *SSHOperationsMock) WriteReadOnlyFileTo(host string, contentReader io.Reader, filePathOnHost string) error {
	s.WrittenReadOnlyFiles = append(s.WrittenReadOnlyFiles, ReadOnlyFiles{Host: host, FilePathOnHost: filePathOnHost})
	return nil
}

func (s *SSHOperationsMock) WriteExecutableFileTo(host string, contentReader io.Reader, filePathOnHost string) error {
	return fmt.Errorf("Not implemented")
}

func EnsureNoCommandsIssued(issuedCommands []Command, host string, t *testing.T) {
	for _, issuedCommand := range issuedCommands {
		if issuedCommand.GetHost() == host {
			t.Errorf("No commands for host '%s' expected, but found some.", host)
		}
	}
}

func EnsureCommandIssued(issuedCommands []Command, commandDescription string, host string, t *testing.T) {
	for _, issuedCommand := range issuedCommands {
		if issuedCommand.GetHost() == host && issuedCommand.GetDescription() == commandDescription {
			return
		}
	}
	t.Errorf("Command '%s' was not executed on host '%s'", commandDescription, host)
}

func EnsureCommandNotIssued(issuedCommands []Command, commandDescription string, host string, t *testing.T) {
	for _, issuedCommand := range issuedCommands {
		if issuedCommand.GetHost() == host && issuedCommand.GetDescription() == commandDescription {
			t.Errorf("Command '%s' was executed on host '%s'", commandDescription, host)
		}
	}
}
