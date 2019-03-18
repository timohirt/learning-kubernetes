package etcd

import (
	"bytes"
	"fmt"
	"kthw/certs"
	"kthw/cmd/common"
	"kthw/cmd/infra/server"
	"kthw/cmd/sshconnect"
	"os"
	"strings"
)

// InstallOnHost selects hosts with role 'etcd' and installs etcd on it.
// Currently, only single node clusters are supported.
func InstallOnHost(hostConfigs []*server.Config, ssh sshconnect.SSHOperations, generateCerts certs.GeneratesCerts) error {
	etcdHosts := selectEtcdHosts(hostConfigs)
	if len(etcdHosts) == 0 {
		return fmt.Errorf("List of provided hosts didn't contain a host with role etcd")
	}

	if len(etcdHosts) > 1 {
		return fmt.Errorf("Several hosts have role etcd. Currently, only single node etcd clusters are allowed")
	}

	for _, etcdHost := range etcdHosts {
		host := etcdHost.PublicIP
		certHostnames := []string{"localhost", etcdHost.PrivateIP}
		etcdCert, err := generateCerts.GenEtcdCertificate(certHostnames)
		if err != nil {
			return fmt.Errorf("Error while generating etcd certificate: %s", err)
		}
		commands := &sshconnect.Commands{
			Commands: []sshconnect.Command{
				downloadEtcd(host),
				unpackAndInstall(host),
				uploadEtcdCertPrivateKey(host, etcdCert),
				uploadEtcdCertPublicKey(host, etcdCert),
				uploadCAPublicKey(host, generateCerts.GetCA()),
				uploadSystemdService(etcdHost),
				enableAndStartEtcdSystemdService(host)},
			LogOutput: true}
		err = ssh.RunCmds(commands)
		if err != nil {
			return err
		}
	}
	return nil
}

func selectEtcdHosts(hostConfigs []*server.Config) []*server.Config {
	var etcdHosts []*server.Config
	for _, host := range hostConfigs {
		if common.ArrayContains(host.Roles, "etcd") {
			etcdHosts = append(etcdHosts, host)
		}
	}
	return etcdHosts
}

func uploadEtcdCertPublicKey(host string, etcdCert *certs.EtcdCert) *sshconnect.CopyFileCommand {
	return &sshconnect.CopyFileCommand{
		Host:        host,
		FileContent: bytes.NewReader(etcdCert.PublicKeyBytes),
		FilePath:    "/etc/kubernetes/pki/etcd.crt",
		Description: "Upload etcd certificate public key to /etc/kubernetes/pki/etcd.crt"}
}

func uploadEtcdCertPrivateKey(host string, etcdCert *certs.EtcdCert) *sshconnect.CopyFileCommand {
	return &sshconnect.CopyFileCommand{
		Host:        host,
		FileContent: bytes.NewReader(etcdCert.PrivateKeyBytes),
		FilePath:    "/etc/kubernetes/pki/etcd.key",
		Description: "Upload etcd certificate private key to /etc/kubernetes/pki/etcd.key"}
}

func uploadCAPublicKey(host string, ca *certs.CA) *sshconnect.CopyFileCommand {
	return &sshconnect.CopyFileCommand{
		Host:        host,
		FileContent: bytes.NewReader(ca.CertBytes),
		FilePath:    "/etc/kubernetes/pki/ca.crt",
		Description: "Upload CA certificate public key to /etc/kubernetes/pki/ca.crt"}
}

func uploadSystemdService(hostConfig *server.Config) *sshconnect.CopyFileCommand {
	params := SystemdServiceParameters{
		PrivateIP: hostConfig.PrivateIP,
		NodeName:  hostConfig.Name}
	systemdService, err := GenerateSystemdService(params)
	if err != nil {
		fmt.Printf("Error generating systemd service! %s\n", err)
		os.Exit(1)
	}

	return &sshconnect.CopyFileCommand{
		Host:        hostConfig.PublicIP,
		FileContent: strings.NewReader(systemdService),
		FilePath:    "/etc/systemd/system/etcd.service",
		Description: "Copy etcd systemd service to host"}
}

func downloadEtcd(host string) *sshconnect.ShellCommand {
	return &sshconnect.ShellCommand{
		CommandLine: "curl -L https://github.com/etcd-io/etcd/releases/download/v3.3.12/etcd-v3.3.12-linux-amd64.tar.gz -o /tmp/etcd-v3.3.12.tar.gz",
		Host:        host,
		Description: "Download etcd binary"}
}

func unpackAndInstall(host string) *sshconnect.ShellCommand {
	return &sshconnect.ShellCommand{
		CommandLine: "tar xzf /tmp/etcd-v3.3.12.tar.gz -C /tmp && mv /tmp/etcd-v3.3.12*/etcd* /usr/local/bin/",
		Host:        host,
		Description: "Untar etcd archive and copy to /usr/local/bin"}
}

func enableAndStartEtcdSystemdService(host string) *sshconnect.ShellCommand {
	return &sshconnect.ShellCommand{
		CommandLine: "systemctl daemon-reload && systemctl enable etcd && systemctl restart etcd",
		Host:        host,
		Description: "Enable and start etcd service"}
}
