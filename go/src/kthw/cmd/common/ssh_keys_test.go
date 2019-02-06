package common

import (
	"testing"

	viper "github.com/spf13/viper"
)

func TestParseSSHPublicKey(t *testing.T) {
	name := "my_ssh_key"
	file := "testdata/id_rsa.pub"
	SSHPublicKey, err := ParseSSHPublicKey(name, file)
	if err != nil {
		t.Fatalf("Error while parsing SSH public key: %s", err)
	}

	expectedPublicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCySPRQod61J1swABdriGr5m0gB testuser"
	if SSHPublicKey.PublicKey != expectedPublicKey {
		t.Errorf("Public key from file '%s' differs from expected public key '%s'", SSHPublicKey.PublicKey, expectedPublicKey)
	}

	if SSHPublicKey.Name != name {
		t.Errorf("Name of public key from file '%s' differs from expected name '%s'", SSHPublicKey.Name, name)
	}
}

func TestWriteSSHPublicKeyToConf(t *testing.T) {
	viper.Reset()

	key := &SSHPublicKey{PublicKey: "key", Name: "name"}

	key.WriteToConfig()

	publicKeyFromConfig := viper.GetString(confSSHKeysPublicKeyKey)
	if publicKeyFromConfig != "key" {
		t.Errorf("Public key from config '%s' differs from expedted public key 'key'", publicKeyFromConfig)
	}

	nameFromConfig := viper.GetString(confSSHKeysNameKey)
	if nameFromConfig != "name" {
		t.Errorf("Name from config '%s' differs from expected name 'name'", nameFromConfig)
	}

	key.ID = 12
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

	key, _ := ReadSSHPublicKeyFromConf()
	if key.Name != "name" || key.PublicKey != "key" {
		t.Error("Could not read public key from conf")
	}

	viper.Set(confSSHKeysIDKey, 12)
	key, _ = ReadSSHPublicKeyFromConf()

	if key.ID != 12 {
		t.Error("Error reading ssh key id from config")
	}
}
