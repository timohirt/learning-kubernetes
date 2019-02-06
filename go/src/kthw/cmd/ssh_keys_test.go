package cmd

import (
	"testing"

	viper "github.com/spf13/viper"
)

func TestParseSSHPublicKey(t *testing.T) {
	name := "my_ssh_key"
	file := "testdata/id_rsa.pub"
	sshPublicKey, err := parseSSHPublicKey(name, file)
	if err != nil {
		t.Fatalf("Error while parsing SSH public key: %s", err)
	}

	expectedPublicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCySPRQod61J1swABdriGr5m0gB testuser"
	if sshPublicKey.publicKey != expectedPublicKey {
		t.Errorf("Public key from file '%s' differs from expected public key '%s'", sshPublicKey.publicKey, expectedPublicKey)
	}

	if sshPublicKey.name != name {
		t.Errorf("Name of public key from file '%s' differs from expected name '%s'", sshPublicKey.name, name)
	}
}

func TestWriteSSHPublicKeyToConf(t *testing.T) {
	viper.Reset()

	key := &sshPublicKey{publicKey: "key", name: "name"}

	key.WriteToConfig()

	publicKeyFromConfig := viper.GetString(confSSHKeysPublicKeyKey)
	if publicKeyFromConfig != "key" {
		t.Errorf("Public key from config '%s' differs from expedted public key 'key'", publicKeyFromConfig)
	}

	nameFromConfig := viper.GetString(confSSHKeysNameKey)
	if nameFromConfig != "name" {
		t.Errorf("Name from config '%s' differs from expected name 'name'", nameFromConfig)
	}

	key.id = 12
	key.WriteToConfig()

	keyFromConfig := viper.GetInt(confSSHKeysIDKey)
	if keyFromConfig != 12 {
		t.Errorf("SSH key id from config '%d' differs from expected ID '12'", keyFromConfig)
	}
}

func TestReadSSHPublicKeyFromConf(t *testing.T) {
	viper.Reset()
	viper.Set(confSSHKeysNameKey, "name")
	viper.Set(confSSHKeysPublicKeyKey, "key")

	key, _ := readSSHPublicKeyFromConf()
	if key.name != "name" || key.publicKey != "key" {
		t.Error("Could not read public key from conf")
	}

	viper.Set(confSSHKeysIDKey, 12)
	key, _ = readSSHPublicKeyFromConf()

	if key.id != 12 {
		t.Error("Error reading ssh key id from config")
	}
}
