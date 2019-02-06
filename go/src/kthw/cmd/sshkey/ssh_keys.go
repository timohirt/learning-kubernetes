package sshkey

import (
	"fmt"

	"github.com/spf13/viper"
)

// SSHPublicKey parameters of a public ssh key. 'id' is > 0 if key already created in hcloud
type SSHPublicKey struct {
	ID        int
	PublicKey string
	Name      string
}

// WriteToConfig writes the state of a key to config without writing the config to disk
func (s *SSHPublicKey) WriteToConfig() {
	viper.Set(confSSHKeysPublicKeyKey, s.PublicKey)
	viper.Set(confSSHKeysNameKey, s.Name)
	viper.Set(confSSHKeysIDKey, s.ID)
}

const (
	confSSHKeysPublicKeyKey = "sshKeys.publicKey"
	confSSHKeysNameKey      = "sshKeys.name"
	confSSHKeysIDKey        = "sshKeys.id"
)

// ReadSSHPublicKeyFromConf reads public ssh key from config and returns error if non is set
func ReadSSHPublicKeyFromConf() (*SSHPublicKey, error) {
	if !viper.IsSet(confSSHKeysNameKey) || !viper.IsSet(confSSHKeysPublicKeyKey) {
		return nil, fmt.Errorf("No ssh keys defined to conf. Add one first")
	}

	key := &SSHPublicKey{
		ID:        viper.GetInt(confSSHKeysIDKey),
		Name:      viper.GetString(confSSHKeysNameKey),
		PublicKey: viper.GetString(confSSHKeysPublicKeyKey)}
	return key, nil
}
