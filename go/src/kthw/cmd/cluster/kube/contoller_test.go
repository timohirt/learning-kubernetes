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
	hostInControllerRole := &server.Config{ID: 1, PublicIP: "192.168.1.1", Roles: []string{"controller", "etcd"}}

	etcdNodes := []*kube.EtcdNode{&kube.EtcdNode{EndpointURL: "irrelevant"}}
	taintController := true

	_, err := kube.InstallControllerNode(hostInControllerRole, etcdNodes, sshMock, certLoaderMock, generatesCerts, taintController)
	if err != nil {
		t.Errorf("InstallOnHosts returned an unexpected error: %s\n", err)
	}

	sshconnect.EnsureCommandIssued(sshMock.RunCmdsCommands, "Upload etcd client certificate public key to /etc/kubernetes/pki/etcd-client.crt", hostInControllerRole.PublicIP, t)
	sshconnect.EnsureCommandIssued(sshMock.RunCmdsCommands, "Upload etcd client certificate private key to /etc/kubernetes/pki/etcd-client.key", hostInControllerRole.PublicIP, t)
	sshconnect.EnsureCommandIssued(sshMock.RunCmdsCommands, "Upload CA certificate public key to /etc/kubernetes/pki/ca.crt", hostInControllerRole.PublicIP, t)
	sshconnect.EnsureCommandIssued(sshMock.RunCmdsCommands, "Upload CA certificate private key to /etc/kubernetes/pki/ca.key", hostInControllerRole.PublicIP, t)
	sshconnect.EnsureCommandIssued(sshMock.RunCmdsCommands, "Copy kubeadm config", hostInControllerRole.PublicIP, t)
	sshconnect.EnsureCommandIssued(sshMock.RunCmdsCommands, "Install kubernetes cluster", hostInControllerRole.PublicIP, t)
	sshconnect.EnsureCommandIssued(sshMock.RunCmdsCommands, "Setup Kubectl", hostInControllerRole.PublicIP, t)
	sshconnect.EnsureCommandIssued(sshMock.RunCmdsCommands, "Open firewall pod network -> public IP and :6443 -> public IP", hostInControllerRole.PublicIP, t)
	sshconnect.EnsureCommandIssued(sshMock.RunCmdsCommands, "Install Calico networking", hostInControllerRole.PublicIP, t)
	sshconnect.EnsureCommandIssued(sshMock.RunCmdsCommands, "Untaint controller, allow pod scheduling on controller node", hostInControllerRole.PublicIP, t)
}

func TestDoNotUntaintController(t *testing.T) {
	sshMock := sshconnect.NewSSHOperationsMock()
	certLoaderMock := certs.NewCertificateLoaderMock()
	generatesCerts := certs.NewGeneratesCertsMock()
	hostInControllerRole := &server.Config{ID: 1, PublicIP: "192.168.1.1", Roles: []string{"controller", "etcd"}}

	etcdNodes := []*kube.EtcdNode{&kube.EtcdNode{EndpointURL: "irrelevant"}}
	taintController := false

	_, err := kube.InstallControllerNode(hostInControllerRole, etcdNodes, sshMock, certLoaderMock, generatesCerts, taintController)
	if err != nil {
		t.Errorf("InstallOnHosts returned an unexpected error: %s\n", err)
	}

	sshconnect.EnsureCommandNotIssued(sshMock.RunCmdsCommands, "Untaint controller, allow pod scheduling on controller node", hostInControllerRole.PublicIP, t)
}
