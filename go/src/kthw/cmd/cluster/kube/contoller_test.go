package kube_test

import (
	"kthw/certs"
	"kthw/cmd/cluster/kube"
	"kthw/cmd/infra/server"
	"kthw/cmd/sshconnect"
	"testing"
)

func TestInstallKubernetes(t *testing.T) {
	sshMock := sshconnect.NewSSHOperationsMock()
	certLoaderMock := certs.NewCertificateLoaderMock()
	generatesCerts := certs.NewGeneratesCertsMock()
	controllerPublicIP := "192.168.1.1"
	controllerNode := &kube.ControllerNode{
		Config: &server.Config{
			ID:       1,
			PublicIP: controllerPublicIP,
			Roles:    []string{"controller", "etcd"}}}
	etcdNodes := []*kube.EtcdNode{&kube.EtcdNode{EndpointURL: "irrelevant"}}
	taintController := true

	err := kube.InstallControllerNode(controllerNode, etcdNodes, sshMock, certLoaderMock, generatesCerts, taintController)
	if err != nil {
		t.Errorf("InstallOnHosts returned an unexpected error: %s\n", err)
	}

	sshconnect.EnsureCommandIssued(sshMock.RunCmdsCommands, "Upload etcd client certificate public key to /etc/kubernetes/pki/etcd-client.crt", controllerPublicIP, t)
	sshconnect.EnsureCommandIssued(sshMock.RunCmdsCommands, "Upload etcd client certificate private key to /etc/kubernetes/pki/etcd-client.key", controllerPublicIP, t)
	sshconnect.EnsureCommandIssued(sshMock.RunCmdsCommands, "Upload CA certificate public key to /etc/kubernetes/pki/ca.crt", controllerPublicIP, t)
	sshconnect.EnsureCommandIssued(sshMock.RunCmdsCommands, "Upload CA certificate private key to /etc/kubernetes/pki/ca.key", controllerPublicIP, t)
	sshconnect.EnsureCommandIssued(sshMock.RunCmdsCommands, "Copy kubeadm config", controllerPublicIP, t)
	sshconnect.EnsureCommandIssued(sshMock.RunCmdsCommands, "Install kubernetes cluster", controllerPublicIP, t)
	sshconnect.EnsureCommandIssued(sshMock.RunCmdsCommands, "Setup Kubectl", controllerPublicIP, t)
	sshconnect.EnsureCommandIssued(sshMock.RunCmdsCommands, "Open firewall pod network -> public IP and :6443 -> public IP", controllerPublicIP, t)
	sshconnect.EnsureCommandIssued(sshMock.RunCmdsCommands, "Install Calico networking", controllerPublicIP, t)
	sshconnect.EnsureCommandIssued(sshMock.RunCmdsCommands, "Untaint controller, allow pod scheduling on controller node", controllerPublicIP, t)
}

func TestDoNotUntaintController(t *testing.T) {
	sshMock := sshconnect.NewSSHOperationsMock()
	certLoaderMock := certs.NewCertificateLoaderMock()
	generatesCerts := certs.NewGeneratesCertsMock()
	controllerPublicIP := "192.168.1.1"
	controllerNode := &kube.ControllerNode{
		Config: &server.Config{
			ID:       1,
			PublicIP: controllerPublicIP,
			Roles:    []string{"controller", "etcd"}}}

	etcdNodes := []*kube.EtcdNode{&kube.EtcdNode{EndpointURL: "irrelevant"}}
	taintController := false

	err := kube.InstallControllerNode(controllerNode, etcdNodes, sshMock, certLoaderMock, generatesCerts, taintController)
	if err != nil {
		t.Errorf("InstallOnHosts returned an unexpected error: %s\n", err)
	}

	sshconnect.EnsureCommandNotIssued(sshMock.RunCmdsCommands, "Untaint controller, allow pod scheduling on controller node", controllerPublicIP, t)
}
