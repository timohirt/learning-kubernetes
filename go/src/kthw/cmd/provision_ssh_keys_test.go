package cmd

import (
	"kthw/cmd/common"
	"kthw/cmd/hcloudclient"
	"reflect"
	"testing"
)

func TestProvisionSSHKeys(t *testing.T) {
	createSSHKeyResult := &hcloudclient.CreateSSHKeyResults{ID: 12}
	hcloudClient := &MockHCloudOperations{
		createSSHKeyResults: createSSHKeyResult}

	key := common.SSHPublicKey{PublicKey: "key", Name: "name"}
	updatedKey := createSSHKey(key, hcloudClient)

	key.ID = 12

	if !reflect.DeepEqual(*updatedKey, key) {
		t.Errorf("Updated key didn't match expected key.")
	}
}
