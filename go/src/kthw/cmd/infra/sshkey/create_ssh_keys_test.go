package sshkey_test

import (
	"kthw/cmd/hcloudclient"
	"kthw/cmd/infra/sshkey"
	"reflect"
	"testing"
)

func TestProvisionSSHKeys(t *testing.T) {
	createSSHKeyResult := &hcloudclient.CreateSSHKeyResults{ID: 12}
	hcloudClient := &hcloudclient.MockHCloudOperations{
		CreateSSHKeyResults: createSSHKeyResult}

	key := sshkey.SSHPublicKey{PublicKey: "key", Name: "name"}
	updatedKey := sshkey.CreateSSHKey(key, hcloudClient)

	key.ID = 12

	if !reflect.DeepEqual(*updatedKey, key) {
		t.Errorf("Updated key didn't match expected key.")
	}
}
