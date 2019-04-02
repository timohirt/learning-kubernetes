package kube

import (
	"fmt"
	"kthw/cmd/common"
	"kthw/cmd/infra/server"
	"kthw/cmd/sshconnect"
)

// InstallWorkerNode gets a kubeadm join command from controllerName and runs it on host
func InstallWorkerNode(host *server.Config, controllerNode *ControllerNode, ssh sshconnect.SSHOperations) error {
	if common.ArrayContains(host.Roles, "controller") {
		return fmt.Errorf("Installing worker to service in role controller is not allowed")
	}

	if !common.ArrayContains(host.Roles, "worker") {
		return fmt.Errorf("Installing worker on server not possible. It is not in role worker")
	}

	res, _ := ssh.RunCmd(getClusterJoinCommand(controllerNode), false)
	fmt.Println(res)

	res, err := ssh.RunCmd(runClusterJoinCommand(host, res), false)
	if err != nil {
		return nil
	}
	fmt.Println("Result - ", res)

	return nil
}

func getClusterJoinCommand(controller *ControllerNode) *sshconnect.ShellCommand {
	return &sshconnect.ShellCommand{
		CommandLine: "kubeadm token create --print-join-command",
		Host:        controller.Config.PublicIP,
		Description: "Get cluster join command from controller"}
}

func runClusterJoinCommand(host *server.Config, joinCommand string) *sshconnect.ShellCommand {
	return &sshconnect.ShellCommand{
		CommandLine: joinCommand,
		Host:        host.PublicIP,
		Description: "Running join cluster command on worker"}
}
