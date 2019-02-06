package sshkey

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// AddSSHPublicKeyToConfig reads SSH public key from file and adds it to config.
func AddSSHPublicKeyToConfig(name string, file string) error {
	publicSSHKey, err := parseSSHPublicKey(name, file)
	if err != nil {
		return fmt.Errorf("Error while parsing ssh key from file '%s' config: %s", file, err)
	}

	publicSSHKey.WriteToConfig()

	return nil
}

func parseSSHPublicKey(name string, file string) (*SSHPublicKey, error) {
	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("Error reading ssh key from file '%s': %s", file, err)
	}
	publicKeyFromFile := strings.TrimSpace(string(fileContent))
	publicKey := &SSHPublicKey{PublicKey: publicKeyFromFile, Name: name}
	return publicKey, nil
}
