package sshkey_test

import (
	"kthw/cmd/sshkey"
	"testing"

	viper "github.com/spf13/viper"
)

func TestWriteSSHPublicKeyToConf(t *testing.T) {
	viper.Reset()

	key := &sshkey.SSHPublicKey{PublicKey: "key", Name: "name"}

	key.WriteToConfig()

	publicKeyFromConfig := viper.GetString("sshKeys.publicKey")
	if publicKeyFromConfig != "key" {
		t.Errorf("Public key from config '%s' differs from expedted public key 'key'", publicKeyFromConfig)
	}

	nameFromConfig := viper.GetString("sshKeys.name")
	if nameFromConfig != "name" {
		t.Errorf("Name from config '%s' differs from expected name 'name'", nameFromConfig)
	}

	key.ID = 12
	key.WriteToConfig()

	keyFromConfig := viper.GetInt("sshKeys.id")
	if keyFromConfig != 12 {
		t.Errorf("SSH key id from config '%d' differs from expected ID '12'", keyFromConfig)
	}
}

func TestReadSSHPublicKeyFromConf(t *testing.T) {
	viper.Reset()
	viper.Set("sshKeys.name", "name")
	viper.Set("sshKeys.publicKey", "key")

	key, _ := sshkey.ReadSSHPublicKeyFromConf()
	if key.Name != "name" || key.PublicKey != "key" {
		t.Error("Could not read public key from conf")
	}

	viper.Set("sshKeys.id", 12)
	key, _ = sshkey.ReadSSHPublicKeyFromConf()

	if key.ID != 12 {
		t.Error("Error reading ssh key id from config")
	}
}
