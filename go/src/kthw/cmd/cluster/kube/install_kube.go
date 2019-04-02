package kube

import (
	"fmt"
	"kthw/certs"
	"kthw/cmd/infra/server"
	"kthw/cmd/sshconnect"
)

// InstallOnHosts installs kubernetes to all hosts in role controller or worker
func InstallOnHosts(
	serverConfigs []*server.Config,
	ssh sshconnect.SSHOperations,
	certsLoader certs.CertificateLoader,
	certsGenerator certs.GeneratesCerts) error {

	controllerConfigs := server.SelectHostsInRole(serverConfigs, "controller")
	if len(controllerConfigs) <= 0 {
		return fmt.Errorf("List of provided hosts didn't contain a host with role controller, but one controller is required")
	}
	if len(controllerConfigs) > 1 {
		return fmt.Errorf("List of provided hosts contains more than one host with role controller, but one controller is allowed")
	}
	controllerNode := &ControllerNode{Config: controllerConfigs[0]}

	etcdHosts := server.SelectHostsInRole(serverConfigs, "etcd")
	var etcdNodes []*EtcdNode
	for _, etcdHost := range etcdHosts {
		node := &EtcdNode{EndpointURL: fmt.Sprintf("https://%s:2379", etcdHost.PrivateIP)}
		etcdNodes = append(etcdNodes, node)
	}

	deployPodsToControllerNode := len(serverConfigs) <= 1
	InstallControllerNode(controllerNode, etcdNodes, ssh, certsLoader, certsGenerator, deployPodsToControllerNode)

	workerConfigs := server.SelectHostsInRole(serverConfigs, "worker")
	for _, workerConfig := range workerConfigs {
		InstallWorkerNode(workerConfig, controllerNode, ssh)
	}

	return nil
}
