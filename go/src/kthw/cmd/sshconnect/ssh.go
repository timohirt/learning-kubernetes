package sshconnect

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"kthw/cmd/common"
	"os"
	"path"

	"github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
)

const (
	ed25519Key = "id_ed25519"
	rsaKey     = "id_rsa"
	sshBaseDir = ".ssh"
)

func loadPrivateKeyFile() ssh.AuthMethod {
	userHome := os.Getenv("HOME")
	var selectedKeyFile string
	ed25519KeyPath := path.Join(userHome, sshBaseDir, ed25519Key)
	rsaKeyPath := path.Join(userHome, sshBaseDir, rsaKey)
	if common.FileExists(ed25519KeyPath) {
		selectedKeyFile = ed25519KeyPath
	} else if common.FileExists(rsaKeyPath) {
		selectedKeyFile = rsaKeyPath
	} else {
		fmt.Println("No supported SSH private key found!")
		os.Exit(1)
	}

	buffer, err := ioutil.ReadFile(selectedKeyFile)
	if err != nil {
		fmt.Printf("Error while reading SSH private key from file '%s': %s", selectedKeyFile, err)
		os.Exit(1)
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		fmt.Printf("Error while parsing SSH private key from file '%s': %s", selectedKeyFile, err)
		os.Exit(1)
	}
	return ssh.PublicKeys(key)
}

// SSHOperations allow running commands on remote hosts and transferring files.
type SSHOperations interface {
	RunCmd(host string, command string) (string, error)
	WriteReadOnlyFileTo(host string, contentReader io.Reader, filePathOnHost string) error
	WriteExecutableFileTo(host string, contentReader io.Reader, filePathOnHost string) error
}

// SSHConnect contains sshConfig used to connect to hosts and allows to run commands on a host and copy files via SCP.
type SSHConnect struct {
	sshConfig ssh.ClientConfig
	SSHOperations
}

// NewSSHConnect created a ssh.ClintConfig for user root and using a private key from ~/.ssh
func NewSSHConnect() *SSHConnect {
	sshConfig := ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{loadPrivateKeyFile()},
		HostKeyCallback: ssh.InsecureIgnoreHostKey()}
	return &SSHConnect{sshConfig: sshConfig}
}

func (c *SSHConnect) connect(host string) (*ssh.Session, error) {
	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", host), &c.sshConfig)
	if err != nil {
		return nil, fmt.Errorf("Error while connecting to server: %s", err)
	}

	session, err := connection.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Failed to create session: %s", err)
	}
	return session, nil
}

// RunCmd connects to host, runs command on this host and returns its output.
func (c *SSHConnect) RunCmd(host string, command string) (string, error) {
	session, err := c.connect(host)
	if err != nil {
		return "", err
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf
	err = session.Run(command)
	if err != nil {
		errFromServer := stderrBuf.String()
		return "", fmt.Errorf("error while running command on remote host. Error output %s. Error %s", errFromServer, err)
	}

	return stdoutBuf.String(), nil
}

// WriteReadOnlyFileTo connects to host, reads from contentReader and writes it to file at filePathOnHost.
// Set permission of this file to 0444.
func (c *SSHConnect) WriteReadOnlyFileTo(host string, contentReader io.Reader, filePathOnHost string) error {
	return c.writeFileTo(host, contentReader, filePathOnHost, "0444")
}

// WriteExecutableFileTo connects to host, reads from contentReader and writes it to file at filePathOnHost.
// Set permission of this file to 0744.
func (c *SSHConnect) WriteExecutableFileTo(host string, contentReader io.Reader, filePathOnHost string) error {
	return c.writeFileTo(host, contentReader, filePathOnHost, "0744")
}

func (c *SSHConnect) writeFileTo(host string, contentReader io.Reader, filePathOnHost string, filePermission string) error {

	client := scp.NewClient(fmt.Sprintf("%s:22", host), &c.sshConfig)
	err := client.Connect()
	if err != nil {
		return err
	}
	defer client.Close()

	return client.CopyFile(contentReader, filePathOnHost, "0655")
}
