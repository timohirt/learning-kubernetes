package sshkey_test

import (
	"kthw/cmd/infra/sshkey"
	"testing"
)

func TestAddSSHPublicKeyToConfig(t *testing.T) {
	name := "my_ssh_key"
	file := "testdata/id_rsa.pub"
	_, err := sshkey.AddSSHPublicKeyToConfig(name, file)
	if err != nil {
		t.Fatalf("Error while parsing SSH public key: %s", err)
	}

	publicKeyFromConf, err := sshkey.ReadSSHPublicKeyFromConf()
	if err != nil {
		t.Fatalf("Error while loading public key from config: %s", err)
	}

	expectedPublicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCySPRQod61J1swABdriGr5m0gB testuser"
	if publicKeyFromConf.PublicKey != expectedPublicKey {
		t.Errorf("Public key from file '%s' differs from expected public key '%s'", publicKeyFromConf.PublicKey, expectedPublicKey)
	}

	if publicKeyFromConf.Name != name {
		t.Errorf("Name of public key from file '%s' differs from expected name '%s'", publicKeyFromConf.Name, name)
	}
}
