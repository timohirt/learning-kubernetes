package kube

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

type ControllerNode struct {
	Config *server.Config
}

// InstallControllerNode installs a Kubernets controller on host.
// If you only have one controller running, consider setting
// `runPodsOnController = true` and deploy pods to controller.
func InstallControllerNode(
	controllerNode *ControllerNode,
	etcdNodes []*EtcdNode,
	ssh sshconnect.SSHOperations,
	certsLoader certs.CertificateLoader,
	certGenerator certs.GeneratesCerts,
	runPodsOnController bool) error {

	host := controllerNode.Config

	allCommands := baseSetup(host, etcdNodes, certsLoader, certGenerator)

	if runPodsOnController {
		allCommands = append(allCommands, untaintController(host))
	}

	allCommands = append(allCommands, NewCalicoNetworkingAddOn().getCommands(host)...)

	allCommands = append(allCommands, NewKubernetesDashboardAddOn().getCommands(host)...)

	commands := &sshconnect.Commands{
		Commands:  allCommands,
		LogOutput: true}

	err := ssh.RunCmds(commands)
	return err
}

func baseSetup(
	config *server.Config,
	etcdNodes []*EtcdNode,
	certsLoader certs.CertificateLoader,
	certGenerator certs.GeneratesCerts) []sshconnect.Command {

	host := config.PublicIP
	ca, err := certsLoader.LoadCA()
	if err != nil {
		fmt.Printf("Error while loading CA certificate: %s", err)
		os.Exit(1)
	}

	commands := uploadEtcdClientCert(host, certGenerator)
	commands = append(commands,
		uploadCAPublicKey(host, ca),
		uploadCAPrivateKey(host, ca),
		uploadKubeadmconfig(config, etcdNodes),
		installKubernetesCluster(config),
		setupKubectl(config),
		openFirewall(config))

	return commands
}

func uploadKubeadmconfig(config *server.Config, etcdNodes []*EtcdNode) *sshconnect.CopyFileCommand {
	kubeAdmParams := NewKubeAdmParams(config, etcdNodes)
	kubeadmConfig, err := GenerateKubeadmControllerConfig(kubeAdmParams)
	if err != nil {
		fmt.Printf("Error generating kubeadm controller config! %s\n", err)
		os.Exit(1)
	}

	return &sshconnect.CopyFileCommand{
		Host:        config.PublicIP,
		FileContent: strings.NewReader(kubeadmConfig),
		FilePath:    "/etc/kubernetes/kubeadm-controller.conf",
		Description: "Copy kubeadm config"}
}

func uploadEtcdClientCert(host string, certGenerator certs.GeneratesCerts) []sshconnect.Command {
	etcdClientCert, err := certGenerator.GenEtcdClientCertificate()
	common.WhenErrPrintAndExit(err)

	return []sshconnect.Command{
		&sshconnect.CopyFileCommand{
			Host:        host,
			FileContent: bytes.NewReader(etcdClientCert.PublicKeyBytes),
			FilePath:    "/etc/kubernetes/pki/etcd-client.crt",
			Description: "Upload etcd client certificate public key to /etc/kubernetes/pki/etcd-client.crt"},
		&sshconnect.CopyFileCommand{
			Host:        host,
			FileContent: bytes.NewReader(etcdClientCert.PrivateKeyBytes),
			FilePath:    "/etc/kubernetes/pki/etcd-client.key",
			Description: "Upload etcd client certificate private key to /etc/kubernetes/pki/etcd-client.key"}}
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

func removeKubernetesCluster(config *server.Config) *sshconnect.ShellCommand {
	return &sshconnect.ShellCommand{
		CommandLine: "kubeadm reset -f",
		Host:        config.PublicIP,
		Description: "Removing kubernetes cluster"}
}

func installKubernetesCluster(config *server.Config) *sshconnect.ShellCommand {
	return &sshconnect.ShellCommand{
		CommandLine: "kubeadm init --config /etc/kubernetes/kubeadm-controller.conf",
		Host:        config.PublicIP,
		Description: "Install kubernetes cluster"}
}

func setupKubectl(config *server.Config) *sshconnect.ShellCommand {
	return &sshconnect.ShellCommand{
		CommandLine: "mkdir -p $HOME/.kube && cp -i /etc/kubernetes/admin.conf $HOME/.kube/config && sudo chown $(id -u):$(id -g) $HOME/.kube/config",
		Host:        config.PublicIP,
		Description: "Setup Kubectl"}
}

func openFirewall(config *server.Config) *sshconnect.ShellCommand {
	return &sshconnect.ShellCommand{
		CommandLine: fmt.Sprintf("ufw allow from %s to %s && ufw allow 6443", podNetworkCIDR, config.PublicIP),
		Host:        config.PublicIP,
		Description: "Open firewall pod network -> public IP and :6443 -> public IP"}
}

func untaintController(config *server.Config) *sshconnect.ShellCommand {
	return &sshconnect.ShellCommand{
		CommandLine: "kubectl taint nodes --all node-role.kubernetes.io/master-",
		Host:        config.PublicIP,
		Description: "Untaint controller, allow pod scheduling on controller node"}
}
