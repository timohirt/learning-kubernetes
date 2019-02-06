package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/viper"
)

// sshPublicKey parameters of a public ssh key. 'id' is > 0 if key already created in hcloud
type sshPublicKey struct {
	id        int
	publicKey string
	name      string
}

func (s *sshPublicKey) WriteToConfig() {
	viper.Set(confSSHKeysPublicKeyKey, s.publicKey)
	viper.Set(confSSHKeysNameKey, s.name)
	viper.Set(confSSHKeysIDKey, s.id)
}

const (
	confSSHKeysPublicKeyKey = "sshKeys.publicKey"
	confSSHKeysNameKey      = "sshKeys.name"
	confSSHKeysIDKey        = "sshKeys.id"
)

func parseSSHPublicKey(name string, file string) (*sshPublicKey, error) {
	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("Error reading ssh key from file '%s': %s", file, err)
	}
	publicKeyFromFile := strings.TrimSpace(string(fileContent))
	publicKey := &sshPublicKey{publicKey: publicKeyFromFile, name: name}
	return publicKey, nil
}

func readSSHPublicKeyFromConf() (*sshPublicKey, error) {
	if !viper.IsSet(confSSHKeysNameKey) || !viper.IsSet(confSSHKeysPublicKeyKey) {
		return nil, fmt.Errorf("No ssh keys defined to conf. Add one first")
	}

	key := &sshPublicKey{
		id:        viper.GetInt(confSSHKeysIDKey),
		name:      viper.GetString(confSSHKeysNameKey),
		publicKey: viper.GetString(confSSHKeysPublicKeyKey)}
	return key, nil
}
