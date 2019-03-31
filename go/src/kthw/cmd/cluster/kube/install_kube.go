package kube

import (
	"fmt"
	"kthw/certs"
	"kthw/cmd/infra/server"
	"kthw/cmd/sshconnect"
)

// InstallOnHosts installs kubernetes to all hosts in role controller or worker
func InstallOnHosts(
	hostConfigs []*server.Config,
	ssh sshconnect.SSHOperations,
	certsLoader certs.CertificateLoader,
	certsGenerator certs.GeneratesCerts) error {

	controllerHosts := server.SelectHostsInRole(hostConfigs, "controller")
	if len(controllerHosts) <= 0 {
		return fmt.Errorf("List of provided hosts didn't contain a host with role controller, but one controller is required")
	}
	if len(controllerHosts) > 1 {
		return fmt.Errorf("List of provided hosts contains more than one host with role controller, but one controller is allowed")
	}

	etcdHosts := server.SelectHostsInRole(hostConfigs, "etcd")
	var etcdNodes []*EtcdNode
	for _, etcdHost := range etcdHosts {
		node := &EtcdNode{EndpointURL: fmt.Sprintf("https://%s:2379", etcdHost.PrivateIP)}
		etcdNodes = append(etcdNodes, node)
	}

	deployPodsToControllerNode := len(hostConfigs) <= 1

	for _, controllerHost := range controllerHosts {
		InstallControllerNode(controllerHost, etcdNodes, ssh, certsLoader, certsGenerator, deployPodsToControllerNode)
	}
	return nil
}
