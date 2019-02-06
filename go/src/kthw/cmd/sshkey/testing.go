package sshkey

// ASSHPublicKeyWithID fixture to be used in tests
var ASSHPublicKeyWithID = SSHPublicKey{ID: 17, PublicKey: "publicKey", Name: "name"}

// ASSHPublicKeyWithIDInConfig writes key to config in scope and return SSHPublicKey
func ASSHPublicKeyWithIDInConfig() SSHPublicKey {
	ASSHPublicKeyWithID.WriteToConfig()
	return ASSHPublicKeyWithID
}
