package sshconnect

import (
	"fmt"
	"io"
)

type SSHOperationsMock struct {
	WrittenReadOnlyFiles []ReadOnlyFiles
	RunCmdCommands       []IssuedCommand
	RunCmdsCommands      []Command
	SSHOperations
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
	s.RunCmdCommands = append(s.RunCmdCommands, IssuedCommand{Host: host, Command: command})
	return "", nil
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
