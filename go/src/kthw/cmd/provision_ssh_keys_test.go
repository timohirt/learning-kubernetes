package cmd

import (
	"kthw/cmd/hcloudclient"
	"reflect"
	"testing"
)

func TestProvisionSSHKeys(t *testing.T) {
	createSSHKeyResult := &hcloudclient.CreateSSHKeyResults{ID: 12}
	hcloudClient := &MockHCloudOperations{
		createSSHKeyResults: createSSHKeyResult}

	key := sshPublicKey{publicKey: "key", name: "name"}
	updatedKey := createSSHKey(key, hcloudClient)

	key.id = 12

	if !reflect.DeepEqual(*updatedKey, key) {
		t.Errorf("Updated key didn't match expected key.")
	}
}
