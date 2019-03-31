package kube_test

import (
	"kthw/certs"
	"kthw/cmd/cluster/kube"
	"kthw/cmd/infra/server"
	"kthw/cmd/sshconnect"
	"testing"
)

func TestFailInstallIfNoHostsWithRoleControllerExist(t *testing.T) {
	mock := sshconnect.NewSSHOperationsMock()
	certLoaderMock := certs.NewCertificateLoaderMock()
	generatesCerts := certs.NewGeneratesCertsMock()
	hostConfigs := []*server.Config{
		&server.Config{ID: 1, PublicIP: "192.168.1.2", Roles: []string{"etcd"}}}

	err := kube.InstallOnHosts(hostConfigs, mock, certLoaderMock, generatesCerts)

	if err == nil {
		t.Errorf("Installing kubernetes if there is no host with role controller is not possible.\n")
	}
}

func TestFailInstallIfMoreThanOneHostWithRoleControllerExist(t *testing.T) {
	mock := sshconnect.NewSSHOperationsMock()
	certLoaderMock := certs.NewCertificateLoaderMock()
	generatesCerts := certs.NewGeneratesCertsMock()
	hostConfigs := []*server.Config{
		&server.Config{ID: 1, PublicIP: "192.168.1.1", Roles: []string{"controller"}},
		&server.Config{ID: 2, PublicIP: "192.168.1.2", Roles: []string{"controller"}}}

	err := kube.InstallOnHosts(hostConfigs, mock, certLoaderMock, generatesCerts)

	if err == nil {
		t.Errorf("Installing kubernetes is currently only supported with on controller.\n")
	}
}

func TestInstallKubernetes(t *testing.T) {
	sshMock := sshconnect.NewSSHOperationsMock()
	certLoaderMock := certs.NewCertificateLoaderMock()
	generatesCerts := certs.NewGeneratesCertsMock()
	hostInControllerRole := &server.Config{ID: 1, PublicIP: "192.168.1.1", Roles: []string{"controller", "etcd"}}
	hostConfigs := []*server.Config{
		hostInControllerRole}

	err := kube.InstallOnHosts(hostConfigs, sshMock, certLoaderMock, generatesCerts)
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
	sshconnect.EnsureCommandIssued(sshMock.RunCmdsCommands, "Untaint master, allow pod scheduling on master node", hostInControllerRole.PublicIP, t)
}
