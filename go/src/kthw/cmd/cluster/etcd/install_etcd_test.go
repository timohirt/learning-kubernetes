package etcd_test

import (
	"fmt"
	"kthw/certs"
	"kthw/cmd/cluster/etcd"
	"kthw/cmd/infra/server"
	"kthw/cmd/sshconnect"
	"testing"
)

func TestFailInstallEtcdIfNoHostsWithRoleEtcdExist(t *testing.T) {
	mock := sshconnect.NewSSHOperationsMock()
	hostConfigs := []*server.Config{
		&server.Config{ID: 1, PublicIP: "192.168.1.2", Roles: []string{"controller"}}}

	generatesCerts := NewGeneratesCertsMock()
	err := etcd.InstallOnHost(hostConfigs, mock, generatesCerts)

	if err == nil {
		t.Errorf("Installing etcd if there is no host with role etcd is not possible.\n")
	}
}

func TestFailInstallEtcdIfNoSingleNodeCluster(t *testing.T) {
	mock := sshconnect.NewSSHOperationsMock()
	hostConfigs := []*server.Config{
		&server.Config{ID: 1, PublicIP: "192.168.1.1", Roles: []string{"etcd"}},
		&server.Config{ID: 2, PublicIP: "192.168.1.2", Roles: []string{"etcd"}}}

	err := etcd.InstallOnHost(hostConfigs, mock, NewGeneratesCertsMock())

	if err == nil {
		t.Errorf("Installing etcd if there is no host with role etcd is not possible.\n")
	}
}

func TestInstallEtcd(t *testing.T) {
	mock := sshconnect.NewSSHOperationsMock()
	hostInEtcdRole := &server.Config{ID: 1, PublicIP: "192.168.1.1", Roles: []string{"etcd", "worker"}}
	hostConfigs := []*server.Config{
		hostInEtcdRole,
		&server.Config{ID: 2, PublicIP: "192.168.1.2", Roles: []string{"controller"}}}

	generatesCerts := NewGeneratesCertsMock()
	err := etcd.InstallOnHost(hostConfigs, mock, generatesCerts)
	if err != nil {
		t.Errorf("InstallEtcd returned an unexpected error: %s\n", err)
	}

	if !generatesCerts.isEtcdCertGenerated {
		t.Errorf("etcd certificate was not generated\n")
	}

	sshconnect.EnsureCommandIssued(mock.RunCmdsCommands, "Download etcd binary", hostInEtcdRole.PublicIP, t)
	sshconnect.EnsureCommandIssued(mock.RunCmdsCommands, "Upload etcd certificate public key to /etc/kubernetes/pki/etcd.crt", hostInEtcdRole.PublicIP, t)
	sshconnect.EnsureCommandIssued(mock.RunCmdsCommands, "Upload etcd certificate private key to /etc/kubernetes/pki/etcd.key", hostInEtcdRole.PublicIP, t)
	sshconnect.EnsureCommandIssued(mock.RunCmdsCommands, "Upload CA certificate public key to /etc/kubernetes/pki/ca.crt", hostInEtcdRole.PublicIP, t)
	sshconnect.EnsureCommandIssued(mock.RunCmdsCommands, "Untar etcd archive and copy to /usr/local/bin", hostInEtcdRole.PublicIP, t)
	sshconnect.EnsureCommandIssued(mock.RunCmdsCommands, "Copy etcd systemd service to host", hostInEtcdRole.PublicIP, t)
	sshconnect.EnsureCommandIssued(mock.RunCmdsCommands, "Enable and start etcd service", hostInEtcdRole.PublicIP, t)

	sshconnect.EnsureNoCommandsIssued(mock.RunCmdsCommands, hostConfigs[1].PublicIP, t)
}

type GeneratesCertsMock struct {
	ca                  *certs.CA
	etcdCert            *certs.EtcdCert
	isEtcdCertGenerated bool
	certs.GeneratesCerts
}

func (g *GeneratesCertsMock) GetCA() *certs.CA { return g.ca }
func (g *GeneratesCertsMock) GenAdminClientCertificate() (*certs.AdminClientCert, error) {
	return nil, fmt.Errorf("Not yet implemented")
}
func (g *GeneratesCertsMock) GenEtcdCertificate(hosts []string) (*certs.EtcdCert, error) {
	g.isEtcdCertGenerated = true
	return g.etcdCert, nil
}

func NewGeneratesCertsMock() *GeneratesCertsMock {
	ca := certs.CA{
		CertBytes: []byte("CA_CERT"),
		KeyBytes:  []byte("CA_KEY")}

	etcdCert := certs.EtcdCert{
		PrivateKeyBytes: []byte("ETCD_KEY"),
		PublicKeyBytes:  []byte("ETCD_CERT")}

	return &GeneratesCertsMock{
		ca:                  &ca,
		etcdCert:            &etcdCert,
		isEtcdCertGenerated: false}
}
