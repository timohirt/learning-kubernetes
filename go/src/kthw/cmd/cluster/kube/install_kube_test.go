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
