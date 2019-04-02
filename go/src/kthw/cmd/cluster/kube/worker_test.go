package kube_test

import (
	"kthw/cmd/cluster/kube"
	"kthw/cmd/infra/server"
	"kthw/cmd/sshconnect"
	"testing"
)

func TestFailIfHostIsInControllerRole(t *testing.T) {
	sshMock := sshconnect.NewSSHOperationsMock()
	hostConfig := &server.Config{ID: 1, PublicIP: "192.168.1.1", Roles: []string{"controller"}}
	controllerNode := &kube.ControllerNode{}

	err := kube.InstallWorkerNode(hostConfig, controllerNode, sshMock)
	if err == nil {
		t.Errorf("Installing worker nodes is only possible on worker nodes and not on servers in role controller.")
	}
}

func TestFailIfHostIsNotInWorkerNode(t *testing.T) {
	sshMock := sshconnect.NewSSHOperationsMock()
	hostConfig := &server.Config{ID: 1, PublicIP: "192.168.1.1", Roles: []string{"etcd"}}
	controllerNode := &kube.ControllerNode{}

	err := kube.InstallWorkerNode(hostConfig, controllerNode, sshMock)
	if err == nil {
		t.Errorf("Installing worker requires a server to be in role worker.")
	}
}
