package kube

import (
	"bytes"
	"fmt"
	"kthw/certs"
	"kthw/cmd/infra/server"
	"kthw/cmd/sshconnect"
	"os"
	"strings"
)

// InstallOnHosts installs kubernetes to all hosts in role controller or worker
func InstallOnHosts(hostConfigs []*server.Config, ssh sshconnect.SSHOperations, certsLoader certs.CertificateLoader) error {
	controllerHosts := server.SelectHostsInRole(hostConfigs, "controller")
	if len(controllerHosts) <= 0 {
		return fmt.Errorf("List of provided hosts didn't contain a host with role controller, but one controller is required")
	}
	if len(controllerHosts) > 1 {
		return fmt.Errorf("List of provided hosts contains more than one host with role controller, but one controller is allowed")
	}

	etcdHosts := server.SelectHostsInRole(hostConfigs, "etcd")
	var etcdNodes []EtcdNode
	for _, etcdHost := range etcdHosts {
		node := EtcdNode{EndpointURL: fmt.Sprintf("https://%s:2379", etcdHost.PrivateIP)}
		etcdNodes = append(etcdNodes, node)
	}

	for _, controllerHost := range controllerHosts {
		allCommands := baseSetup(controllerHost, etcdNodes, certsLoader)

		if len(controllerHosts) == 1 {
			allCommands = append(allCommands, untaintMaster(controllerHost))
		}

		allCommands = append(allCommands, NewCalicoNetworkingAddOn().getCommands(controllerHost)...)

		allCommands = append(allCommands, NewKubernetesDashboardAddOn().getCommands(controllerHost)...)

		commands := &sshconnect.Commands{
			Commands:  allCommands,
			LogOutput: true}

		err := ssh.RunCmds(commands)
		if err != nil {
			return err
		}
	}
	return nil
}

func baseSetup(controllerHost *server.Config, etcdNodes []EtcdNode, certsLoader certs.CertificateLoader) []sshconnect.Command {
	host := controllerHost.PublicIP

	etcdClientCert, err := certsLoader.LoadEtcdClientCert()
	if err != nil {
		fmt.Printf("Error while loading etcd client certificate: %s", err)
		os.Exit(1)
	}
	ca, err := certsLoader.LoadCA()
	if err != nil {
		fmt.Printf("Error while loading CA certificate: %s", err)
		os.Exit(1)
	}

	return []sshconnect.Command{
		uploadEtcdClientCertPrivateKey(host, etcdClientCert),
		uploadEtcdClientCertPublicKey(host, etcdClientCert),
		uploadCAPublicKey(host, ca),
		uploadCAPrivateKey(host, ca),
		uploadKubeadmMasterConfig(controllerHost, etcdNodes),
		installKubernetesCluster(controllerHost),
		setupKubectl(controllerHost),
		openFirewall(controllerHost),
	}
}

func uploadKubeadmMasterConfig(hostConfig *server.Config, etcdNodes []EtcdNode) *sshconnect.CopyFileCommand {
	kubeAdmParams := NewKubeAdmParams(hostConfig, etcdNodes)
	kubeadmConfig, err := GenerateKubeadmControllerConfig(kubeAdmParams)
	if err != nil {
		fmt.Printf("Error generating kubeadm controller config! %s\n", err)
		os.Exit(1)
	}

	return &sshconnect.CopyFileCommand{
		Host:        hostConfig.PublicIP,
		FileContent: strings.NewReader(kubeadmConfig),
		FilePath:    "/etc/kubernetes/kubeadm-controller.conf",
		Description: "Copy kubeadm config"}
}

func uploadEtcdClientCertPublicKey(host string, etcdClientCert *certs.EtcdClientCert) *sshconnect.CopyFileCommand {
	return &sshconnect.CopyFileCommand{
		Host:        host,
		FileContent: bytes.NewReader(etcdClientCert.PublicKeyBytes),
		FilePath:    "/etc/kubernetes/pki/etcd-client.crt",
		Description: "Upload etcd client certificate public key to /etc/kubernetes/pki/etcd-client.crt"}
}

func uploadEtcdClientCertPrivateKey(host string, etcdClientCert *certs.EtcdClientCert) *sshconnect.CopyFileCommand {
	return &sshconnect.CopyFileCommand{
		Host:        host,
		FileContent: bytes.NewReader(etcdClientCert.PrivateKeyBytes),
		FilePath:    "/etc/kubernetes/pki/etcd-client.key",
		Description: "Upload etcd client certificate private key to /etc/kubernetes/pki/etcd-client.key"}
}

func uploadCAPublicKey(host string, ca *certs.CA) *sshconnect.CopyFileCommand {
	return &sshconnect.CopyFileCommand{
		Host:        host,
		FileContent: bytes.NewReader(ca.CertBytes),
		FilePath:    "/etc/kubernetes/pki/ca.crt",
		Description: "Upload CA certificate public key to /etc/kubernetes/pki/ca.crt"}
}

func uploadCAPrivateKey(host string, ca *certs.CA) *sshconnect.CopyFileCommand {
	return &sshconnect.CopyFileCommand{
		Host:        host,
		FileContent: bytes.NewReader(ca.KeyBytes),
		FilePath:    "/etc/kubernetes/pki/ca.key",
		Description: "Upload CA certificate private key to /etc/kubernetes/pki/ca.key"}
}

func removeKubernetesCluster(hostConfig *server.Config) *sshconnect.ShellCommand {
	return &sshconnect.ShellCommand{
		CommandLine: "kubeadm reset -f",
		Host:        hostConfig.PublicIP,
		Description: "Removing kubernetes cluster"}
}

func installKubernetesCluster(hostConfig *server.Config) *sshconnect.ShellCommand {
	return &sshconnect.ShellCommand{
		CommandLine: "kubeadm init --config /etc/kubernetes/kubeadm-controller.conf",
		Host:        hostConfig.PublicIP,
		Description: "Install kubernetes cluster"}
}

func setupKubectl(hostConfig *server.Config) *sshconnect.ShellCommand {
	return &sshconnect.ShellCommand{
		CommandLine: "mkdir -p $HOME/.kube && cp -i /etc/kubernetes/admin.conf $HOME/.kube/config && sudo chown $(id -u):$(id -g) $HOME/.kube/config",
		Host:        hostConfig.PublicIP,
		Description: "Setup Kubectl"}
}

func openFirewall(hostConfig *server.Config) *sshconnect.ShellCommand {
	return &sshconnect.ShellCommand{
		CommandLine: fmt.Sprintf("ufw allow from %s to %s && ufw allow 6443", podNetworkCIDR, hostConfig.PublicIP),
		Host:        hostConfig.PublicIP,
		Description: "Open firewall pod network -> public IP and :6443 -> public IP"}
}

func untaintMaster(hostConfig *server.Config) *sshconnect.ShellCommand {
	return &sshconnect.ShellCommand{
		CommandLine: "kubectl taint nodes --all node-role.kubernetes.io/master-",
		Host:        hostConfig.PublicIP,
		Description: "Untaint master, allow pod scheduling on master node"}
}
